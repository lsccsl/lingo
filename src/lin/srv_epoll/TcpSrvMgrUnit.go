package main

import (
	"lin/lin_common"
	"lin/msgpacket"
	"time"
)

type MAP_TCPSRV map[int64/*srv id*/]*TcpSrv

type TcpSrvMgrUnitStatic struct {
	totalRPCOut int64
	totalRPCIn int64
}
type TcpSrvMgrUnit struct {
	chSrv chan interface{}
	tcpSrvMgr *TcpSrvMgr
	mapSrv MAP_TCPSRV

	TcpSrvMgrUnitStatic
}


func (pthis*TcpSrvMgrUnit)_go_srvProcess_unit(){
	for {
		msg := <- pthis.chSrv
		switch t := msg.(type) {
		case *srvEvt_addremote:
			pthis.process_srvEvt_addremote(t)
		case *srvEvt_TcpDialSuc:
			pthis.process_srvEvt_DialSuc(t)
		case *srvEvt_TcpClose:
			pthis.process_srvEvt_TcpClose(t)
		case *srvEvt_TcpAcpt:
			pthis.process_srvEvt_TcpAcpt(t)
		case *srvEvt_protoMsg:
			pthis.process_srvEvt_protoMsg(t)
		case *srvEvt_RPC:
			pthis.process_srvEvt_RPC(t)
		case *srvEvt_RPC_Del:
			pthis.process_srvEvt_RPC_Del(t)
		case *srvEvt_static:
			pthis.process_srvEvt_static(t)
		}
	}
}

func (pthis*TcpSrvMgrUnit)addSrv(srv *TcpSrv){
	pthis.mapSrv[srv.srvID] = srv
}
func (pthis*TcpSrvMgrUnit)getSrv(srvID int64) *TcpSrv {
	srv, _ := pthis.mapSrv[srvID]
	return srv
}

func (pthis*TcpSrvMgrUnit)process_srvEvt_addremote(evt * srvEvt_addremote){
	lin_common.LogDebug("add srv:", evt.srvID, " addr:", evt.addr)
	oldSrv := pthis.getSrv(evt.srvID)
	bDial := true
	if oldSrv != nil {
		if oldSrv.addr == evt.addr {
			if !oldSrv.fdDial.IsNull() {
				bDial = false
			}
		}
	} else {
		pthis.addSrv(ConstructorTcpSrv(evt.srvID, evt.addr, pthis))
	}

	if bDial {
		fd, err := pthis.tcpSrvMgr.eSrvMgr.lsn.EPollListenerDial(evt.addr, &TcpAttachData{srvID : evt.srvID})
		if err != nil {
			lin_common.LogErr("connect to srv:", evt.srvID, " dial err")
		}
		lin_common.LogDebug("srv:", evt.srvID, " fd:", fd.String())
	}
}

func (pthis*TcpSrvMgrUnit)process_srvEvt_DialSuc(evt * srvEvt_TcpDialSuc){
	oldSrv := pthis.getSrv(evt.srvID)
	if oldSrv == nil {
		lin_common.LogDebug(" dial suc, no srv:", evt.srvID, " dial suc, fd:", evt.fdDial.String())
		return
	}
	lin_common.LogDebug("srv:", evt.srvID, " dial suc, fd:", evt.fdDial.String())
	oldFD := oldSrv.fdDial
	oldSrv.fdDial = evt.fdDial

	if !oldFD.IsSame(&evt.fdDial) && !oldFD.IsNull(){
		if oldFD.FD != evt.fdDial.FD {
			pthis.tcpSrvMgr.eSrvMgr.lsn.EPollListenerCloseTcp(oldFD)
		}
	}

	if oldSrv.timerDialClose != nil {
		oldSrv.timerDialClose.Reset(oldSrv.durationClose)
	}
}

func (pthis*TcpSrvMgrUnit)process_srvEvt_TcpAcpt(evt *srvEvt_TcpAcpt) {
	lin_common.LogDebug("srv:", evt.srvID, " acpt fd:", evt.fdAcpt.String())
	oldSrv := pthis.getSrv(evt.srvID)
	if oldSrv == nil {
		return
	}
	oldFD := oldSrv.fdAcpt
	oldSrv.fdAcpt = evt.fdAcpt

	if !oldFD.IsSame(&evt.fdAcpt) && !oldFD.IsNull(){
		if oldFD.FD != evt.fdAcpt.FD {
			pthis.tcpSrvMgr.eSrvMgr.lsn.EPollListenerCloseTcp(oldFD)
		}
	}

	if oldSrv.timerAcptClose != nil {
		oldSrv.timerAcptClose.Reset(oldSrv.durationClose)
	}

	msgRes := &msgpacket.MSG_SRV_REPORT_RES{
		SrvId:evt.srvID,
		TcpConnId:int64(oldSrv.fdAcpt.Magic),
	}
	pthis.tcpSrvMgr.eSrvMgr.SendProtoMsg(evt.fdAcpt, msgpacket.MSG_TYPE__MSG_SRV_REPORT_RES, msgRes)
}

func (pthis*TcpSrvMgrUnit)process_srvEvt_TcpClose(evt * srvEvt_TcpClose){
	lin_common.LogDebug("add srv:", evt.srvID, " fd:", evt.fd.String())
	oldSrv := pthis.getSrv(evt.srvID)
	if oldSrv == nil {
		return
	}

	if oldSrv.fdDial.IsSame(&evt.fd) || oldSrv.fdDial.IsNull() {
		oldSrv.fdDial = lin_common.FD_DEF_NIL

/*		fd, err := pthis.tcpSrvMgr.eSrvMgr.lsn.EPollListenerDial(oldSrv.addr, &TcpAttachData{srvID : oldSrv.srvID})
		if err != nil {
			lin_common.LogErr("srv:", oldSrv.srvID, " dial err")
		}
		lin_common.LogDebug(" reconnect to srv:", oldSrv.srvID, " fd:", fd.String())*/
	}
	if oldSrv.fdAcpt.IsSame(&evt.fd) {
		oldSrv.fdAcpt = lin_common.FD_DEF_NIL
	}
}

