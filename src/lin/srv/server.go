package main

import (
	"github.com/golang/protobuf/proto"
	"lin/lin_common"
	msgpacket "lin/msgpacket"
	"lin/tcp"
	"sync/atomic"
	"time"
)

type Server struct {
	srvMgr *ServerMgr
	srvID         int64
	connDialID    tcp.TCP_CONNECTION_ID
	connAcptID    tcp.TCP_CONNECTION_ID
	chSrvProtoMsg chan *interProtoMsg
	chInterMsg chan interface{}
	heartbeatIntervalSec int
	rpcMgr *RPCManager

	isStopProcess int32
}

type interMsgSrvReport struct {
	tcpAccept *tcp.TcpConnection
}
type interMsgConnDial struct {
	tcpDial *tcp.TcpConnection
}

func ConstructServer(srvMgr *ServerMgr, srvID int64, heartbeatIntervalSec int)*Server {
	s := &Server{
		srvMgr:srvMgr,
		srvID:srvID,
		chSrvProtoMsg:make(chan *interProtoMsg, 100),
		chInterMsg:make(chan interface{}, 100),
		heartbeatIntervalSec:heartbeatIntervalSec,
		rpcMgr:ConstructRPCManager(),
		isStopProcess:0,
	}
	go s.go_serverProcess()
	return s
}

func (pthis*Server) go_serverProcess() {
	defer func() {
		atomic.StoreInt32(&pthis.isStopProcess, 1)
		err := recover()
		if err != nil {
			lin_common.LogErr(err)
		}
	}()

	chTimer := time.After(time.Second * time.Duration(pthis.heartbeatIntervalSec))

	MSG_LOOP:
	for {
		select {
		case ProtoMsg := <- pthis.chSrvProtoMsg:
			{
				if ProtoMsg == nil {
					break MSG_LOOP
				}
				pthis.processServerMsg(ProtoMsg)
			}

		case interMsg := <- pthis.chInterMsg:
			{
				switch t:= interMsg.(type) {
				case *interMsgSrvReport:
					pthis.processSrvReport(t.tcpAccept)
				case *interMsgConnDial:
					pthis.processDailConnect(t.tcpDial)
				}
			}

		case <-chTimer:
			{
				//log.LogErr("send test from dial, srvid:", pthis.srvID, pthis.heartbeatIntervalSec)
				chTimer = time.After(time.Second * time.Duration(pthis.heartbeatIntervalSec))
				//send heartbeat
				msgTest := &msgpacket.MSG_HEARTBEAT{}
				msgTest.Id = pthis.srvID
				pthis.srvMgr.tcpMgr.TcpConnectSendProtoMsg(pthis.connDialID, msgpacket.MSG_TYPE__MSG_HEARTBEAT, msgTest)
			}
		}
	}

	atomic.StoreInt32(&pthis.isStopProcess, 1)
	close(pthis.chSrvProtoMsg)
	close(pthis.chInterMsg)
}

func (pthis*Server) ServerClose() {
	pthis.srvMgr.tcpMgr.TcpMgrCloseConn(pthis.connAcptID)
	pthis.srvMgr.tcpMgr.TcpMgrCloseConn(pthis.connDialID)

	pthis.chSrvProtoMsg <- nil
}

func (pthis*Server) ServerCloseAndDelDialData() {
	pthis.srvMgr.tcpMgr.TcpDialDelDialData(pthis.srvID)
	pthis.ServerClose()
}

func (pthis*Server)processSrvReport(tcpAccept *tcp.TcpConnection){
	pthis.connAcptID = tcpAccept.TcpConnectionID()

	lin_common.LogDebug(pthis.srvID, " ", pthis)
}

func (pthis*Server)processDailConnect(tcpDial *tcp.TcpConnection){
	pthis.connDialID = tcpDial.TcpConnectionID()

	lin_common.LogDebug(pthis.srvID, " ", pthis)
}

func (pthis*Server)PushInterMsg(msg interface{}){
	if atomic.LoadInt32(&pthis.isStopProcess) == 1 {
		return
	}
	pthis.chInterMsg <- msg
}
func (pthis*Server)PushProtoMsg(msgType msgpacket.MSG_TYPE, protoMsg proto.Message, tcpConnID tcp.TCP_CONNECTION_ID){
	if atomic.LoadInt32(&pthis.isStopProcess) == 1 {
		return
	}
	pthis.chSrvProtoMsg <- &interProtoMsg{
		msgType:msgType,
		protoMsg:protoMsg,
		tcpConnID:tcpConnID,
	}
}

