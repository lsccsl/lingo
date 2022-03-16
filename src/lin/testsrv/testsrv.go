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
}

func ConstructTestSrv(addrLocal string, addrRemote string, srvId int64) *TestSrv {
	s := &TestSrv{
		srvId:srvId,
	}

	var err error
	lsn, err := net.Listen("tcp", addrLocal)
	s.tcpAcpt, err = lsn.Accept()
	if err != nil {
		lin_common.LogDebug(err)
	}

	msg := recvProtoMsg(s.tcpAcpt)
	msgReport := msg.(*msgpacket.MSG_SRV_REPORT)
	if msgReport == nil {
		return nil
	}

	s.tcpDial, err = net.Dial("tcp", addrRemote)
	if err != nil {
		lin_common.LogDebug(err)
	}

	msgReport = &msgpacket.MSG_SRV_REPORT{SrvId: srvId}
	s.tcpDial.Write(msgpacket.ProtoPacketToBin(msgpacket.MSG_TYPE__MSG_SRV_REPORT, msgReport))

	Global_wg.Add(2)

	go s.go_tcpDial()
	go s.go_tcpAcpt()

	return s
}

func (pthis*TestSrv)go_tcpDial() {
	recvBuf := bytes.NewBuffer(make([]byte, 0, MAX_PACK_LEN))
	TmpBuf := make([]byte, G_MTU)

	var seq int64 = 0

	Loop:
	for{
		msgRPC := &msgpacket.MSG_RPC{
			MsgId:lin_common.GenUUID64_V4(),
			MsgType:int32(msgpacket.MSG_TYPE__MSG_TEST),
		}
		seq ++
		msgTest := &msgpacket.MSG_TEST{
			Id:pthis.srvId,
			Seq:seq,
		}
		var err error
		msgRPC.MsgBin, err = proto.Marshal(msgTest)
		if err != nil {
			continue
		}
		pthis.tcpDial.Write(msgpacket.ProtoPacketToBin(msgpacket.MSG_TYPE__MSG_RPC, msgRPC))


		readSize, err := pthis.tcpDial.Read(TmpBuf)
		if nil != err {
			break Loop
		}
		recvBuf.Write(TmpBuf[0:readSize])
		//fmt.Println("tcp read:", readSize, err)

		READ_LOOP:
		for ; recvBuf.Len() >= PACK_HEAD_SIZE; {
			binHead := recvBuf.Bytes()[0:PACK_HEAD_SIZE]

			packLen := binary.LittleEndian.Uint32(binHead[0:4])
			packType := binary.LittleEndian.Uint16(binHead[4:6])

			if recvBuf.Len() < int(packLen){
				break READ_LOOP
			}

			binBody := recvBuf.Bytes()[6:packLen]

			protoMsg := msgpacket.ParseProtoMsg(binBody, int32(packType))
			fmt.Println(protoMsg)

			recvBuf.Next(int(packLen))
		}
	}

	Global_wg.Done()
}

func (pthis*TestSrv)go_tcpAcpt() {
	recvBuf := bytes.NewBuffer(make([]byte, 0, MAX_PACK_LEN))
	TmpBuf := make([]byte, G_MTU)

	Loop:
	for{
		readSize, err := pthis.tcpAcpt.Read(TmpBuf)
		if nil != err {
			break Loop
		}
		recvBuf.Write(TmpBuf[0:readSize])
		//fmt.Println("tcp read:", readSize, err)

		READ_LOOP:
		for ; recvBuf.Len() >= PACK_HEAD_SIZE; {
			binHead := recvBuf.Bytes()[0:PACK_HEAD_SIZE]

			packLen := binary.LittleEndian.Uint32(binHead[0:4])
			packType := binary.LittleEndian.Uint16(binHead[4:6])

			if recvBuf.Len() < int(packLen){
				break READ_LOOP
			}

			binBody := recvBuf.Bytes()[6:packLen]

			protoMsg := msgpacket.ParseProtoMsg(binBody, int32(packType))
			recvBuf.Next(int(packLen))

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
				pthis.tcpDial.Write(msgpacket.ProtoPacketToBin(msgpacket.MSG_TYPE__MSG_RPC_RES, msgRPCRes))
			}

		}
	}

	Global_wg.Done()
}