package main

import (
	"github.com/golang/protobuf/proto"
	"goserver/msgpacket"
	"goserver/test/tcp"
)

func (pthis*Server)ServerProcessRPC(tcpConn *tcp.TcpConnection, msgBody proto.Message) proto.Message {
	var msgRes proto.Message = nil
	switch t:= msgBody.(type) {
	case *msgpacket.MSG_TEST:
		{
			msgRes = pthis.processRPCTest(t)
		}
	}

	return msgRes
}

func (pthis*Server)processRPCTest(msg *msgpacket.MSG_TEST) *msgpacket.MSG_TEST_RES {
	return &msgpacket.MSG_TEST_RES{Id: msg.Id, Str:msg.Str, Seq: msg.Seq}
}
