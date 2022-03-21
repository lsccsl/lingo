package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"lin/lin_common"
	"lin/msgpacket"
	"net"
	"os"
	"syscall"
	"time"
)

const MAX_PACK_LEN int = 65535
const G_MTU int = 1536
const PACK_HEAD_SIZE int = 6


type INTER_MSG_TYPE int32

type interMsg struct {
	msgtype INTER_MSG_TYPE
	msgdata interface{}
}
type interSendMsg struct {
	msgtype msgpacket.MSG_TYPE
	msgproto proto.Message
}
type interSendMsgLoop struct {
	loopCount int
	loopSmall int
}

const (
	INTER_MSG_TYPE_sendmsg = 1
	INTER_MSG_TYPE_sendmsg_loop = 2
	INTER_MSG_TYPE_recvmsg = 3
)

const CHAN_BUF_LEN = 100

type ClientTcpInfo struct{
	id int64
	addr string
	tcpCon *net.TCPConn
	needDecrypt bool
	msgChan chan *interMsg
	ByteSend int
	ByteRecv int
	seq int64

	rttTotal int64
	rttAver int64
	rttMax int64
	rttMin int64
	diffArrive int64
	diffProcess int64
	diffBack int64
	testCount int64
	reconnectCount int64
	testCountTotal int64
}
var globalTcpInfo *ClientTcpInfo

type PackHead struct{
	packLen uint32
	packType uint16
}




func CheckError(err error)bool{
	if nil == err{
		return true
	}

	switch t:=err.(type){
	case net.Error:
		{
			if t.Timeout() {
				//lin_common.LogDebug(" time out")
			} else if t.Temporary() {
				lin_common.LogDebug(" temporary")
			} else {
				lin_common.LogDebug(" other err:", t)
				netOpErr, ok := t.(*net.OpError)
				if ok {
					lin_common.LogDebug(" net op err:", netOpErr)
					switch st := netOpErr.Err.(type){
					case *os.SyscallError:
						{
							lin_common.LogDebug("syscall err", st)
							switch sterr := st.Err.(type) {
							case syscall.Errno:
								lin_common.LogDebug("syscall errno:", sterr)
							default:
								lin_common.LogDebug("unknow sys call err:", sterr)
							}
						}
					default:
						lin_common.LogDebug("unknow net op err", st)
					}
				} else {
					lin_common.LogDebug(" unkonw other :", t)
				}
			}
		}
/*	case *net.OpError:
		lin_common.LogDebug(t)*/
	default:
		lin_common.LogDebug(t)
	}

	//lin_common.LogErr(err)
	if err == io.EOF{
		lin_common.LogDebug("io eof")
		return false
	}
	netErr, ok := err.(net.Error)
	if ok{
		if netErr.Timeout() {
			//lin_common.LogDebug("time out")
			return true
		}
		if netErr.Temporary() {
			//lin_common.LogDebug("temporary")
			return true
		}
		netOpErr, ok := netErr.(*net.OpError)
		if ok{
			switch t := netOpErr.Err.(type){
			case *net.DNSError:
				//lin_common.LogDebug("net.DNSError:", t)
				break
			case *os.SyscallError:
				if errno, ok := t.Err.(syscall.Errno); ok {
					switch errno {
					case syscall.ECONNREFUSED:
						//lin_common.LogDebug("connect refused")
						break
					case syscall.ETIMEDOUT:
						//lin_common.LogDebug("timeout")
						return true
						break
					case syscall.WSAECONNRESET:
						//lin_common.LogDebug("connection reset")
						break
					default:
						//lin_common.LogDebug("unknow err num:", errno)
						break
					}
				}
			case *net.UnknownNetworkError:
				//lin_common.LogDebug("net.UnknownNetworkError", t)
				break
			case *os.LinkError:
				//lin_common.LogDebug("os.LinkError", t)
				break
			case *os.PathError:
				//lin_common.LogDebug("os.PathError", t)
				break
			default:
				//lin_common.LogDebug("unknown err", t)
				break
			}
		}
	}
	return false
}

func (tcpInfo *ClientTcpInfo)GetSeq() int64 {
	tcpInfo.seq ++
	return tcpInfo.seq
}

func (tcpInfo *ClientTcpInfo)GoClientTcpProcess() {
	fmt.Println("GoClientTcpSend")

	count := 0

	defer func() {
		err := recover()
		if err != nil {
			lin_common.LogDebug("recover get err:", err)
		}
		lin_common.LogDebug("call defer exit coroutine")
		Global_wg.Done()
	}()

	//chTimer := time.After(time.Second * time.Duration(10))

	PROCESS_LOOP:
	for {
		select {
		case msg := <-tcpInfo.msgChan:
			{
				if msg == nil {
					break PROCESS_LOOP
				}
				if (msg.msgtype == INTER_MSG_TYPE_sendmsg) {
					count ++
				}
				tcpInfo.processMsg(msg)
			}

/*		case <-chTimer:
			{
				//fmt.Println("timeout")
				msg := &msgpacket.MSG_HEARTBEAT{Id:tcpInfo.id}
				go tcpInfo.TcpSend(msg)
				chTimer = time.After(time.Second * time.Duration(10))
			}*/
		}
	}
}

