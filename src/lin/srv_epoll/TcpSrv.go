package main

import (
	"github.com/golang/protobuf/proto"
	"lin/lin_common"
	"lin/msgpacket"
	"runtime"
	"time"
)

type CHAN_RPC_ROUTEBACK chan proto.Message
type RPCReq struct {
	rpcID int64
	chRouteBack CHAN_RPC_ROUTEBACK
}
type MAP_RPC_REQ map[int64/* rpc msg id */]*RPCReq

type TcpSrvStatic struct {
	timestampLastHeartbeat int64
}
type TcpSrvRPC struct {
	mapRPC MAP_RPC_REQ
}
type TcpSrv struct {
	srvID int64
	addr string
	fdDial lin_common.FD_DEF
	fdAcpt lin_common.FD_DEF

	timerDialClose * time.Timer
	timerAcptClose * time.Timer
	durationClose time.Duration
	timerHB * time.Timer
	durationHB time.Duration
	pu *TcpSrvMgrUnit

	TcpSrvRPC

	TcpSrvStatic
}

func (pthis*TcpSrv)Destructor() {
	lin_common.LogDebug(" srv:", pthis.srvID, " fdDial:", pthis.fdDial.String(), " fdAcpt:", pthis.fdAcpt.String())
	runtime.SetFinalizer(pthis, nil)
	if pthis.timerDialClose != nil {
		pthis.timerDialClose.Stop()
		pthis.timerDialClose = nil
	}
	if pthis.timerAcptClose != nil {
		pthis.timerAcptClose.Stop()
		pthis.timerAcptClose = nil
	}
	if pthis.timerHB != nil {
		pthis.timerHB.Stop()
		pthis.timerHB = nil
	}
}

func (pthis*TcpSrv)sendHeartBeat(){
	lin_common.LogDebug("send heartbeat to dial:", pthis.srvID, " fdDial:", pthis.fdDial.String())
	msgHeartBeat := &msgpacket.MSG_HEARTBEAT{}
	msgHeartBeat.Id = pthis.srvID

	pthis.pu.tcpSrvMgr.eSrvMgr.SendProtoMsg(pthis.fdDial, msgpacket.MSG_TYPE__MSG_HEARTBEAT, msgHeartBeat)
	pthis.timerHB = time.AfterFunc(pthis.durationHB,
		func(){
			pthis.sendHeartBeat()
		})
}

func (pthis*TcpSrv)addRPC(rpcID int64, chRouteBack CHAN_RPC_ROUTEBACK) {
	pthis.mapRPC[rpcID] = &RPCReq{rpcID:rpcID,chRouteBack: chRouteBack}
}
func (pthis*TcpSrv)getRPC(rpcID int64) *RPCReq {
	return pthis.mapRPC[rpcID]
}
func (pthis*TcpSrv)delRPC(rpcID int64){
	delete(pthis.mapRPC, rpcID)
}


func (pthis*TcpSrv)TcpSrvSendRPC(evt *srvEvt_RPC){
	msgRPC := &msgpacket.MSG_RPC{
		MsgId:evt.rpcUUID,
		MsgType:int32(evt.msgType),
		Timestamp:time.Now().UnixMilli(),
		TimeoutWait:evt.timeoutMills,
	}
	var err error
	msgRPC.MsgBin, err = proto.Marshal(evt.msg)
	if err != nil {
		lin_common.LogErr(err)
		evt.chRouteBack <- nil
		return
	}

	pthis.addRPC(msgRPC.MsgId, evt.chRouteBack)
	pthis.pu.tcpSrvMgr.eSrvMgr.SendProtoMsg(pthis.fdDial, msgpacket.MSG_TYPE__MSG_RPC, msgRPC)
}

func (pthis*TcpSrv)TcpSrvDelRPC(rpcUUID int64){
	pthis.delRPC(rpcUUID)
}

func (pthis*TcpSrv)TcpSrvProcessRPCMsg(fd lin_common.FD_DEF, protoMsg *msgpacket.MSG_RPC){
	rreq := pthis.getRPC(protoMsg.MsgId)
	if rreq == nil {
		lin_common.LogDebug(" can't find rpc:", protoMsg.MsgId, " srv:", pthis.srvID)
	}


	// todo : send rpc res
}

func ConstructorTcpSrv(srvID int64, addr string, pu *TcpSrvMgrUnit) *TcpSrv {
	timeSec := pu.tcpSrvMgr.eSrvMgr.clientCloseTimeoutSec
	if timeSec < 6 {
		timeSec = 6
	}
	srv := &TcpSrv{
		srvID : srvID,
		pu : pu,
		addr: addr,
		durationClose : time.Second*time.Duration(timeSec),
		durationHB : time.Second*time.Duration(timeSec / 2),
		TcpSrvRPC:TcpSrvRPC{
			mapRPC : make(MAP_RPC_REQ),
		},
	}
	runtime.SetFinalizer(srv, (*TcpSrv).Destructor)

	srv.timerDialClose = time.AfterFunc(srv.durationClose,
		func(){
			lin_common.LogDebug("timeout close srv dial:", srv.srvID, " fdDial:", srv.fdDial.String())
			srv.pu.tcpSrvMgr.eSrvMgr.lsn.EPollListenerCloseTcp(srv.fdDial)
		})
	srv.timerAcptClose = time.AfterFunc(srv.durationClose,
		func(){
			lin_common.LogDebug("timeout close srv acpt:", srv.srvID, " fdAcpt:", srv.fdAcpt.String())
			srv.pu.tcpSrvMgr.eSrvMgr.lsn.EPollListenerCloseTcp(srv.fdAcpt)
		})
	srv.timerHB = time.AfterFunc(srv.durationHB,
		func(){
			srv.sendHeartBeat()
		})
	return srv
}

