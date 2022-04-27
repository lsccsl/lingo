package main

import (
	"lin/lin_common"
	"lin/msgpacket"
	"runtime"
	"time"
)

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

	timestampLastHeartbeat int64
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

