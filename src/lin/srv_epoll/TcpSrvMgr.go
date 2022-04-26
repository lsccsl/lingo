package main

import (
	"github.com/golang/protobuf/proto"
	"lin/lin_common"
)

/* begin srv event */
type srvEvt_addremote struct {
	srvID int64
	addr string
	closeExpireSec int
}
type srvEvt_TcpDialSuc struct {
	srvID int64
	fdDial lin_common.FD_DEF
}
type srvEvt_TcpClose struct {
	srvID int64
	fd lin_common.FD_DEF
}
type srvEvt_TcpAcpt struct {
	srvID int64
	fdAcpt lin_common.FD_DEF
}
type srvEvt_protoMsg struct {
	srvID int64
	fd lin_common.FD_DEF
	msg proto.Message
}
/* end srv event */

type TcpSrvMgr struct {
	eSrvMgr *EpollServerMgr

	mgrUnit []*TcpSrvMgrUnit
}


func (pthis*TcpSrvMgr)getSrvProcessUnit(srvID int64)*TcpSrvMgrUnit{
	processUnitCount := int64(len(pthis.mgrUnit))
	idx := srvID % processUnitCount
	if idx >= processUnitCount {
		return nil
	}
	pu := pthis.mgrUnit[idx]
	if pu == nil {
		return nil
	}
	return pu
}

func (pthis*TcpSrvMgr)TcpSrvMgrPushMsgToUnit(srvID int64, msg interface{}) {
	pu := pthis.getSrvProcessUnit(srvID)
	if pu != nil {
		pu.chSrv <- msg
	} else {
		lin_common.LogErr("srv:", srvID, " push msg to srv process unit err")
	}
}

func (pthis*TcpSrvMgr)TcpSrvMgrAddRemoteSrv(srvID int64, addr string, closeExpireSec int){
	pthis.TcpSrvMgrPushMsgToUnit(srvID, &srvEvt_addremote{
		srvID : srvID,
		addr : addr,
		closeExpireSec : closeExpireSec,
	})

	fd, err := pthis.eSrvMgr.lsn.EPollListenerDial(addr, &TcpAttachData{srvID : srvID})
	if err != nil {
		lin_common.LogErr("connect to srv:", srvID, " dial err")
	}
	lin_common.LogDebug("srv:", srvID, " fd:", fd.String())
}

func (pthis*TcpSrvMgr)TcpSrvMgrRPC(srvID int64){

}

func ConstructorTcpSrvMgr(eSrvMgr *EpollServerMgr, srvProcessUnitCount int) *TcpSrvMgr {
	tcpSrvMgr := &TcpSrvMgr{
		eSrvMgr : eSrvMgr,
	}

/*	for i := 0; i < srvProcessUnitCount; i ++ {
		processUnit := ConstructorTcpSrvMgrUnit(tcpSrvMgr)
		tcpSrvMgr.mgrUnit = append(tcpSrvMgr.mgrUnit, processUnit)
		processUnit._go_Process_unit()
	}*/

	return tcpSrvMgr
}

