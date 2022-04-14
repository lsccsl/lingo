package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"lin/lin_common"
	"lin/msgpacket"
	"lin/tcp"
	"math"
	"net"
	"os"
	"time"
)
const MAX_PACK_LEN int = 65535
const G_MTU int = 1536

type TestSrvStatic struct {
	totalWriteRpc int64
	totalRpcDial int64

	totalRpcRecv int64

	totalRedial int64
	totalReAcpt int64

	minRTTDialRpc int64
	maxRTTDialRpc int64
	totalRTTRpc int64
}

type TestSrv struct {
	tcpDial *net.TCPConn
	tcpAcpt *net.TCPConn
	srvId int64
	recvBuf *bytes.Buffer
	TmpBuf []byte
	seq int64

	DialConnectionID tcp.TCP_CONNECTION_ID
	AcptConnectionID tcp.TCP_CONNECTION_ID

	addrRemote string
	addrLocal string
	local_port int

	TestSrvStatic

	AutoRedial bool
}

func ConstructTestSrv(addrLocal string, local_port int, addrRemote string, srvId int64) *TestSrv {
	s := &TestSrv{
		srvId:srvId,
		recvBuf : bytes.NewBuffer(make([]byte, 0, MAX_PACK_LEN)),
		TmpBuf : make([]byte, G_MTU),
		seq : 0,
		addrRemote : addrRemote,
		addrLocal : addrLocal,
		local_port : local_port,
		AutoRedial : true,
	}
	s.minRTTDialRpc = math.MaxInt64
	s.maxRTTDialRpc = 0

	Global_TestSrvMgr.TestSrvMgrAdd(s)

	Global_wg.Add(2)
	go s.go_tcpAcpt()

	return s
}

func (pthis*TestSrv)TestSrvBeginDial(){
	go pthis.go_tcpDial()
}

const RPC_RETRY_COUNT = 18
const RPC_READ_TIMEOUT = 10

func (pthis*TestSrv)TestSrvDial() (err interface{}) {
	if pthis.tcpDial == nil {
		return lin_common.GenErr(0, "no tcp dial conn")
	}
	defer func() {
		errR := recover()
		if errR != nil {
			err = errR
		}
	}()
	msgRPC := &msgpacket.MSG_RPC{
		MsgId:lin_common.GenUUID64_V4(),
		MsgType:int32(msgpacket.MSG_TYPE__MSG_TEST),
		Timestamp:time.Now().UnixMilli(),
		TimeoutWait: RPC_RETRY_COUNT * RPC_READ_TIMEOUT * 1000,
	}
	pthis.seq ++
	msgTest := &msgpacket.MSG_TEST{
		Id:pthis.srvId,
		Seq:pthis.seq,
	}
	msgRPC.MsgBin, err = proto.Marshal(msgTest)
	if err != nil {
		return
	}
	pthis.totalWriteRpc ++
	tBeginTime := time.Now().UnixMilli()
	_, err = pthis.tcpDial.Write(msgpacket.ProtoPacketToBin(msgpacket.MSG_TYPE__MSG_RPC, msgRPC))
	if err != nil {
		lin_common.LogErr("write tcp err:", pthis.DialConnectionID, " srv:", pthis.srvId, " err:", err)
		return err
	}
	//lin_common.LogDebug(" write", msgRPC, err)

	RSP_LOOP:
	for {
		var msgRsp proto.Message
		msgRsp, err = recvProtoMsg(pthis.tcpDial, RPC_RETRY_COUNT, RPC_READ_TIMEOUT)
		//lin_common.LogDebug(msgRsp, err)
		if err != nil {
			//lin_common.LogDebug(" err:", err, " fail count:", pthis.totalWriteRpc - pthis.totalRpcDial, " total rpc:", pthis.totalWriteRpc)
			return err
		}

		switch t := msgRsp.(type) {
		case *msgpacket.MSG_RPC_RES:
			if t.MsgId == msgRPC.MsgId {
				pthis.totalRpcDial ++
				tEndTime := time.Now().UnixMilli()
				tDiff := tEndTime - tBeginTime
				if tDiff < pthis.minRTTDialRpc {
					pthis.minRTTDialRpc = tDiff
				}
				if tDiff > pthis.maxRTTDialRpc {
					pthis.maxRTTDialRpc = tDiff
				}
				pthis.totalRTTRpc += tDiff
				break RSP_LOOP
			} else {
				lin_common.LogDebug(" other id:", msgRPC.MsgId, " conn:", pthis.DialConnectionID, " srv:", pthis.srvId)
			}
		default:
			lin_common.LogDebug(t, " err:", err, " conn:", pthis.DialConnectionID, " srv:", pthis.srvId)
		}
	}

	return nil
}

