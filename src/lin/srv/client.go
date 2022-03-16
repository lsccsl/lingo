package main

import (
	"github.com/golang/protobuf/proto"
	"lin/lin_common"
	"lin/msgpacket"
	"lin/tcp"
	"sync/atomic"
)


type MAP_CLIENT_STATIC map[msgpacket.MSG_TYPE]int64
type Client struct {
	srvMgr    *ServerMgr
	tcpConnID tcp.TCP_CONNECTION_ID
	tcpConn   *tcp.TcpConnection
	clientID  int64
	chClientProtoMsg chan *interProtoMsg
	isStopProcess int32
	mapStaticMsgRecv MAP_CLIENT_STATIC

	goRoutineID uint64
}

func ConstructClient(srvMgr *ServerMgr, tcpConn *tcp.TcpConnection,clientID int64) *Client {
	if tcpConn == nil {
		return nil
	}
	c := &Client{
		srvMgr:srvMgr,
		tcpConnID:tcpConn.TcpConnectionID(),
		tcpConn:tcpConn,
		clientID:clientID,
		chClientProtoMsg:make(chan *interProtoMsg, 100),
		isStopProcess:0,
		mapStaticMsgRecv:make(MAP_CLIENT_STATIC),
	}

	tcpConn.ConnData = c
	tcpConn.ConnType = tcp.TCP_CONNECTIOON_TYPE_client

	go c.go_clientProcess()

	return c
}

func (pthis*Client) go_clientProcess() {
	defer func() {
		err := recover()
		if err != nil {
			lin_common.LogErr(err)
		}
	}()

	pthis.goRoutineID = lin_common.GetGID()

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
						lin_common.LogErr(err)
					}
				}()
				pthis.processClientMsg(ProtoMsg)
			}()
		}
	}

	atomic.StoreInt32(&pthis.isStopProcess, 1)
	if pthis.tcpConn != nil {
		pthis.tcpConn.TcpConnectClose()
	}
	if pthis.chClientProtoMsg != nil {
		close(pthis.chClientProtoMsg)
	}
}

func (pthis*Client) ClientGetConnectionID() tcp.TCP_CONNECTION_ID {
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

/*func (pthis*Client) ProcessProtoMsg(msgType msgpacket.MSG_TYPE, protoMsg proto.Message) {
	pthis.mapStaticMsgRecv[msgType] = pthis.mapStaticMsgRecv[msgType] + 1
	switch t:= protoMsg.(type){
	case *msgpacket.MSG_HEARTBEAT:
		pthis.process_MSG_HEARTBEAT(t)
	case *msgpacket.MSG_TEST:
		pthis.process_MSG_TEST(t)
	case *msgpacket.MSG_TCP_STATIC:
		pthis.process_MSG_TCP_STATIC(t)
	}
}*/

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
	msgRes := &msgpacket.MSG_HEARTBEAT_RES{}
	msgRes.Id = protoMsg.Id
	pthis.tcpConn.TcpConnectSendBin(msgpacket.ProtoPacketToBin(msgpacket.MSG_TYPE__MSG_HEARTBEAT_RES, msgRes))
}

func (pthis*Client) process_MSG_TEST (protoMsg * msgpacket.MSG_TEST) {
	msgRes := &msgpacket.MSG_TEST_RES{}
	msgRes.Id = protoMsg.Id
	msgRes.Str = protoMsg.Str
	msgRes.Seq = protoMsg.Seq
	pthis.tcpConn.TcpConnectSendBin(msgpacket.ProtoPacketToBin(msgpacket.MSG_TYPE__MSG_TEST_RES, msgRes))
}

func (pthis*Client) process_MSG_TCP_STATIC(protoMsg * msgpacket.MSG_TCP_STATIC) {
	lin_common.LogDebug(" seq:", protoMsg.Seq)
	tcpConn := pthis.srvMgr.tcpMgr.GetTcpConnection(pthis.tcpConnID)
	msgRes := &msgpacket.MSG_TCP_STATIC_RES{
		ByteRecv:tcpConn.ByteRecv,
		ByteProc:tcpConn.ByteProc,
		ByteSend:tcpConn.ByteSend,
	}
	msgRes.MapStaticMsgRecv = make(map[int32]int64)
	for key, val := range pthis.mapStaticMsgRecv {
		msgRes.MapStaticMsgRecv[int32(key)] = val
	}
	pthis.tcpConn.TcpConnectSendBin(msgpacket.ProtoPacketToBin(msgpacket.MSG_TYPE__MSG_TCP_STATIC_RES, msgRes))
}

