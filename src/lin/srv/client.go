package main

import (
	"github.com/golang/protobuf/proto"
	"lin/log"
	"lin/msgpacket"
	"sync/atomic"
)


type MAP_CLIENT_STATIC map[msgpacket.MSG_TYPE]int64
type Client struct {
	srvMgr *ServerMgr
	tcpConnID TCP_CONNECTION_ID
	clientID int64
	chClientProtoMsg chan *interProtoMsg
	isStopProcess int32
	mapStaticMsgRecv MAP_CLIENT_STATIC
}

func ConstructClient(srvMgr *ServerMgr, tcpConnID TCP_CONNECTION_ID,clientID int64) *Client {
	c := &Client{
		srvMgr:srvMgr,
		tcpConnID:tcpConnID,
		clientID:clientID,
		chClientProtoMsg:make(chan *interProtoMsg, 100),
		isStopProcess:0,
		mapStaticMsgRecv:make(MAP_CLIENT_STATIC),
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
			pthis.mapStaticMsgRecv[ProtoMsg.msgType] = pthis.mapStaticMsgRecv[ProtoMsg.msgType] + 1
			func(){
				defer func() {
					err := recover()
					if err != nil {
						log.LogErr(err)
					}
				}()
				pthis.processClientMsg(ProtoMsg)
			}()
		}
	}

	atomic.StoreInt32(&pthis.isStopProcess, 1)
	close(pthis.chClientProtoMsg)
}

func (pthis*Client) ClientGetConnectionID() TCP_CONNECTION_ID{
	return pthis.tcpConnID
}
func (pthis*Client) ClientGetClientID()int64{
	return pthis.clientID
}

func (pthis*Client) ClientClose() {
	pthis.srvMgr.tcpMgr.TcpMgrCloseConn(pthis.tcpConnID)
	pthis.chClientProtoMsg <- nil
}

func (pthis*Client) PushProtoMsg(msgType msgpacket.MSG_TYPE, protoMsg proto.Message) {
	if atomic.LoadInt32(&pthis.isStopProcess) == 1 {
		return
	}

	pthis.chClientProtoMsg <- &interProtoMsg{
		msgType:msgType,
		protoMsg:protoMsg,
	}
}

func (pthis*Client)Go_processRPC(tcpConn * TcpConnection, msg *msgpacket.MSG_RPC, msgBody proto.Message) {
}
func (pthis*Client)processRPCRes(tcpConn * TcpConnection, msg *msgpacket.MSG_RPC_RES, msgBody proto.Message) {
}

func (pthis*Client) processClientMsg (interMsg * interProtoMsg) {
	switch t:=interMsg.protoMsg.(type){
	case *msgpacket.MSG_HEARTBEAT:
		pthis.process_MSG_HEARTBEAT(t)
	case *msgpacket.MSG_TEST:
		pthis.process_MSG_TEST(t)
	case *msgpacket.MSG_TCP_STATIC:
		pthis.process_MSG_TCP_STATIC(t)
	}
}

func (pthis*Client) process_MSG_HEARTBEAT (protoMsg * msgpacket.MSG_HEARTBEAT) {
	//log.LogDebug(protoMsg)

	msgRes := &msgpacket.MSG_HEARTBEAT_RES{}
	msgRes.Id = protoMsg.Id
	pthis.srvMgr.tcpMgr.TcpConnectSendProtoMsg(pthis.tcpConnID, msgpacket.MSG_TYPE__MSG_HEARTBEAT_RES, msgRes)
}

func (pthis*Client) process_MSG_TEST (protoMsg * msgpacket.MSG_TEST) {
	//log.LogDebug(protoMsg)

	msgRes := &msgpacket.MSG_TEST_RES{}
	msgRes.Id = protoMsg.Id
	msgRes.Str = protoMsg.Str
	pthis.srvMgr.tcpMgr.TcpConnectSendProtoMsg(pthis.tcpConnID, msgpacket.MSG_TYPE__MSG_TEST_RES, msgRes)
}

func (pthis*Client) process_MSG_TCP_STATIC(protoMsg * msgpacket.MSG_TCP_STATIC) {
	tcpConn := pthis.srvMgr.tcpMgr.getTcpConnection(pthis.tcpConnID)
	msgRes := &msgpacket.MSG_TCP_STATIC_RES{
		ByteRecv:tcpConn.ByteRecv,
		ByteProc:tcpConn.ByteProc,
		ByteSend:tcpConn.ByteSend,
	}
	msgRes.MapStaticMsgRecv = make(map[int32]int64)
	for key, val := range pthis.mapStaticMsgRecv {
		msgRes.MapStaticMsgRecv[int32(key)] = val
	}
	pthis.srvMgr.tcpMgr.TcpConnectSendProtoMsg(pthis.tcpConnID, msgpacket.MSG_TYPE__MSG_TCP_STATIC_RES, msgRes)
}