func (tcpInfo *ClientTcpInfo)processMsg(msg *interMsg) {
	switch msg.msgtype {
	case INTER_MSG_TYPE_sendmsg:
		sendMsg, ok := msg.msgdata.(*interSendMsg)
		if ok && sendMsg != nil {
			tcpInfo.processSendMsg(sendMsg)
		}
	case INTER_MSG_TYPE_sendmsg_loop:
		sendMsgLoop, ok := msg.msgdata.(*interSendMsgLoop)
		if ok && sendMsgLoop != nil {
			tcpInfo.processSendMsgLoop(sendMsgLoop)
		}
	}
}
func (tcpInfo *ClientTcpInfo)processSendMsg(msg *interSendMsg) {
	bin := tcpInfo.FormatMsg(msg.msgtype, msg.msgproto)
	tcpInfo.ByteSend += len(bin)
	tcpInfo.tcpCon.Write(bin)
}

func (pthis *ClientTcpInfo)processSendMsgLoop(msg *interSendMsgLoop) {
	defer func() {
		err := recover()
		if err != nil {
			lin_common.LogDebug("recover get err:", err)
		}
	}()

	pthis.testCount = 0
	pthis.rttTotal = 0
	pthis.rttAver = 0
	pthis.rttMax = 0
	pthis.rttMin = 0
	pthis.diffArrive = 0
	pthis.diffProcess = 0
	pthis.diffBack = 0

	var seq int64 = 0
	var maxSeq int64 = 0
	for i := 0; i < msg.loopCount; i ++ {
		for j := 0; j < msg.loopSmall; j ++ {
			msgTest := &msgpacket.MSG_TEST{}
			msgTest.Id = pthis.id
			msgTest.Str = fmt.Sprintf("%v_%v_%v", pthis.id, j, i)
			seq++
			msgTest.Seq = seq
			msgTest.Timestamp = time.Now().UnixMilli()
			//lin_common.LogDebug("send test:", msgTest)
			bin := pthis.FormatMsg(msgpacket.MSG_TYPE__MSG_TEST, msgTest)
			pthis.ByteSend += len(bin)
			pthis.tcpCon.Write(bin)
		}

		READ_LOOP:
		for k := 0; maxSeq < seq;  {
			select {
			case msgRes := <-pthis.msgChan:
			{
				tnow := time.Now().UnixMilli()
				msgTestRes, ok := msgRes.msgdata.(*msgpacket.MSG_TEST_RES)
				if !ok {
					continue
				}
				k ++
				maxSeq = msgTestRes.Seq
				diff := tnow - msgTestRes.Timestamp

				diffArrive := msgTestRes.TimestampArrive - msgTestRes.Timestamp
				diffProcess := msgTestRes.TimestampProcess - msgTestRes.TimestampArrive
				diffBack := tnow - msgTestRes.TimestampProcess

				pthis.diffArrive += diffArrive
				pthis.diffProcess += diffProcess
				pthis.diffBack += diffBack

				pthis.rttTotal += diff
				if pthis.rttMin == 0 {
					pthis.rttMin = diff
				} else if pthis.rttMin > diff {
					pthis.rttMin = diff
				}
				if pthis.rttMax < diff {
					pthis.rttMax = diff
				}
				pthis.testCount ++
				pthis.testCountTotal ++
				//lin_common.LogDebug("recv res:", msgRes.msgdata)
			}
			case <- time.After(time.Second * 15):
				break READ_LOOP
			}
		}

		if pthis.testCount >= 1{
			pthis.rttAver = pthis.rttTotal / pthis.testCount
		}

/*		if maxSeq < seq {
			lin_common.LogDebug("~~~~~~err seq:", maxSeq)
		}*/
	}
}

func (tcpInfo *ClientTcpInfo)FormatMsg(msgtype msgpacket.MSG_TYPE, msg proto.Message)[]byte{
	binMsg, _ := proto.Marshal(msg)
	var wb []byte
	var buf bytes.Buffer
	_ = binary.Write(&buf,binary.LittleEndian,uint32(6 + len(binMsg)))
	_ = binary.Write(&buf,binary.LittleEndian,uint16(msgtype))
	wb = buf.Bytes()
	wb = append(wb, binMsg...)

	return wb
}

// todo 省略 msgtype参数
func (tcpInfo *ClientTcpInfo)TcpSend(msg proto.Message) {
	var msgtype = msgpacket.MSG_TYPE(msgpacket.GetMsgTypeByMsgInstance(msg))
	tcpInfo.msgChan <- &interMsg{
		msgtype:INTER_MSG_TYPE_sendmsg,
		msgdata: &interSendMsg{
			msgtype:msgtype,
			msgproto:msg,
		},
	}
}

