package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"lin/lin_common"
	"lin/msgpacket"
	"net"
)
const MAX_PACK_LEN int = 65535
const G_MTU int = 1536

type TestSrv struct {
	tcpDial net.Conn
	tcpAcpt net.Conn
	srvId int64
	recvBuf *bytes.Buffer
	TmpBuf []byte
	seq int64

	addrRemote string
	addrLocal string
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

/*	var err error
	lsn, err := net.Listen("tcp", addrLocal)
	s.tcpAcpt, err = lsn.Accept()
	if err != nil {
		lin_common.LogDebug(err)
	}

	msg := recvProtoMsg(s.tcpAcpt)
	msgReport := msg.(*msgpacket.MSG_SRV_REPORT)
	if msgReport == nil {
		return nil
	}*/

/*	s.tcpDial, err = net.Dial("tcp", addrRemote)
	if err != nil {
		lin_common.LogDebug(err)
	}

	msgReport = &msgpacket.MSG_SRV_REPORT{SrvId: srvId}
	s.tcpDial.Write(msgpacket.ProtoPacketToBin(msgpacket.MSG_TYPE__MSG_SRV_REPORT, msgReport))
*/
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
		err = recover()
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
	_, err = pthis.tcpDial.Write(msgpacket.ProtoPacketToBin(msgpacket.MSG_TYPE__MSG_RPC, msgRPC))

	RSP_LOOP:
	for {
		msgRsp, err := recvProtoMsg(pthis.tcpDial)
		if err != nil {
			return err
		}

		switch t := msgRsp.(type) {
		case *msgpacket.MSG_RPC_RES:
			if t.MsgId == msgRPC.MsgId {
				break RSP_LOOP
			}
		default:
			fmt.Println(t)
		}
	}

	return nil
}

func (pthis*TestSrv)go_tcpDial() {

	for{
		err := pthis.TestSrvDial()
		if err != nil || pthis.tcpDial == nil{
			pthis.tcpDial, err = net.Dial("tcp", pthis.addrRemote)
			if err != nil {
				continue
			}
			msgReport := &msgpacket.MSG_SRV_REPORT{SrvId: pthis.srvId}
			pthis.tcpDial.Write(msgpacket.ProtoPacketToBin(msgpacket.MSG_TYPE__MSG_SRV_REPORT, msgReport))
		}
	}

	Global_wg.Done()
}

func (pthis*TestSrv)TestSrvAcpt() (err interface{}) {
	if pthis.tcpAcpt == nil {
		return lin_common.GenErr(0, "no tcp acpt conn")
	}
	defer func() {
		err = recover()
	}()
	readSize, err := pthis.tcpAcpt.Read(pthis.TmpBuf)
	if nil != err {
		return err
	}
	pthis.recvBuf.Write(pthis.TmpBuf[0:readSize])
	//fmt.Println("tcp read:", readSize, err)

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

	lsn, _ := net.Listen("tcp", pthis.addrLocal)

	for{
		err := pthis.TestSrvAcpt()
		if err != nil || pthis.tcpAcpt == nil {
			if pthis.tcpAcpt != nil {
				pthis.tcpAcpt.Close()
			}
			pthis.tcpAcpt, err = lsn.Accept()
			if err != nil {
				lin_common.LogDebug(err)
			}
			REPORT_LOOP:
			for {
				msg, err := recvProtoMsg(pthis.tcpAcpt)
				if err != nil {
					break REPORT_LOOP
				}
				switch msg.(type) {
				case *msgpacket.MSG_SRV_REPORT:
					lin_common.LogDebug(" suc recv srv report")
					break REPORT_LOOP
				default:
				}
			}
		}
	}

	Global_wg.Done()
}