func (pthis*TcpSrvMgrUnit)process_srvEvt_protoMsg(evt *srvEvt_protoMsg) {
	//lin_common.LogDebug("proto msg srv:", evt.srvID, " fd:", evt.fd.String(), " msg:", evt.msg, " msgtype:", evt.msgType)
	oldSrv := pthis.getSrv(evt.srvID)
	if oldSrv == nil {
		return
	}

	if evt.fd.IsSame(&oldSrv.fdDial) {
		if oldSrv.timerDialClose != nil {
			oldSrv.timerDialClose.Reset(oldSrv.durationClose)
		}
	}
	if evt.fd.IsSame(&oldSrv.fdAcpt) {
		if oldSrv.timerAcptClose != nil {
			oldSrv.timerAcptClose.Reset(oldSrv.durationClose)
		}
	}

	switch evt.msgType {
	case msgpacket.MSG_TYPE__MSG_HEARTBEAT:
		t := evt.msg.(*msgpacket.MSG_HEARTBEAT)
		if t != nil {
			pthis.process_MSG_HEARTBEAT(oldSrv, evt.fd, t)
		}
	case msgpacket.MSG_TYPE__MSG_HEARTBEAT_RES:
		t := evt.msg.(*msgpacket.MSG_HEARTBEAT_RES)
		if t != nil {
			pthis.process_MSG_HEARTBEAT_RES(oldSrv, evt.fd, t)
		}
	case msgpacket.MSG_TYPE__MSG_RPC:
		t := evt.msg.(*msgpacket.MSG_RPC)
		if t != nil {
			pthis.process_MSG_RPC(oldSrv, evt.fd, t)
		}
	case msgpacket.MSG_TYPE__MSG_RPC_RES:
		t := evt.msg.(*msgpacket.MSG_RPC_RES)
		if t != nil {
			pthis.process_MSG_RPC_RES(oldSrv, evt.fd, t)
		}
	default:
		oldSrv.process_ProtoMsg(evt.fd, evt.msg)
	}
}

func (pthis*TcpSrvMgrUnit)process_srvEvt_RPC(evt *srvEvt_RPC) {
	lin_common.LogDebug("rpc msg srv:", evt.srvID, " rpcUUID:", evt.rpcUUID, " msgtype:", evt.msgType, " msg:", evt.msg)
	oldSrv := pthis.getSrv(evt.srvID)
	if oldSrv == nil {
		return
	}

	pthis.totalRPCOut ++
	oldSrv.TcpSrvSendRPC(evt)
}

func (pthis*TcpSrvMgrUnit)process_srvEvt_RPC_Del(evt *srvEvt_RPC_Del) {
	lin_common.LogDebug("del rpc msg srv:", evt.srvID, " rpcUUID:", evt.rpcUUID)
	oldSrv := pthis.getSrv(evt.srvID)
	if oldSrv == nil {
		return
	}

	oldSrv.TcpSrvDelRPC(evt.rpcUUID)
}
func (pthis*TcpSrvMgrUnit)process_srvEvt_static(evt *srvEvt_static) {
	evt.chBack <- *pthis
}

func (pthis*TcpSrvMgrUnit)process_MSG_HEARTBEAT(srv * TcpSrv, fd lin_common.FD_DEF, protoMsg *msgpacket.MSG_HEARTBEAT) {
	lin_common.LogDebug(" HB srv:", srv.srvID, " fd", fd.String(), " from srv:", protoMsg.Id)
	srv.timestampLastHeartbeat = time.Now().Unix()

	msgRes := &msgpacket.MSG_HEARTBEAT_RES{}
	msgRes.Id = protoMsg.Id
	pthis.tcpSrvMgr.eSrvMgr.SendProtoMsg(fd, msgpacket.MSG_TYPE__MSG_HEARTBEAT_RES, msgRes)
}

func (pthis*TcpSrvMgrUnit)process_MSG_HEARTBEAT_RES(srv * TcpSrv, fd lin_common.FD_DEF, protoMsg *msgpacket.MSG_HEARTBEAT_RES) {
	lin_common.LogDebug(" HBRES srv:", srv.srvID, " fd", fd.String(), " from srv:", protoMsg.Id)
	srv.timestampLastHeartbeat = time.Now().Unix()
}

func (pthis*TcpSrvMgrUnit)process_MSG_RPC(srv * TcpSrv, fd lin_common.FD_DEF, protoMsg *msgpacket.MSG_RPC) {
	pthis.totalRPCIn ++

	if srv == nil {
		return
	}
	srv.TcpSrvProcessRPCMsg(fd, protoMsg)
}

func (pthis*TcpSrvMgrUnit)process_MSG_RPC_RES(srv * TcpSrv, fd lin_common.FD_DEF, protoMsg *msgpacket.MSG_RPC_RES) {
	if srv == nil {
		return
	}
	srv.TcpSrvProcessRPCResMsg(fd, protoMsg)
}

func ConstructorTcpSrvMgrUnit(tcpSrvMgr*TcpSrvMgr)*TcpSrvMgrUnit {
	return &TcpSrvMgrUnit{
		chSrv : make(chan interface{}, 100),
		tcpSrvMgr : tcpSrvMgr,
		mapSrv : make(MAP_TCPSRV),
	}
}