func (pthis*TestSrv)tcpReDial() {

	httpAddDial(&ServerFromHttp{
		SrvID: pthis.srvId,
		IP: Global_testCfg.local_ip,
		Port: pthis.local_port,
	})

	pthis.DialConnectionID = 0
	conn, err := net.Dial("tcp", pthis.addrRemote)
	pthis.totalRedial ++
	if err != nil || conn == nil{
		lin_common.LogDebug("dial err:", err, " conn:", pthis.DialConnectionID, " srv:", pthis.srvId)
		return
	}
	pthis.tcpDial = conn.(*net.TCPConn)
	msgReport := &msgpacket.MSG_SRV_REPORT{SrvId: pthis.srvId}
	pthis.tcpDial.Write(msgpacket.ProtoPacketToBin(msgpacket.MSG_TYPE__MSG_SRV_REPORT, msgReport))

REPORT_RES_LOOP:
	for {
		msg, err := recvProtoMsg(pthis.tcpDial, 3, 10)
		if err != nil {
			break REPORT_RES_LOOP
		}
		switch t := msg.(type) {
		case *msgpacket.MSG_SRV_REPORT_RES:
			pthis.DialConnectionID = tcp.TCP_CONNECTION_ID(t.TcpConnId)
			//lin_common.LogDebug(" suc recv srv report", " conn:", pthis.DialConnectionID, " srv:", pthis.srvId)
			break REPORT_RES_LOOP
		default:
		}
	}
	//time.Sleep(time.Second * 1000)

}

func (pthis*TestSrv)go_tcpDial() {
	defer func() {
		lin_common.LogDebug("exit dial")
		err := recover()
		if err != nil {
			lin_common.LogDebug(err)
		}
	}()

	pthis.tcpReDial()

	for{
		if !pthis.AutoRedial {
			time.Sleep(time.Second * 3)
			continue
		}
		err := pthis.TestSrvDial()
		if err != nil || pthis.tcpDial == nil{
			if pthis.tcpDial != nil {
				pthis.tcpDial.Close()
			}
			lin_common.LogDebug("rpc err:", err, " conn:", pthis.DialConnectionID, " srv:", pthis.srvId)
			if pthis.AutoRedial {
				pthis.tcpReDial()
			}
		}
	}

	Global_wg.Done()
}

