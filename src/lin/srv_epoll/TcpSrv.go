package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"lin/lin_common"
	cor_pool "lin/lin_cor_pool"
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

func (pthis*TcpSrv)TcpSrvProcessRPCResMsg(fd lin_common.FD_DEF, msgRPCRes *msgpacket.MSG_RPC_RES){
	defer func() {
		err := recover()
		if err != nil {
			lin_common.LogErr(" rpc res err:", fd.String(), " err:", err, " rpc:", msgRPCRes.MsgId, " rpctype:", msgRPCRes.MsgType)
		}
	}()

	rreq := pthis.getRPC(msgRPCRes.MsgId)
	pthis.delRPC(msgRPCRes.MsgId)

	if rreq == nil {
		return
	}
	msgBody := msgpacket.ParseProtoMsg(msgRPCRes.MsgBin, msgRPCRes.MsgType)
	rreq.chRouteBack <- msgBody
}

func (pthis*TcpSrv)TcpSrvProcessRPCMsg(fd lin_common.FD_DEF, msgRPC *msgpacket.MSG_RPC){
	msgRPC.TimestampArrive = time.Now().UnixMilli()

	tDiff := msgRPC.TimestampArrive - msgRPC.Timestamp
	if tDiff > msgRPC.TimeoutWait {
		lin_common.LogDebug("recv rpc timeout, tDiff:", tDiff, " timeout wait:", msgRPC.TimeoutWait,
			" srv:", pthis.srvID, " fd:", fd.String())
	}

	rpcFunc := func() {
		msgRPCRes := &msgpacket.MSG_RPC_RES{
			MsgId:msgRPC.MsgId,
			MsgType:msgRPC.MsgType,
			ResCode:msgpacket.RESPONSE_CODE_RESPONSE_CODE_OK,
			Timestamp:msgRPC.Timestamp,
			TimestampArrive: msgRPC.TimestampArrive,
			TimestampProcess: time.Now().UnixMilli(),
		}

		msgRPCBody := msgpacket.ParseProtoMsg(msgRPC.MsgBin, msgRPC.MsgType)
		var msgResBody proto.Message = nil
		switch t:= msgRPCBody.(type) {
		case *msgpacket.MSG_TEST:
			{
				msgResBody = &msgpacket.MSG_TEST_RES{Id: t.Id, Str:t.Str, Seq: t.Seq}
			}
		}

		if msgResBody != nil {
			var err error
			msgRPCRes.MsgBin, err = proto.Marshal(msgResBody)
			if err != nil {
				lin_common.LogErr(err)
			}
		}

		pthis.pu.tcpSrvMgr.eSrvMgr.SendProtoMsg(fd, msgpacket.MSG_TYPE__MSG_RPC_RES, msgRPCRes)
	}
	pthis.pu.tcpSrvMgr.rpcPool.CorPoolAddJob(&cor_pool.CorPoolJobData{
		JobType_ : EN_CORPOOL_JOBTYPE_Rpc_req,
		JobData_ : pthis.srvID,
		JobCB_ : func(jd cor_pool.CorPoolJobData){
			rpcFunc()
		},
	})
}

func (pthis*TcpSrv)process_ProtoMsg(fd lin_common.FD_DEF, protoMsg proto.Message) {
	switch t:=protoMsg.(type){
	case *msgpacket.MSG_TEST_RPC:
		fmt.Println(t)
	}
}

func (pthis*TcpSrv)process_Timer(evt *srvEvt_timer) {
	switch evt.timerType {
	case EN_TIMER_TYPE_close_dial:
		{
			if !pthis.fdDial.IsNull() {
				lin_common.LogDebug("timeout close srv dial:", pthis.srvID, " fdDial:", pthis.fdDial.String())
				pthis.pu.tcpSrvMgr.eSrvMgr.lsn.EPollListenerCloseTcp(pthis.fdDial)
			}
			pthis.timerDialClose = time.AfterFunc(pthis.durationClose, pthis.startCloseDailTimer)
		}
	case EN_TIMER_TYPE_close_acpt:
		{
			if !pthis.fdAcpt.IsNull() {
				lin_common.LogDebug("timeout close srv acpt:", pthis.srvID, " fdAcpt:", pthis.fdAcpt.String())
				pthis.pu.tcpSrvMgr.eSrvMgr.lsn.EPollListenerCloseTcp(pthis.fdAcpt)
			}
			pthis.timerAcptClose = time.AfterFunc(pthis.durationClose, pthis.startCloseAcptTimer)
		}
	case EN_TIMER_TYPE_heartbeat:
		{
			//lin_common.LogDebug("send heartbeat to dial, srv:", pthis.srvID, " fdDial:", pthis.fdDial.String())
			msgHeartBeat := &msgpacket.MSG_HEARTBEAT{}
			msgHeartBeat.Id = pthis.srvID

			if !pthis.fdDial.IsNull() {
				pthis.pu.tcpSrvMgr.eSrvMgr.SendProtoMsg(pthis.fdDial, msgpacket.MSG_TYPE__MSG_HEARTBEAT, msgHeartBeat)
			}
			pthis.timerHB = time.AfterFunc(pthis.durationHB, pthis.startHeartBeatTimer)
		}
	}
}

func (pthis*TcpSrv)startHeartBeatTimer(){
	pthis.pu.chSrv <- &srvEvt_timer{srvID:pthis.srvID,
		timerType:EN_TIMER_TYPE_heartbeat,
		timerData:nil,
	}
}

func (pthis*TcpSrv)startCloseAcptTimer() {
	pthis.pu.chSrv <- &srvEvt_timer{srvID:pthis.srvID,
		timerType:EN_TIMER_TYPE_close_acpt,
		timerData:nil,
	}
}

func (pthis*TcpSrv)startCloseDailTimer() {
	pthis.pu.chSrv <- &srvEvt_timer{srvID:pthis.srvID,
		timerType:EN_TIMER_TYPE_close_dial,
		timerData:nil,
	}
}


func ConstructorTcpSrv(srvID int64, addr string, pu *TcpSrvMgrUnit) *TcpSrv {
	timeSec := pu.tcpSrvMgr.eSrvMgr.clientCloseTimeoutSec
	if timeSec < 6 {
		timeSec = 6
	}
	srv := &TcpSrv{
		fdDial: lin_common.FD_DEF_NIL,
		fdAcpt: lin_common.FD_DEF_NIL,
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

	lin_common.LogDebug(" srv:", srvID, " close timeout:", srv.durationClose)

	srv.timerAcptClose = time.AfterFunc(srv.durationClose, srv.startCloseAcptTimer)
	srv.timerDialClose = time.AfterFunc(srv.durationClose, srv.startCloseDailTimer)
	srv.timerHB        = time.AfterFunc(srv.durationHB,    srv.startHeartBeatTimer)

	return srv
}