func (tcpInfo *ClientTcpInfo)TcpSendLoop(loopCount int, loopSmall int) {
	tcpInfo.msgChan <- &interMsg{
		msgtype:INTER_MSG_TYPE_sendmsg_loop,
		msgdata:&interSendMsgLoop{loopCount, loopSmall},
	}
}



func (tcpInfo *ClientTcpInfo)GoClientTcpRead(){
	fmt.Println("GoClientTcpRead")

	defer func() {
		err := recover()
		if err != nil {
			lin_common.LogDebug("recover get err:", err)
		}
		lin_common.LogDebug("call defer exit coroutine")
		Global_wg.Done()
	}()

	recvBuf := bytes.NewBuffer(make([]byte, 0, MAX_PACK_LEN))
	TmpBuf := make([]byte, G_MTU)

	curHead := PackHead{0,0}

	for{
		tcpInfo.tcpCon.SetReadDeadline(time.Now().Add(time.Second * 120))
		readSize, err := tcpInfo.tcpCon.Read(TmpBuf)
		if !CheckError(err){
			tcpInfo.reconnectCount ++
			con , _ := net.Dial("tcp", tcpInfo.addr)
			tcpInfo.tcpCon = con.(*net.TCPConn)

			msg := &msgpacket.MSG_LOGIN{}
			msg.Id = tcpInfo.id
			bin := tcpInfo.FormatMsg(msgpacket.MSG_TYPE__MSG_LOGIN, msg)
			tcpInfo.tcpCon.Write(bin)
			continue
		}
		tcpInfo.ByteRecv += readSize
		recvBuf.Write(TmpBuf[0:readSize])
		//fmt.Println("tcp read:", readSize, err)

		READ_LOOP:
		for ; recvBuf.Len() >= PACK_HEAD_SIZE; {
			binHead := recvBuf.Bytes()[0:PACK_HEAD_SIZE]

			curHead.packLen = binary.LittleEndian.Uint32(binHead[0:4])
			curHead.packType = binary.LittleEndian.Uint16(binHead[4:6])

			if recvBuf.Len() < int(curHead.packLen){
				curHead.packLen = 0
				break READ_LOOP
			}

			binBody := recvBuf.Bytes()[6:curHead.packLen]

			protoMsg := msgpacket.ParseProtoMsg(binBody, int32(curHead.packType))
			switch msgpacket.MSG_TYPE(curHead.packType) {
			case msgpacket.MSG_TYPE__MSG_HEARTBEAT_RES:
			case msgpacket.MSG_TYPE__MSG_TEST_RES:
			default:
				lin_common.LogDebug(msgpacket.MSG_TYPE(curHead.packType), " proto msg:", protoMsg)
			}

			recvBuf.Next(int(curHead.packLen))

			tcpInfo.msgChan <- &interMsg{
				msgtype:INTER_MSG_TYPE_recvmsg,
				msgdata:protoMsg,
			}

			curHead.packLen = 0
		}
	}

	fmt.Println("exit GoClientTcpRead", time.Now())
}

func (pthis *ClientTcpInfo)ClientDump() (str string) {
	count := pthis.testCount
	if count <= 0 {
		count = 1
	}
	str = fmt.Sprintf("aver:%v min:%v max:%v reconnect:%v" +
		" diffArrive:%v diffProcess:%v count:%v",
		pthis.rttAver, pthis.rttMin, pthis.rttMax, pthis.reconnectCount,
		pthis.diffArrive / count,
		pthis.diffProcess / count,
		pthis.testCountTotal)
	return
}

func StartClient(id int64, addr string) *ClientTcpInfo {
	conn, err := net.Dial("tcp", addr/*"192.168.2.129:2003"*/)
	fmt.Println(conn, err)

	tcpInfo := &ClientTcpInfo{
		id:id,
		addr:addr,
		tcpCon : conn.(*net.TCPConn),
		msgChan : make(chan * interMsg, CHAN_BUF_LEN),
		ByteSend:0,
		ByteRecv:0,
		rttTotal:0,
		rttAver:0,
		rttMax:0,
		rttMin:0,
		diffArrive:0,
		diffProcess:0,
		diffBack:0,
		reconnectCount:0,
		testCountTotal:0,
	}
	globalTcpInfo = tcpInfo

	Global_wg.Add(2)
	go tcpInfo.GoClientTcpRead()
	go tcpInfo.GoClientTcpProcess()

	msg := &msgpacket.MSG_LOGIN{}
	msg.Id = id
	tcpInfo.TcpSend(msg)

	Global_cliMgr.ClientMgrAdd(tcpInfo)

	return tcpInfo
}