func (pthis*TestSrv)TestSrvAcpt() (err interface{}) {
	if pthis.tcpAcpt == nil {
		return lin_common.GenErr(0, "no tcp acpt conn")
	}
	defer func() {
		errR := recover()
		if errR != nil {
			lin_common.LogDebug("acpt read:", errR)
			err = errR
		}
	}()
	readSize, err := pthis.tcpAcpt.Read(pthis.TmpBuf)
	if nil != err {
		return err
	}
	pthis.totalRpcRecv ++
	_, err = pthis.recvBuf.Write(pthis.TmpBuf[0:readSize])
	if err != nil {
		return err
	}

	READ_LOOP:
	for ; pthis.recvBuf.Len() >= PACK_HEAD_SIZE; {
		binHead := pthis.recvBuf.Bytes()[0:PACK_HEAD_SIZE]

		packLen := binary.LittleEndian.Uint32(binHead[0:4])
		packType := binary.LittleEndian.Uint16(binHead[4:6])

		if pthis.recvBuf.Len() < int(packLen){
			break READ_LOOP
		}

		binBody := pthis.recvBuf.Bytes()[6:packLen]

		protoMsg := msgpacket.ParseProtoMsg(binBody, int32(packType))
		pthis.recvBuf.Next(int(packLen))
		switch t := protoMsg.(type) {
		case *msgpacket.MSG_RPC:
			{

			}
		case *msgpacket.MSG_HEARTBEAT:
			{
				//lin_common.LogDebug("MSG_HEARTBEAT:", t.Id)
				msgHBRsp := &msgpacket.MSG_HEARTBEAT_RES{Id:t.Id}
				pthis.tcpAcpt.Write(msgpacket.ProtoPacketToBin(msgpacket.MSG_TYPE__MSG_HEARTBEAT_RES, msgHBRsp))
				continue
			}
		default:
			continue
		}

		//fmt.Println(protoMsg)
		var msgTest *msgpacket.MSG_TEST
		var msgRPC *msgpacket.MSG_RPC
		{
			msgRPC = protoMsg.(*msgpacket.MSG_RPC)
			if msgRPC == nil {
				continue
			}
			msgR := msgpacket.ParseProtoMsg(msgRPC.MsgBin, int32(msgRPC.MsgType))
			if msgR == nil {
				continue
			}
			msgTest = msgR.(*msgpacket.MSG_TEST)
			if msgTest == nil {
				continue
			}
		}

		{
			msgTestRes := &msgpacket.MSG_TEST_RES{
				Id:msgTest.Id,
				Seq:msgTest.Seq,
				Str:"````msgTest.Str!!!!",
				//Str:msgTest.Str,
			}
			msgRPCRes := &msgpacket.MSG_RPC_RES{
				MsgId:msgRPC.MsgId,
				MsgType:int32(msgpacket.MSG_TYPE__MSG_TEST_RES),
			}
			msgRPCRes.MsgBin, err = proto.Marshal(msgTestRes)
			if err != nil {
				continue
			}
			pthis.tcpAcpt.Write(msgpacket.ProtoPacketToBin(msgpacket.MSG_TYPE__MSG_RPC_RES, msgRPCRes))
		}
	}

	return nil
}

func (pthis*TestSrv)go_tcpAcpt() {
	defer func() {
		lin_common.LogDebug("exit dial")
		err := recover()
		if err != nil {
			lin_common.LogDebug(err)
		}
	}()

	if pthis.srvId == 599 {
		lin_common.LogDebug("srv:", pthis.srvId, " begin listen:", pthis.addrLocal)
	}

	lsn, err := net.Listen("tcp", pthis.addrLocal)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	for{
		var err interface{}

		if pthis.srvId == 599 {
			lin_common.LogDebug("srv:", pthis.srvId, " begin accept")
		}

		conn, err := lsn.Accept()

		if pthis.srvId == 599 {
			lin_common.LogDebug("srv:", pthis.srvId, " accept err:", err)
		}

/*		if pthis.tcpAcpt != nil {
			pthis.tcpAcpt.Close()
		}*/
		pthis.totalReAcpt ++
		if err != nil || conn == nil{
			lin_common.LogDebug("acpt err:", err, " conn:", pthis.AcptConnectionID, " srv:", pthis.srvId)
			continue
		}

		pthis.tcpAcpt = conn.(*net.TCPConn)

		go func(){
			corConn := conn
			defer func() {
				corConn.Close()
			}()
			REPORT_LOOP:
			for {
				msg, err := recvProtoMsg(pthis.tcpAcpt, 3, 10)
				if err != nil {
					break REPORT_LOOP
				}
				switch t := msg.(type) {
				case *msgpacket.MSG_SRV_REPORT:
					pthis.AcptConnectionID = tcp.TCP_CONNECTION_ID(t.TcpConnId)
					//lin_common.LogDebug(" suc recv srv report", " conn:", pthis.AcptConnectionID, " srv:", pthis.srvId)
					break REPORT_LOOP
				default:
				}
			}
			for {
				err = pthis.TestSrvAcpt()
				if err != nil{
					lin_common.LogDebug("acpt read err:", err, " conn:", pthis.AcptConnectionID, " srv:", pthis.srvId)
					break
				}
			}
			pthis.AcptConnectionID = 0
		}()
	}

	Global_wg.Done()
}