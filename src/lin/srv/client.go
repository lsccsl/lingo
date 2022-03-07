package main

import (
	"github.com/golang/protobuf/proto"
	"lin/log"
	"lin/msgpacket"
	"sync/atomic"
)



type Client struct {
	srvMgr *ServerMgr
	tcpConn *TcpConnection
	clientID int64
	chClientProtoMsg chan *interProtoMsg
	isStopProcess int32
}

func ConstructClient(srvMgr *ServerMgr, tcpConn *TcpConnection,clientID int64) *Client {
	c := &Client{
		srvMgr:srvMgr,
		tcpConn:tcpConn,
		clientID:clientID,
		chClientProtoMsg:make(chan *interProtoMsg, 100),
		isStopProcess:0,
	}

	go c.go_clientProcess()

	return c
}

func (pthis*Client) go_clientProcess() {
	defer func() {
		err := recover()
		if err != nil {
			log.LogErr(err)
		}
	}()

	MSG_LOOP:
	for {
		select {
		case ProtoMsg := <- pthis.chClientProtoMsg:
			if ProtoMsg == nil {
				break MSG_LOOP
			}
			pthis.processClientMsg(ProtoMsg)
		}
	}

	atomic.StoreInt32(&pthis.isStopProcess, 1)
	close(pthis.chClientProtoMsg)
}

func (pthis*Client) ClientGetConnection()*TcpConnection{
	return pthis.tcpConn
}
func (pthis*Client) ClientGetClientID()int64{
	return pthis.clientID
}

func (pthis*Client) ClientClose() {
	pthis.srvMgr.tcpMgr.TcpMgrCloseConn(pthis.tcpConn.TcpConnectionID())
	pthis.chClientProtoMsg <- nil
}

func (pthis*Client) PushClientMsg(msgType msgpacket.MSG_TYPE, protoMsg proto.Message) {
	if atomic.LoadInt32(&pthis.isStopProcess) == 1 {
		return
	}

	pthis.chClientProtoMsg <- &interProtoMsg{
		msgType:msgType,
		protoMsg:protoMsg,
	}
}

func (pthis*Client)processRPC(tcpConn * TcpConnection, msg proto.Message) proto.Message {
	return nil
}
func (pthis*Client)processRPCRes(tcpConn * TcpConnection, msg proto.Message) {
}

func (pthis*Client) processClientMsg (interMsg * interProtoMsg) {
	switch t:=interMsg.protoMsg.(type){
	case *msgpacket.MSG_TEST:
		pthis.process_MSG_TEST(t)
	}
}

func (pthis*Client) process_MSG_TEST (protoMsg * msgpacket.MSG_TEST) {
	log.LogDebug(protoMsg)

	msgRes := &msgpacket.MSG_TEST_RES{}
	msgRes.Id = protoMsg.Id
	pthis.tcpConn.TcpConnectSendProtoMsg(msgpacket.MSG_TYPE__MSG_TEST_RES, msgRes)
}