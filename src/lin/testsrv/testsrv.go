package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"lin/lin_common"
	"lin/msgpacket"
	"lin/tcp"
	"net"
)
const MAX_PACK_LEN int = 65535
const G_MTU int = 1536

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

	totalWriteRpc int64
	totalRpcDial int64
}

func ConstructTestSrv(addrLocal string, addrRemote string, srvId int64) *TestSrv {
	s := &TestSrv{
		srvId:srvId,
		recvBuf : bytes.NewBuffer(make([]byte, 0, MAX_PACK_LEN)),
		TmpBuf : make([]byte, G_MTU),
		seq : 0,
		addrRemote : addrRemote,
		addrLocal : addrLocal,
	}
	Global_TestSrvMgr.TestSrvMgrAdd(s)

	Global_wg.Add(2)
	go s.go_tcpAcpt()
	go s.go_tcpDial()

	return s
}

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
	_, err = pthis.tcpDial.Write(msgpacket.ProtoPacketToBin(msgpacket.MSG_TYPE__MSG_RPC, msgRPC))
	if err != nil {
		lin_common.LogErr("write tcp err:", pthis.DialConnectionID)
		return err
	}
	//lin_common.LogDebug(" write", msgRPC, err)

	RSP_LOOP:
	for {
		var msgRsp proto.Message
		msgRsp, err = recvProtoMsg(pthis.tcpDial, 3, 3)
		//lin_common.LogDebug(msgRsp, err)
		if err != nil {
			lin_common.LogDebug(" err:", err, " fail count:", pthis.totalWriteRpc - pthis.totalRpcDial, " total rpc:", pthis.totalWriteRpc)
			return err
		}

		switch t := msgRsp.(type) {
		case *msgpacket.MSG_RPC_RES:
			if t.MsgId == msgRPC.MsgId {
				pthis.totalRpcDial ++
				break RSP_LOOP
			} else {
				lin_common.LogDebug(" other id:", msgRPC.MsgId)
			}
		default:
			lin_common.LogDebug(t, " err:", err)
		}
	}

	return nil
}

func (pthis*TestSrv)go_tcpDial() {
	defer func() {
		lin_common.LogDebug("exit dial")
		err := recover()
		if err != nil {
			lin_common.LogDebug(err)
		}
	}()

	for{
		err := pthis.TestSrvDial()
		if err != nil || pthis.tcpDial == nil{
			lin_common.LogDebug("rpc err:", err)
			conn, err := net.Dial("tcp", pthis.addrRemote)
			if err != nil || conn == nil{
				lin_common.LogDebug("dial err:", err)
				continue
			}
			pthis.tcpDial = conn.(*net.TCPConn)
			msgReport := &msgpacket.MSG_SRV_REPORT{SrvId: pthis.srvId}
			pthis.tcpDial.Write(msgpacket.ProtoPacketToBin(msgpacket.MSG_TYPE__MSG_SRV_REPORT, msgReport))

			REPORT_RES_LOOP:
			for {
				msg, err := recvProtoMsg(pthis.tcpDial, 3, 3)
				if err != nil {
					break REPORT_RES_LOOP
				}
				switch t := msg.(type) {
				case *msgpacket.MSG_SRV_REPORT_RES:
					pthis.DialConnectionID = tcp.TCP_CONNECTION_ID(t.TcpConnId)
					lin_common.LogDebug(" suc recv srv report:", pthis.DialConnectionID)
					break REPORT_RES_LOOP
				default:
				}
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
				lin_common.LogDebug("MSG_HEARTBEAT:", t.Id)
				msgHBRsp := &msgpacket.MSG_HEARTBEAT_RES{Id:t.Id}
				pthis.tcpAcpt.Write(msgpacket.ProtoPacketToBin(msgpacket.MSG_TYPE__MSG_HEARTBEAT_RES, msgHBRsp))
				continue
			}
		default:
			continue
		}

		fmt.Println(protoMsg)
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

	lsn, _ := net.Listen("tcp", pthis.addrLocal)

	for{
		err := pthis.TestSrvAcpt()
		if err != nil || pthis.tcpAcpt == nil {
			lin_common.LogDebug("acpt read err:", err)
			if pthis.tcpAcpt != nil {
				pthis.tcpAcpt.Close()
			}
			conn, err := lsn.Accept()
			if err != nil || conn == nil{
				lin_common.LogDebug("acpt err:", err)
				lin_common.LogDebug(err)
				continue
			}
			pthis.tcpAcpt = conn.(*net.TCPConn)
			REPORT_LOOP:
			for {
				msg, err := recvProtoMsg(pthis.tcpAcpt, 3, 3)
				if err != nil {
					break REPORT_LOOP
				}
				switch t := msg.(type) {
				case *msgpacket.MSG_SRV_REPORT:
					pthis.AcptConnectionID = tcp.TCP_CONNECTION_ID(t.TcpConnId)
					lin_common.LogDebug(" suc recv srv report")
					break REPORT_LOOP
				default:
				}
			}
		}
	}

	Global_wg.Done()
}