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
const PACK_HEAD_SIZE int = 4


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
}

const (
	INTER_MSG_TYPE_sendmsg = 1
	INTER_MSG_TYPE_sendmsg_loop = 2
	INTER_MSG_TYPE_recvmsg = 3
)

const CHAN_BUF_LEN = 100

type ClientTcpInfo struct{
	id int64
	con net.Conn
	needDecrypt bool
	msgChan chan *interMsg
	ByteSend int
	ByteRecv int
	seq int64
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
	fmt.Println("CheckError:", err)
	if err == io.EOF{
		return false
	}
	netErr, ok := err.(net.Error)
	if ok{
		netOpErr, ok := netErr.(*net.OpError)
		if ok{
			switch t := netOpErr.Err.(type){
			case *net.DNSError:
				fmt.Println("net.DNSError:", t)
				break
			case *os.SyscallError:
				if errno, ok := t.Err.(syscall.Errno); ok {
					switch errno {
					case syscall.ECONNREFUSED:
						fmt.Println("connect refused")
						break
					case syscall.ETIMEDOUT:
						fmt.Println("timeout")
						return true
						break
					case syscall.WSAECONNRESET:
						fmt.Println("connection reset")
						break
					default:
						fmt.Println("unknow err num", errno)
						break
						//case syscall.E
					}
				}
			case *net.UnknownNetworkError:
				fmt.Println("net.UnknownNetworkError", t)
				break
			case *os.LinkError:
				fmt.Println("os.LinkError", t)
				break
			case *os.PathError:
				fmt.Println("os.PathError", t)
				break
			default:
				fmt.Println("unknown err", t)
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
			fmt.Println("recover get err:", err)
		}
		fmt.Println("call defer exit coroutine")
		Global_wg.Done()
	}()

	chTimer := time.After(time.Second * time.Duration(10))

	PROCESS_LOOP:
	for {
		//fmt.Println("begin wait msg", lin_common.GetGID(), count)
		select {
		case msg := <-tcpInfo.msgChan:
			{
				if msg == nil {
					break PROCESS_LOOP
				}
				if (msg.msgtype == INTER_MSG_TYPE_sendmsg) {
					count ++
					//fmt.Println("write", count)
				}
				tcpInfo.processMsg(msg)
			}

		case <-chTimer:
			{
				//fmt.Println("timeout")
				msg := &msgpacket.MSG_HEARTBEAT{Id:tcpInfo.id}
				go tcpInfo.TcpSend(msg)
				chTimer = time.After(time.Second * time.Duration(10))
			}
		}
		//fmt.Println("end wait msg", count)
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
	tcpInfo.con.Write(bin)
}

func (pthis *ClientTcpInfo)processSendMsgLoop(msg *interSendMsgLoop) {
	for i := 0; i < msg.loopCount; i ++ {
		for j := 0; j < 200; j ++ {
			msgTest := &msgpacket.MSG_TEST{}
			msgTest.Id = pthis.id
			msgTest.Str = fmt.Sprintf("%v_%v_%v", pthis.id, j, i)
			bin := pthis.FormatMsg(msgpacket.MSG_TYPE__MSG_TEST, msgTest)
			pthis.ByteSend += len(bin)
			pthis.con.Write(bin)
		}

		for k := 0; k < 200; k ++ {
			//msgRes := <-pthis.msgChan
			_ = <-pthis.msgChan
			//lin_common.LogDebug("recv res:", msgRes.msgdata)
		}

		if i % 1000 == 0 {
			lin_common.LogDebug(i)
		}
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

func (tcpInfo *ClientTcpInfo)TcpSendLoop(loopCount int) {
	tcpInfo.msgChan <- &interMsg{
		msgtype:INTER_MSG_TYPE_sendmsg_loop,
		msgdata:&interSendMsgLoop{loopCount},
	}
}



func (tcpInfo *ClientTcpInfo)GoClientTcpRead(){
	fmt.Println("GoClientTcpRead")

	defer func() {
		err := recover()
		if err != nil {
			fmt.Println("recover get err:", err)
		}
		fmt.Println("call defer exit coroutine")
		Global_wg.Done()
	}()

	recvBuf := bytes.NewBuffer(make([]byte, 0, MAX_PACK_LEN))
	TmpBuf := make([]byte, G_MTU)

	curHead := PackHead{0,0}

	Loop:
	for{
		readSize, err := tcpInfo.con.Read(TmpBuf)
		if !CheckError(err){
			break Loop
		}
		tcpInfo.ByteRecv += readSize
		recvBuf.Write(TmpBuf[0:readSize])
		//fmt.Println("tcp read:", readSize, err)

		READ_LOOP:
		for ; recvBuf.Len() >= PACK_HEAD_SIZE; {
			binHead := recvBuf.Bytes()[0:6]

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

func StartClient(id int64, addr string) *ClientTcpInfo {
	conn, err := net.Dial("tcp", addr/*"192.168.2.129:2003"*/)
	fmt.Println(conn, err)

	tcpInfo := &ClientTcpInfo{
		id:id,
		con : conn,
		msgChan : make(chan * interMsg, CHAN_BUF_LEN),
		ByteSend:0,
		ByteRecv:0,
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
