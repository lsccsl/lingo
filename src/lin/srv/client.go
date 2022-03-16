package main

import (
	"github.com/golang/protobuf/proto"
	"lin/lin_common"
	"lin/msgpacket"
	"sync/atomic"
)

type interClientMsgTcpWrite struct {
	bin []byte
}

type MAP_CLIENT_STATIC map[msgpacket.MSG_TYPE]int64
type Client struct {
	srvMgr *ServerMgr
	tcpConnID TCP_CONNECTION_ID
	tcpConn *TcpConnection
	clientID int64
	chClientProtoMsg chan *interProtoMsg
	chMsgWrite chan *interClientMsgTcpWrite
	isStopProcess int32
	mapStaticMsgRecv MAP_CLIENT_STATIC

	goRoutineID uint64
}

func ConstructClient(srvMgr *ServerMgr, tcpConn *TcpConnection,clientID int64) *Client {
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
	tcpConn.ConnType = TCP_CONNECTIOON_TYPE_client

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

		case tcpW := <- pthis.chMsgWrite:
			{
				_, err := pthis.tcpConn.TcpConnectSendBin(tcpW.bin)
				if err != nil {
					pthis.tcpConn.TcpConnectSetCloseReason(TCP_CONNECTION_CLOSE_REASON_writeerr)
					break MSG_LOOP
				}
			}
		}
	}

	atomic.StoreInt32(&pthis.isStopProcess, 1)
	if pthis.tcpConn != nil {
		pthis.tcpConn.TcpConnectClose()
	}
	if pthis.chClientProtoMsg != nil {
		close(pthis.chClientProtoMsg)
	}
	if pthis.chMsgWrite != nil {
		close(pthis.chMsgWrite)
	}
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

func (pthis*Client) ProcessProtoMsg(msgType msgpacket.MSG_TYPE, protoMsg proto.Message) {
	pthis.mapStaticMsgRecv[msgType] = pthis.mapStaticMsgRecv[msgType] + 1
	switch t:= protoMsg.(type){
	case *msgpacket.MSG_HEARTBEAT:
		pthis.process_MSG_HEARTBEAT(t)
	case *msgpacket.MSG_TEST:
		pthis.process_MSG_TEST(t)
	case *msgpacket.MSG_TCP_STATIC:
		pthis.process_MSG_TCP_STATIC(t)
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
	pthis.ClientSendProtoMsg(msgpacket.MSG_TYPE__MSG_HEARTBEAT_RES, msgRes)
}

func (pthis*Client) process_MSG_TEST (protoMsg * msgpacket.MSG_TEST) {
	msgRes := &msgpacket.MSG_TEST_RES{}
	msgRes.Id = protoMsg.Id
	msgRes.Str = protoMsg.Str
	msgRes.Seq = protoMsg.Seq
	pthis.ClientSendProtoMsg(msgpacket.MSG_TYPE__MSG_TEST_RES, msgRes)
}

func (pthis*Client) process_MSG_TCP_STATIC(protoMsg * msgpacket.MSG_TCP_STATIC) {
	lin_common.LogDebug(" seq:", protoMsg.Seq)
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
	pthis.ClientSendProtoMsg(msgpacket.MSG_TYPE__MSG_TCP_STATIC_RES, msgRes)
	//pthis.srvMgr.tcpMgr.TcpConnectSendProtoMsg(pthis.tcpConnID, msgpacket.MSG_TYPE__MSG_TCP_STATIC_RES, msgRes)
}

func (pthis*Client) ClientSendProtoMsg(msgType msgpacket.MSG_TYPE, protoMsg proto.Message) {
	if lin_common.GetGID() != pthis.goRoutineID {
		lin_common.LogErr(" should call ClientSendProtoMsg in the client process coroutine:", msgType)
		return
	}
	pthis.tcpConn.TcpConnectSendBin(msgpacket.ProtoPacketToBin(msgType, protoMsg))
}

func (pthis * Client)ClientSendProtoMsgSync(msgType msgpacket.MSG_TYPE, protoMsg proto.Message) {
	if atomic.LoadInt32(&pthis.isStopProcess) != 0 {
		return
	}
	bin := msgpacket.ProtoPacketToBin(msgType, protoMsg)
	tcpW := &interClientMsgTcpWrite{
		make([]byte,len(bin)),
	}
	copy(tcpW.bin, bin)
	pthis.chMsgWrite <- tcpW
}
