package main

import (
	"context"
	"github.com/golang/protobuf/proto"
	"lin/lin_common"
	msgpacket "lin/msgpacket"
	"lin/tcp"
	"sync/atomic"
	"time"
)

type ServerStatic struct {
	totalRPCIn int64
	totalRPCInLast int64

	totalRPCOut int64
	totalRPCOutLast int64
	totalRPCOutFail int64

	timestamp float64

	timestampLastHeartbeat int64
	timestampReport int64
}
type ServerDialData struct {
	dialTimeoutSec int
	closeExpireSec int
	ip string
	port int
	needRedial bool
	redialCount int

	tcpConnIDCurDial tcp.TCP_CONNECTION_ID

	DialCancelFunc context.CancelFunc
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

	ServerDialData
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
type interMsgBeginRedial struct {
	ip string
	port int
	closeExpireSec int
	dialTimeoutSec int
	needRedial bool
	redialCount int
}

func ConstructServer(srvMgr *ServerMgr, srvID int64, heartbeatIntervalSec int)*Server {
	s := &Server{
		srvMgr:srvMgr,
		srvID:srvID,
		connDial:nil,
		connAcpt:nil,
		chSrvProtoMsg:make(chan *interProtoMsg, 100),
		chInterMsg:make(chan interface{}, 100),
		heartbeatIntervalSec:heartbeatIntervalSec,
		rpcMgr:ConstructRPCManager(),
		isStopProcess:0,
	}
	srvMgr.addServer(s)
	go s.go_serverProcess()

	return s
}

func (pthis*Server) ServerSetDialData(ip string, port int, closeExpireSec int,
	dialTimeoutSec int,
	needRedial bool, redialCount int) {
	pthis.chInterMsg <- &interMsgBeginRedial{
		ip:ip,
		port:port,
		closeExpireSec:closeExpireSec,
		dialTimeoutSec:dialTimeoutSec,
		needRedial:needRedial,
		redialCount:redialCount,
	}
}

func (pthis*Server) go_serverProcess() {
	defer func() {
		lin_common.LogDebug("srv:", pthis.srvID, " exit process")
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
				case *interMsgBeginRedial:
					pthis.processRedial(t)
				}
			}

		case <-chTimer:
			{
				chTimer = time.After(time.Second * time.Duration(pthis.heartbeatIntervalSec))
				//send heartbeat
				if pthis.connDial != nil {
/*					lin_common.LogDebug("send heartbeat from dial, srv:", pthis.srvID, pthis.heartbeatIntervalSec,
						" connection id:", pthis.connDial.TcpConnectionID())*/
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

/*func (pthis*Server) ServerCloseAndDelDialData() {
	pthis.srvMgr.tcpMgr.TcpDialDelDialData(pthis.srvID)
	pthis.ServerClose()
}*/

func (pthis*Server)processSrvReport(tcpAccept *tcp.TcpConnection){
	if tcpAccept == nil {
		return
	}
	if pthis.connAcpt != nil {
		if pthis.connAcpt.TcpConnectionID() != tcpAccept.TcpConnectionID() {
			pthis.connAcpt.TcpConnectSetCloseReason(tcp.TCP_CONNECTION_CLOSE_REASON_new_acpt)
			pthis.connAcpt.TcpConnectClose()
		}
	}
	pthis.connAcpt = tcpAccept

	lin_common.LogDebug("srv:", pthis.srvID, " conn:", pthis.connAcpt.TcpConnectionID())
	pthis.timestampReport = time.Now().Unix()

	msgRes := &msgpacket.MSG_SRV_REPORT_RES{
		SrvId:pthis.srvID,
		TcpConnId:int64(pthis.connAcpt.TcpConnectionID()),
	}
	TcpConnectSendProtoMsg(pthis.connAcpt, msgpacket.MSG_TYPE__MSG_SRV_REPORT_RES, msgRes)
}

func (pthis*Server)processDailConnect(tcpDial *tcp.TcpConnection){
	if tcpDial == nil {
		return
	}
	if pthis.connDial != nil {
		if pthis.connDial.TcpConnectionID() != tcpDial.TcpConnectionID() {
			pthis.connDial.TcpConnectSetCloseReason(tcp.TCP_CONNECTION_CLOSE_REASON_new_dial)
			pthis.connDial.TcpConnectClose()
		}
	}
	pthis.connDial = tcpDial

	lin_common.LogDebug(" srv:", pthis.srvID, " conn:", pthis.connDial.TcpConnectionID())
	pthis.srvMgr.tcpMgr.TcpDialDelDialData(pthis.srvID)

	msgR := &msgpacket.MSG_SRV_REPORT{}
	msgR.SrvId = pthis.srvMgr.srvID
	msgR.TcpConnId = int64(pthis.connDial.TcpConnectionID())
	pthis.connDial.TcpConnectSendBin(msgpacket.ProtoPacketToBin(msgpacket.MSG_TYPE__MSG_SRV_REPORT, msgR))
}

func (pthis*Server)processConnClose(tcpConn *tcp.TcpConnection){
	if tcpConn == nil {
		return
	}

	if tcpConn.IsAccept {
		pthis.connAcpt = nil
		return
	}

	if pthis.connDial != nil {
		if tcpConn.TcpConnectionID() != pthis.connDial.TcpConnectionID() {
			return
		}
	}

	lin_common.LogDebug("srv:", pthis.srvID, " will redial", " tcpConn:", tcpConn.TcpConnectionID(), " addr:", pthis.ip, ":", pthis.port)
	pthis.connDial = nil
	pthis.srvMgr.tcpMgr.TcpDialMgrDial(pthis.srvID, pthis.ip, pthis.port,
		pthis.closeExpireSec, pthis.dialTimeoutSec, pthis.needRedial, pthis.redialCount,
		pthis.srvMgr.dialPool)
}

func (pthis*Server)processRedial(dialMsg *interMsgBeginRedial){
	bRedial := false
	if pthis.ip != dialMsg.ip || pthis.port != dialMsg.port {
		bRedial = true
	}
	lin_common.LogDebug("srv:", pthis.srvID, " ", dialMsg.ip, ":", dialMsg.port, " ", pthis.ip, ":", pthis.port, " bRedial:", bRedial)

	pthis.dialTimeoutSec = dialMsg.dialTimeoutSec
	pthis.closeExpireSec = dialMsg.closeExpireSec
	pthis.ip = dialMsg.ip
	pthis.port = dialMsg.port
	pthis.needRedial = dialMsg.needRedial
	pthis.redialCount = dialMsg.redialCount

	if bRedial {
		if pthis.connDial != nil {
			pthis.connDial.TcpConnectClose()
		}
		srvMgr.tcpMgr.TcpDialMgrDial(pthis.srvID, pthis.ip, pthis.port, pthis.closeExpireSec,
			pthis.dialTimeoutSec,
			pthis.needRedial, pthis.redialCount,
			pthis.srvMgr.dialPool)
	}
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
	msgRes := pthis.ServerProcessRPC(tcpConn, msgBody)

	msgRPCRes := &msgpacket.MSG_RPC_RES{
		MsgId:msg.MsgId,
		MsgType:msg.MsgType,
		ResCode:msgpacket.RESPONSE_CODE_RESPONSE_CODE_OK,
		Timestamp:msg.Timestamp,
		TimestampArrive: msg.TimestampArrive,
		TimestampProcess: time.Now().UnixMilli(),
	}

	tDiff := msgRPCRes.TimestampProcess - msgRPCRes.TimestampArrive
	if tDiff > msg.TimeoutWait {
		lin_common.LogErr("rpc timeout, tDiff:", tDiff, " timeout wait:", msg.TimeoutWait,
			" srv:", tcpConn.SrvID, " conn:", tcpConn.TcpConnectionID())
	}

	if msgRes != nil {
		var err error
		msgRPCRes.MsgBin, err = proto.Marshal(msgRes)
		if err != nil {
			lin_common.LogErr(err)
		}
	}

	atomic.AddInt64(&pthis.totalRPCIn, 1)
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
		lin_common.LogErr("fail find rpc:", msg.MsgId, " srv:", pthis.srvID)
		return
	}
	if rreq.chNtf != nil {
		rreq.chNtf <- msgBody
	}
}



// SendRPC_Async @brief will block timeoutMilliSec
func (pthis*Server)SendRPC_Async(msgType msgpacket.MSG_TYPE, protoMsg proto.Message, timeoutMilliSec int) (proto.Message, error) {

	defer func() {
		err := recover()
		if err != nil {
			lin_common.LogErr(err)
		}
	}()

	if pthis.rpcMgr == nil{
		return nil, lin_common.GenErr(lin_common.ERR_sys, "no rpc mgr")
	}

	msgRPC := &msgpacket.MSG_RPC{
		MsgId:lin_common.GenUUID64_V4(),
		MsgType:int32(msgType),
	}
	var err error
	msgRPC.MsgBin, err = proto.Marshal(protoMsg)
	if err != nil {
		lin_common.LogErr(err)
		return nil, lin_common.GenErr(lin_common.ERR_sys, "packet err")
	}

	rreq := pthis.rpcMgr.RPCManagerAddReq(msgRPC.MsgId)

	atomic.AddInt64(&pthis.totalRPCOut, 1)
	pthis.connDial.TcpConnectSendBin(msgpacket.ProtoPacketToBin(msgpacket.MSG_TYPE__MSG_RPC, msgRPC))

	var res proto.Message = nil
	select{
	case resCh := <-rreq.chNtf:
		res, _ = resCh.(proto.Message)
	case <-time.After(time.Millisecond * time.Duration(timeoutMilliSec)):
		atomic.AddInt64(&pthis.totalRPCOutFail, 1)
		err = lin_common.GenErr(lin_common.ERR_rpc_timeout, " rpc time out srv:", pthis.srvID, " rpcid:", msgRPC.MsgId)
	}

	pthis.rpcMgr.RPCManagerDelReq(msgRPC.MsgId)

	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, lin_common.GenErr(lin_common.ERR_sys, "msg is nil")
	}

	return res, nil
}

func (pthis*Server)processServerMsg (interMsg * interProtoMsg){
	switch t:=interMsg.protoMsg.(type){
	case *msgpacket.MSG_HEARTBEAT:
		pthis.process_MSG_HEARTBEAT(interMsg.tcpConn, t)
	case *msgpacket.MSG_HEARTBEAT_RES:
		pthis.process_MSG_HEARTBEAT_RES(interMsg.tcpConn, t)
	default:
		pthis.processOtherServerMsg(interMsg)
	}
}

func (pthis*Server) process_MSG_HEARTBEAT (tcpConn *tcp.TcpConnection, protoMsg * msgpacket.MSG_HEARTBEAT) {
	lin_common.LogDebug(" srv:", pthis.srvID, " conn", tcpConn.TcpConnectionID(), " from srv:", protoMsg.Id)
	if tcpConn != nil {
		msgRes := &msgpacket.MSG_HEARTBEAT_RES{}
		msgRes.Id = protoMsg.Id
		TcpConnectSendProtoMsg(tcpConn, msgpacket.MSG_TYPE__MSG_HEARTBEAT_RES, msgRes)
	}
	pthis.timestampLastHeartbeat = time.Now().Unix()
}

func (pthis*Server) process_MSG_HEARTBEAT_RES (tcpConn *tcp.TcpConnection, protoMsg * msgpacket.MSG_HEARTBEAT_RES) {
	lin_common.LogDebug(" srv:", pthis.srvID, " conn", tcpConn.TcpConnectionID(), " from srv:", protoMsg.Id)
	pthis.timestampLastHeartbeat = time.Now().Unix()
}