func (pthis*Server)Go_ProcessRPC(tcpConn *tcp.TcpConnection, msg *msgpacket.MSG_RPC, msgBody proto.Message) {
	var msgRes proto.Message = nil
	switch t:= msgBody.(type) {
	case *msgpacket.MSG_TEST:
		{
			msgRes = pthis.processRPCTest(t)
		}
	}

	msgRPCRes := &msgpacket.MSG_RPC_RES{
		MsgId:msg.MsgId,
		MsgType:msg.MsgType,
		ResCode:msgpacket.RESPONSE_CODE_RESPONSE_CODE_NONE,
	}

	if msgRes != nil {
		var err error
		msgRPCRes.MsgBin, err = proto.Marshal(msgRes)
		if err != nil {
			lin_common.LogErr(err)
		}
	}
	tcpConn.TcpConnectSendBin(msgpacket.ProtoPacketToBin(msgpacket.MSG_TYPE__MSG_RPC_RES, msgRPCRes))
}
func (pthis*Server)processRPCRes(tcpConn *tcp.TcpConnection, msg *msgpacket.MSG_RPC_RES, msgBody proto.Message) {
	defer func() {
		err := recover()
		if err != nil {
			lin_common.LogErr(err)
		}
	}()
	if pthis.rpcMgr == nil {
		return
	}
	rreq := pthis.rpcMgr.RPCManagerFindReq(msg.MsgId)
	if rreq == nil {
		return
	}
	if rreq.chNtf != nil {
		rreq.chNtf <- msgBody
	}
}

func (pthis*Server)processRPCTest(msg *msgpacket.MSG_TEST) *msgpacket.MSG_TEST_RES {
	lin_common.LogDebug(msg)
	return &msgpacket.MSG_TEST_RES{Id: msg.Id, Str:msg.Str}
}

// SendRPC_Async @brief will block timeoutMilliSec
func (pthis*Server)SendRPC_Async(msgType msgpacket.MSG_TYPE, protoMsg proto.Message, timeoutMilliSec int) proto.Message {

	defer func() {
		err := recover()
		if err != nil {
			lin_common.LogErr(err)
		}
	}()

	if pthis.rpcMgr == nil{
		return nil
	}

	msgRPC := &msgpacket.MSG_RPC{
		MsgId:lin_common.GenUUID64_V4(),
		MsgType:int32(msgType),
	}
	var err error
	msgRPC.MsgBin, err = proto.Marshal(protoMsg)
	if err != nil {
		lin_common.LogErr(err)
		return nil
	}

	rreq := pthis.rpcMgr.RPCManagerAddReq(msgRPC.MsgId)

	pthis.srvMgr.tcpMgr.TcpConnectSendProtoMsg(pthis.connDialID, msgpacket.MSG_TYPE__MSG_RPC, msgRPC)

	var res proto.Message = nil
	select{
	case resCh := <-rreq.chNtf:
		res, _ = resCh.(proto.Message)
	case <-time.After(time.Millisecond * time.Duration(timeoutMilliSec)):
	}

	pthis.rpcMgr.RPCManagerDelReq(msgRPC.MsgId)

	return res
}

func (pthis*Server)processServerMsg (interMsg * interProtoMsg){
	switch t:=interMsg.protoMsg.(type){
	case *msgpacket.MSG_HEARTBEAT:
		pthis.process_MSG_HEARTBEAT(interMsg.tcpConnID, t)
	case *msgpacket.MSG_HEARTBEAT_RES:
		pthis.process_MSG_HEARTBEAT_RES(t)
	}
}

func (pthis*Server) process_MSG_HEARTBEAT (tcpConnID tcp.TCP_CONNECTION_ID, protoMsg * msgpacket.MSG_HEARTBEAT) {
	lin_common.LogDebug(protoMsg)

	msgRes := &msgpacket.MSG_HEARTBEAT_RES{}
	msgRes.Id = protoMsg.Id
	pthis.srvMgr.tcpMgr.TcpConnectSendProtoMsg(tcpConnID, msgpacket.MSG_TYPE__MSG_HEARTBEAT_RES, msgRes)
}

func (pthis*Server) process_MSG_HEARTBEAT_RES (protoMsg * msgpacket.MSG_HEARTBEAT_RES) {
	lin_common.LogDebug(protoMsg)
}
