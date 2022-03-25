package main

import (
	"github.com/golang/protobuf/proto"
	"lin/lin_common"
	msgpacket "lin/msgpacket"
	"lin/tcp"
	"sync/atomic"
	"time"
)

type ServerStatic struct {
	totalRPCPacket int64
	totalRPCPacketLast int64
	timestamp float64
}

type Server struct {
	srvMgr *ServerMgr
	srvID         int64
	connDial    *tcp.TcpConnection
	connAcpt    *tcp.TcpConnection
	chSrvProtoMsg chan *interProtoMsg
	chInterMsg chan interface{}
	heartbeatIntervalSec int
	rpcMgr *RPCManager

	isStopProcess int32

	ServerStatic
}

type interMsgSrvReport struct {
	tcpAccept *tcp.TcpConnection
}
type interMsgConnDial struct {
	tcpDial *tcp.TcpConnection
}
type interMsgConnClose struct {
	tcpConn *tcp.TcpConnection
}

func ConstructServer(srvMgr *ServerMgr, connDial *tcp.TcpConnection, connAcpt *tcp.TcpConnection, srvID int64, heartbeatIntervalSec int)*Server {
	s := &Server{
		srvMgr:srvMgr,
		srvID:srvID,
		connDial:connDial,
		connAcpt:connAcpt,
		chSrvProtoMsg:make(chan *interProtoMsg, 100),
		chInterMsg:make(chan interface{}, 100),
		heartbeatIntervalSec:heartbeatIntervalSec,
		rpcMgr:ConstructRPCManager(),
		isStopProcess:0,
	}
	if s.connDial == nil && s.connAcpt == nil {
		lin_common.LogErr("connDial and connAcpt is nil:", srvID)
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
				case *interMsgConnClose:
					pthis.processConnClose(t.tcpConn)
				}
			}

		case <-chTimer:
			{
				//log.LogErr("send test from dial, srvid:", pthis.srvID, pthis.heartbeatIntervalSec)
				chTimer = time.After(time.Second * time.Duration(pthis.heartbeatIntervalSec))
				//send heartbeat
				if pthis.connDial != nil {
					msgHeartBeat := &msgpacket.MSG_HEARTBEAT{}
					msgHeartBeat.Id = pthis.srvMgr.srvID
					pthis.connDial.TcpConnectSendBin(msgpacket.ProtoPacketToBin(msgpacket.MSG_TYPE__MSG_HEARTBEAT, msgHeartBeat))
				}
			}
		}
	}

	atomic.StoreInt32(&pthis.isStopProcess, 1)
	close(pthis.chSrvProtoMsg)
	close(pthis.chInterMsg)
}

func (pthis*Server) ServerClose() {
	if pthis.connAcpt != nil {
		pthis.connAcpt.TcpConnectClose()
	}
	if pthis.connDial != nil {
		pthis.connDial.TcpConnectClose()
	}

	pthis.chSrvProtoMsg <- nil
}

func (pthis*Server) ServerCloseAndDelDialData() {
	pthis.srvMgr.tcpMgr.TcpDialDelDialData(pthis.srvID)
	pthis.ServerClose()
}

func (pthis*Server)processSrvReport(tcpAccept *tcp.TcpConnection){
	pthis.connAcpt = tcpAccept
	lin_common.LogDebug(pthis.srvID, " ", pthis)

	msgRes := &msgpacket.MSG_SRV_REPORT_RES{
		SrvId:pthis.srvID,
		TcpConnId:int64(tcpAccept.TcpConnectionID()),
	}
	TcpConnectSendProtoMsg(tcpAccept, msgpacket.MSG_TYPE__MSG_SRV_REPORT_RES, msgRes)
}

func (pthis*Server)processDailConnect(tcpDial *tcp.TcpConnection){
	pthis.connDial = tcpDial
	lin_common.LogDebug(pthis.srvID, " ", pthis)

	msgR := &msgpacket.MSG_SRV_REPORT{}
	msgR.SrvId = pthis.srvMgr.srvID
	tcpDial.TcpConnectSendBin(msgpacket.ProtoPacketToBin(msgpacket.MSG_TYPE__MSG_SRV_REPORT, msgR))
}

func (pthis*Server)processConnClose(tcpConn *tcp.TcpConnection){
	if tcpConn == nil {
		return
	}

	if pthis.connDial == nil && pthis.connAcpt == nil {
		lin_common.LogErr("connDial and connAcpt is nil:", pthis.srvID, " connection id:", tcpConn.TcpConnectionID())
	}
	bRedial := false
	if pthis.connAcpt != nil {
		if tcpConn.TcpConnectionID() == pthis.connAcpt.TcpConnectionID() {
			bRedial = true
		}
	}
	if pthis.connDial == nil {
		bRedial = true
	} else {
		if tcpConn.TcpConnectionID() == pthis.connDial.TcpConnectionID() {
			bRedial = true
		}
	}

	if !bRedial {
		lin_common.LogErr("not redial:", pthis.connAcpt, " connDial:", pthis.connDial, " tcpConn:", tcpConn)
		return
	}

	lin_common.LogDebug(pthis.srvID, " will redial", pthis)
	if pthis.connAcpt != nil {
		pthis.connAcpt.TcpConnectClose()
		pthis.connAcpt = nil
	}
	if pthis.connDial != nil {
		pthis.connDial.TcpConnectClose()
		pthis.connDial = nil
	}
	pthis.connDial, _ = pthis.srvMgr.tcpMgr.TcpDialMgrCheckReDial(tcpConn.SrvID)
}

func (pthis*Server)PushInterMsg(msg interface{}){
	if atomic.LoadInt32(&pthis.isStopProcess) == 1 {
		return
	}
	pthis.chInterMsg <- msg
}
func (pthis*Server)PushProtoMsg(msgType msgpacket.MSG_TYPE, protoMsg proto.Message, tcpConn *tcp.TcpConnection){
	if atomic.LoadInt32(&pthis.isStopProcess) == 1 {
		return
	}
	pthis.chSrvProtoMsg <- &interProtoMsg{
		msgType:msgType,
		protoMsg:protoMsg,
		tcpConn:tcpConn,
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

	atomic.AddInt64(&pthis.totalRPCPacket, 1)

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
	//lin_common.LogDebug(msg)
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

	pthis.connDial.TcpConnectSendBin(msgpacket.ProtoPacketToBin(msgpacket.MSG_TYPE__MSG_RPC, msgRPC))

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
		pthis.process_MSG_HEARTBEAT(interMsg.tcpConn, t)
	case *msgpacket.MSG_HEARTBEAT_RES:
		pthis.process_MSG_HEARTBEAT_RES(t)
	}
}

func (pthis*Server) process_MSG_HEARTBEAT (tcpConn *tcp.TcpConnection, protoMsg * msgpacket.MSG_HEARTBEAT) {
	lin_common.LogDebug(protoMsg)

	msgRes := &msgpacket.MSG_HEARTBEAT_RES{}
	msgRes.Id = protoMsg.Id
	TcpConnectSendProtoMsg(tcpConn, msgpacket.MSG_TYPE__MSG_HEARTBEAT_RES, msgRes)
}

func (pthis*Server) process_MSG_HEARTBEAT_RES (protoMsg * msgpacket.MSG_HEARTBEAT_RES) {
	lin_common.LogDebug(protoMsg)
}
