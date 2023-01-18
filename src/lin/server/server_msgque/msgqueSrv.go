package main

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"lin/lin_common"
	"lin/msgpacket"
	"lin/server/server_common"
	"net"
	"strconv"
	"time"
)

type MsgQueSrv struct {
	lsn *lin_common.EPollListener

	msgqueCenterAddr string

	fdCenter lin_common.FD_DEF
	addrOut string
	queSrvID server_common.SRV_ID

	timerReconnMsgQueCenter time.Timer // timer reconnect to msg que server

	otherMgr otherMsgQueSrvMgr
	smgr *SrvMgr
}

type otherMsgQueSrvInfo struct {
	fdDial lin_common.FD_DEF
	fdAccept lin_common.FD_DEF
	ip string
	port int32
	queSrvID server_common.SRV_ID
}

func (pthis *otherMsgQueSrvInfo)String() (str string) {
	str += "que srv id:" + pthis.queSrvID.String() +
		"[" + pthis.ip + ":" + strconv.FormatInt(int64(pthis.port), 10) + "]" +
		" fdDial:" + pthis.fdDial.String() +
		" fdConn:" + pthis.fdAccept.String()

	return
}

// tcpAttachDataMsgQueSrvDial tcp dial to other msg que
type tcpAttachDataMsgQueSrvDial struct{
	queSrvID server_common.SRV_ID
}
// tcpAttachDataMsgQueSrvAccept tcp accept from other msg que
type tcpAttachDataMsgQueSrvAccept struct {
	queSrvID server_common.SRV_ID
}
// tcpAttachDataMsgQueCenter tcp dial to msg que center
type tcpAttachDataMsgQueCenterDial struct {
}
// tcpAttachDataSrv tcp accept from srv
type tcpAttachDataSrvAccept struct {
	srvUUID server_common.SRV_ID
}

func (pthis*MsgQueSrv)TcpAcceptConnection(fd lin_common.FD_DEF, addr net.Addr, inAttachData interface{})(outAttachData interface{}) {
	lin_common.LogDebug(fd, addr, inAttachData)
	return nil
}

func (pthis*MsgQueSrv)TcpDialConnection(fd lin_common.FD_DEF, addr net.Addr, inAttachData interface{})(outAttachData interface{}) {

	switch inAttachData.(type) {
	case *tcpAttachDataMsgQueCenterDial: // dial to msg que center tcp connection ok
		{
			if !fd.IsSame(&pthis.fdCenter) {
				pthis.lsn.EPollListenerCloseTcp(fd, server_common.EN_TCP_CLOSE_REASON_repeated_msgque_center)
				return
			}
			// send reg msg to msq que center
			tcpAddr, err := net.ResolveTCPAddr("tcp", pthis.addrOut)
			if err != nil {
				lin_common.LogErr(err)
			}
			pbMsgReg := &msgpacket.PB_MSG_INTER_QUESRV_REGISTER{
				QueSrvId: int64(pthis.queSrvID),
				Ip: tcpAddr.IP.String(),
				Port: int32(tcpAddr.Port),
			}
			pthis.SendProtoMsg(fd, msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_REGISTER, pbMsgReg)
		}

	case *tcpAttachDataMsgQueSrvDial: //  dial to other msg que srv tcp connection ok
		{
			// send conn msg to other msg que
			pbMsgConn := &msgpacket.PB_MSG_INTER_QUESRV_CONNECT{
				QueSrvId:int64(pthis.queSrvID),
				AllSrv:&msgpacket.PB_SRV_INFO_ALL{},
			}
			pthis.smgr.getAllSrvNetPB(pbMsgConn.AllSrv)
			pthis.SendProtoMsg(fd, msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_CONNECT, pbMsgConn)
		}

	default:
		{
			lin_common.LogDebug(fd, addr, inAttachData)
		}
	}

	return
}

func (pthis*MsgQueSrv)TcpClose(fd lin_common.FD_DEF, closeReason lin_common.EN_TCP_CLOSE_REASON, inAttachData interface{}) {
	lin_common.LogInfo(fd, " attach data:", inAttachData, " closeReason:", closeReason)

	switch t := inAttachData.(type) {
	case *tcpAttachDataMsgQueSrvDial:
		{
			queSrvID := t.queSrvID
			lin_common.LogInfo(fd, "dial fd close, attach data:", inAttachData, " closeReason:", closeReason, " ", queSrvID.String())

			qsi := &otherMsgQueSrvInfo{}
			ok := pthis.otherMgr.Load(queSrvID, qsi)
			if !ok {
				lin_common.LogInfo(queSrvID.String(), "receive dial tcp close, not exist")
				return
			}
			if !qsi.fdDial.IsSame(&fd) {
				lin_common.LogInfo(queSrvID.String(), "receive dial tcp close, fd is not the same", "now:", qsi.fdAccept, "close:", fd)
				return
			}
			pthis.deleteMsgQueSrvAndRedia(queSrvID)
		}

	case *tcpAttachDataMsgQueSrvAccept:
		{
			queSrvID := t.queSrvID
			lin_common.LogInfo(fd, "accept fd close, attach data:", inAttachData, " closeReason:", closeReason, " ", queSrvID.String())

			qsi := &otherMsgQueSrvInfo{}
			ok := pthis.otherMgr.Load(queSrvID, qsi)
			if !ok {
				lin_common.LogInfo(queSrvID.String(), "receive accept tcp close, not exist")
				return
			}
			if !qsi.fdAccept.IsSame(&fd) {
				lin_common.LogInfo(queSrvID.String(), "receive accept tcp close, fd is not the same", "now:", qsi.fdAccept, "close:", fd)
				return
			}

			pthis.otherMgr.updateQueSrvAccept(queSrvID, lin_common.FD_DEF_NIL)
		}

	case *tcpAttachDataMsgQueCenterDial:
		{
			lin_common.LogInfo(fd, "dial msg que center close, attach data:", inAttachData, " closeReason:", closeReason)
			if !fd.IsSame(&pthis.fdCenter) {
				return
			}

			time.AfterFunc(time.Second * 3, func() {
				pthis.fdCenter, _ = pthis.lsn.EPollListenerDial(pthis.msgqueCenterAddr, &tcpAttachDataMsgQueCenterDial{}, false)
			})
		}

	case *tcpAttachDataSrvAccept:
		{
			// report to all other que srv
			pthis.process_TcpClose_SrvReg(fd, inAttachData)
		}
	}
}

func (pthis*MsgQueSrv)TcpOutBandData(fd lin_common.FD_DEF, data interface{}, inAttachData interface{}) {
	lin_common.LogDebug(fd, data, inAttachData)
}

func (pthis*MsgQueSrv)TcpTick(fd lin_common.FD_DEF, tNowMill int64, inAttachData interface{}){

	switch t := inAttachData.(type) {
	case *tcpAttachDataMsgQueCenterDial:
		{
			lin_common.LogDebug("send heart beat to msg que center ", pthis.queSrvID.String())
			// send heartbeat
			pbHB := &msgpacket.PB_MSG_INTER_QUECENTER_HEARTBEAT{
				QueSrvId: int64(pthis.queSrvID),
			}
			pthis.SendProtoMsg(pthis.fdCenter, msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUECENTER_HEARTBEAT, pbHB)
		}

	case *tcpAttachDataMsgQueSrvDial:
		{
			lin_common.LogDebug("send heart beat to other que srv ", t.queSrvID.String())
			pbHB := &msgpacket.PB_MSG_INTER_QUESRV_HEARTBEAT{
				QueSrvId: int64(pthis.queSrvID),
			}
			pthis.SendProtoMsg(fd, msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_HEARTBEAT, pbHB)
		}

	default:
		{
			lin_common.LogDebug(fd, " tNowMill:", tNowMill, " inAttachData:", inAttachData)
		}
	}
}

func (pthis*MsgQueSrv)TcpData(fd lin_common.FD_DEF, readBuf *bytes.Buffer, inAttachData interface{})(bytesProcess int, outAttachData interface{}) {
	packType, bytesProcess, protoMsg := msgpacket.ProtoUnPacketFromBin(readBuf)
	if protoMsg == nil {
		return
	}

	switch msgpacket.PB_MSG_INTER_TYPE(packType) {
	case msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_REGISTER_RES:
		{
			lin_common.LogDebug(fd, "PB_MSG_INTER_QUESRV_REGISTER_RES:", protoMsg, " attach:", inAttachData)
			pthis.process_PB_MSG_INTER_QUESRV_REGISTER_RES(protoMsg)
		}

	case msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_ONLINE_NTF:
		{
			lin_common.LogDebug(fd, "PB_MSG_INTER_QUESRV_ONLINE_NTF:", protoMsg, " attach:", inAttachData)
			pthis.process_PB_MSG_INTER_QUESRV_ONLINE_NTF(protoMsg)
		}

	case msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_OFFLINE_NTF:
		{
			lin_common.LogDebug(fd, "PB_MSG_INTER_QUESRV_OFFLINE_NTF:", protoMsg, " attach:", inAttachData)
			pthis.process_PB_MSG_INTER_QUESRV_OFFLINE_NTF(protoMsg)
		}

	case msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_CONNECT:
		{
			outAttachData = pthis.process_PB_MSG_INTER_QUESRV_CONNECT(fd, protoMsg, inAttachData)
		}
	case msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_CONNECT_RES:
		{
			pthis.process_PB_MSG_INTER_QUESRV_CONNECT_RES(fd, protoMsg, inAttachData)
		}

	case msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUECENTER_HEARTBEAT_RES:
		{
			lin_common.LogDebug(fd, "PB_MSG_INTER_QUECENTER_HEARTBEAT_RES:", protoMsg, " attach:", inAttachData)
		}

	case msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_HEARTBEAT:
		{
			pthis.process_PB_MSG_INTER_QUESRV_HEARTBEAT(fd, protoMsg, inAttachData)
		}

	case msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_HEARTBEAT_RES:
		{
			lin_common.LogDebug(fd, "PB_MSG_INTER_QUESRV_HEARTBEAT_RES:", protoMsg, " attach:", inAttachData)
		}

	case msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_SRV_REG_TO_QUE:
		{
			outAttachData = pthis.process_PB_MSG_INTER_SRV_REG_TO_QUE(fd, protoMsg)
		}

	case msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_SRV_REPORT_TO_OTHER_QUE:
		{
			pthis.process_PB_MSG_INTER_SRV_REPORT_TO_OTHER_QUE(fd, protoMsg, inAttachData)
		}

	case msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_SRV_HEARTBEAT:
		{
			lin_common.LogDebug(fd, "PB_MSG_INTER_SRV_HEARTBEAT", " proto msg", protoMsg, " attach:", inAttachData)
			pthis.SendProtoMsg(fd, msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_SRV_HEARTBEAT_RES, &msgpacket.PB_MSG_INTER_SRV_HEARTBEAT{})
		}

	default:
		{
			lin_common.LogDebug(fd, "packType:", packType, " bytesProcess:", bytesProcess, " proto msg", protoMsg, " attach:", inAttachData)
		}
	}

	return
}




func (pthis*MsgQueSrv)process_PB_MSG_INTER_SRV_REPORT_TO_OTHER_QUE(fd lin_common.FD_DEF, pbMsg proto.Message, inAttachData interface{}) {
	lin_common.LogDebug(fd, "PB_MSG_INTER_SRV_REPORT_TO_OTHER_QUE:", pbMsg, " inAttachData:", inAttachData)

	attach, ok := inAttachData.(*tcpAttachDataMsgQueSrvAccept)
	if !ok || attach == nil {
		lin_common.LogErr("attach data err:",  inAttachData, pbMsg, fd)
		return
	}

	pbReport, ok := pbMsg.(*msgpacket.PB_MSG_INTER_SRV_REPORT_TO_OTHER_QUE)
	if !ok || pbReport == nil {
		return
	}

	// add other srv
	pthis.smgr.addOtherQueAllSrvFromPB(attach.queSrvID, pbReport.AllSrv)
}

func (pthis*MsgQueSrv)broadCastSrvReportToOtherQueSrv() {
	// ntf to all other que srv
	pbReport := &msgpacket.PB_MSG_INTER_SRV_REPORT_TO_OTHER_QUE{
		AllSrv : &msgpacket.PB_SRV_INFO_ALL{},
	}
	pthis.smgr.getAllSrvNetPB(pbReport.AllSrv)
	pthis.otherMgr.Range(func(key, value any) bool{
		qsi, ok := value.(otherMsgQueSrvInfo)
		if !ok {
			return true
		}
		pthis.SendProtoMsg(qsi.fdDial, msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_SRV_REPORT_TO_OTHER_QUE, pbReport)
		return true
	})
}


func (pthis*MsgQueSrv)process_TcpClose_SrvReg(fd lin_common.FD_DEF, inAttachData interface{}) {
	lin_common.LogInfo(fd, inAttachData)

	attachData, ok := inAttachData.(*tcpAttachDataSrvAccept)
	if !ok || attachData == nil {
		lin_common.LogInfo("tcp attach data err,", fd, inAttachData)
		return
	}

	pthis.smgr.delSrv(attachData.srvUUID)
	pthis.broadCastSrvReportToOtherQueSrv()
}

func (pthis*MsgQueSrv)process_PB_MSG_INTER_SRV_REG_TO_QUE(fd lin_common.FD_DEF, pbMsg proto.Message) interface{} {
	lin_common.LogDebug(fd, "PB_MSG_INTER_TYPE__PB_MSG_INTER_SRV_REG_TO_QUE:", pbMsg)

	pbReg, ok := pbMsg.(*msgpacket.PB_MSG_INTER_SRV_REG_TO_QUE)
	if !ok || pbReg == nil{
		return nil
	}

	// add srv to mgr,
	si := &SrvNetInfo{
		srvUUID :server_common.SRV_ID(pbReg.SrvUuid),
		srvType: server_common.SRV_TYPE(pbReg.SrvType),
		fd :fd,
		addr : lin_common.TcpGetPeerName(fd.FD),
	}

	pthis.smgr.addSrv(si)
	pthis.broadCastSrvReportToOtherQueSrv()

	// send response
	pbRes := &msgpacket.PB_MSG_INTER_SRV_REG_TO_QUE{
		SrvUuid :pbReg.SrvUuid,
		SrvType: pbReg.SrvType,
	}
	pthis.SendProtoMsg(fd, msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_SRV_REG_TO_QUE, pbRes)

	return &tcpAttachDataSrvAccept{srvUUID: si.srvUUID}
}

func (pthis*MsgQueSrv)process_PB_MSG_INTER_QUESRV_HEARTBEAT(fd lin_common.FD_DEF, pbMsg proto.Message, inAttachData interface{}) {
	pbHB := pbMsg.(*msgpacket.PB_MSG_INTER_QUESRV_HEARTBEAT)

	pbRes := &msgpacket.PB_MSG_INTER_QUESRV_HEARTBEAT_RES{
		QueSrvId : pbHB.QueSrvId,
	}
	lin_common.LogDebug("receive heartbeat from other msg que srv:", server_common.SRV_ID(pbHB.QueSrvId).String(), " inAttachData:", inAttachData)
	pthis.SendProtoMsg(fd, msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_HEARTBEAT_RES, pbRes)
}

func (pthis*MsgQueSrv)dialToOtherMsgQue(queSrvID server_common.SRV_ID, ip string, port int32){
	addr := ip + ":" + strconv.FormatInt(int64(port), 10)
	lin_common.LogInfo("addr", addr, " ", queSrvID.String())
	attachData := &tcpAttachDataMsgQueSrvDial{
		queSrvID: queSrvID,
	}
	fdDial, err := pthis.lsn.EPollListenerDial(addr, attachData, false)
	if err != nil {
		lin_common.LogErr("can't connect to:", addr, " ", queSrvID.String())
	}
	pthis.otherMgr.updateQueSrv(queSrvID,
		fdDial,
		ip,
		port)
}

func (pthis*MsgQueSrv)process_PB_MSG_INTER_QUESRV_REGISTER_RES(pbMsg proto.Message) {
	regRes, ok := pbMsg.(*msgpacket.PB_MSG_INTER_QUESRV_REGISTER_RES)
	if !ok || regRes == nil {
		return
	}

	pthis.otherMgr.Range(func(key, value any) bool{
		//close all current tcp connect to other msg que server
		qsi := value.(otherMsgQueSrvInfo)
		pthis.lsn.EPollListenerCloseTcp(qsi.fdDial, server_common.EN_TCP_CLOSE_REASON_reg_reconnect)
		pthis.lsn.EPollListenerCloseTcp(qsi.fdAccept, server_common.EN_TCP_CLOSE_REASON_reg_reconnect)
		return true
	})

	pthis.otherMgr.Clear()
	pthis.queSrvID = server_common.SRV_ID(regRes.QueSrvId)

	//connect to all other msg que srv
	for _, pbqsi := range regRes.QueSrvInfo {
		if pbqsi.QueSrvId == int64(pthis.queSrvID) {
			continue
		}

		pthis.dialToOtherMsgQue(server_common.SRV_ID(pbqsi.QueSrvId), pbqsi.Ip, pbqsi.Port)
	}
}

func (pthis*MsgQueSrv)process_PB_MSG_INTER_QUESRV_ONLINE_NTF(pbMsg proto.Message) {
	//connect to msg que srv
	pbNtf, ok := pbMsg.(*msgpacket.PB_MSG_INTER_QUESRV_ONLINE_NTF)
	if !ok || pbNtf == nil || nil == pbNtf.QueSrvInfo {
		return
	}

	pbqsi := pbNtf.QueSrvInfo
	pthis.dialToOtherMsgQue(server_common.SRV_ID(pbqsi.QueSrvId), pbqsi.Ip, pbqsi.Port)
}

func (pthis*MsgQueSrv)process_PB_MSG_INTER_QUESRV_OFFLINE_NTF(pbMsg proto.Message) {
	// disconnect to msg que srv and delete from map
	pbNtf, ok := pbMsg.(*msgpacket.PB_MSG_INTER_QUESRV_OFFLINE_NTF)
	if !ok || pbNtf == nil {
		return
	}

	queSrvID := server_common.SRV_ID(pbNtf.QueSrvId)

	qsi := &otherMsgQueSrvInfo{}
	ok = pthis.otherMgr.LoadAndDelete(queSrvID, qsi)
	if !ok {
		lin_common.LogInfo(queSrvID.String(), "where receive dial tcp close, not exist")
		return
	}
	pthis.lsn.EPollListenerCloseTcp(qsi.fdDial, server_common.EN_TCP_CLOSE_REASON_msgque_center_ntf_offline)
	pthis.lsn.EPollListenerCloseTcp(qsi.fdAccept, server_common.EN_TCP_CLOSE_REASON_recv_ntf_offline)
}

func (pthis*MsgQueSrv)deleteMsgQueSrvAndRedia(queSrvID server_common.SRV_ID) {
	qsi1 := otherMsgQueSrvInfo{}
	ok := pthis.otherMgr.Load(queSrvID, &qsi1)
	if !ok {
		return
	}

	pthis.smgr.delOtherQueAllSrv(queSrvID)

	pthis.lsn.EPollListenerCloseTcp(qsi1.fdDial, server_common.EN_TCP_CLOSE_REASON_recv_ntf_offline)
	pthis.lsn.EPollListenerCloseTcp(qsi1.fdAccept, server_common.EN_TCP_CLOSE_REASON_recv_ntf_offline)

	qsi1.fdDial = lin_common.FD_DEF_NIL
	qsi1.fdAccept = lin_common.FD_DEF_NIL
	pthis.otherMgr.Store(&qsi1)

	time.AfterFunc(time.Second * 3, func() {
		qsi := &otherMsgQueSrvInfo{}
		ok := pthis.otherMgr.Load(queSrvID, qsi)
		if !ok {
			lin_common.LogInfo(queSrvID.String(), "que srv not exist", queSrvID.String())
			return
		}
		if !qsi.fdDial.IsNull() {
			lin_common.LogInfo(queSrvID.String(), "fdDial is not null", "now:", qsi.fdDial)
			return
		}

		pthis.dialToOtherMsgQue(qsi.queSrvID, qsi.ip, qsi.port)
	})
}


func (pthis*MsgQueSrv)process_PB_MSG_INTER_QUESRV_CONNECT(fd lin_common.FD_DEF, pbMsg proto.Message, inAttachData interface{}) interface{} {
	pbConn, ok := pbMsg.(*msgpacket.PB_MSG_INTER_QUESRV_CONNECT)
	if !ok || pbConn == nil {
		return nil
	}

	queSrvID := server_common.SRV_ID(pbConn.QueSrvId)
	lin_common.LogInfo("que srv connect, fd:", fd, " ", queSrvID.String(), " pbConn:", pbConn, " inAttachData:", inAttachData)

	pthis.otherMgr.updateQueSrvAccept(queSrvID, fd)

	// send res pack
	pbRes := &msgpacket.PB_MSG_INTER_QUESRV_CONNECT_RES{
		QueSrvId: int64(queSrvID),
		AllSrv:&msgpacket.PB_SRV_INFO_ALL{},
	}
	pthis.smgr.getAllSrvNetPB(pbRes.AllSrv)
	pthis.SendProtoMsg(fd, msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_CONNECT_RES, pbRes)

	// add other srv
	pthis.smgr.addOtherQueAllSrvFromPB(queSrvID, pbConn.AllSrv)

	return &tcpAttachDataMsgQueSrvAccept{queSrvID: queSrvID}
}

func (pthis*MsgQueSrv)process_PB_MSG_INTER_QUESRV_CONNECT_RES(fd lin_common.FD_DEF, pbMsg proto.Message, inAttachData interface{}) {
	pbConnRes, ok := pbMsg.(*msgpacket.PB_MSG_INTER_QUESRV_CONNECT_RES)
	if !ok || pbConnRes == nil {
		return
	}

	attachData, ok := inAttachData.(*tcpAttachDataMsgQueSrvDial)
	if !ok || attachData == nil {
		lin_common.LogInfo("tcp attach data is nil, fd:", fd, " ", " pbConnRes:", pbConnRes, " inAttachData:", inAttachData)
	}

	lin_common.LogInfo("que srv connect, fd:", fd, " ", attachData.queSrvID.String(), " pbConnRes:", pbConnRes, " inAttachData:", inAttachData)

	// add other srv
	pthis.smgr.addOtherQueAllSrvFromPB(attachData.queSrvID, pbConnRes.AllSrv)
}


func (pthis*MsgQueSrv)SendProtoMsg(fd lin_common.FD_DEF, msgType msgpacket.PB_MSG_INTER_TYPE, protoMsg proto.Message){
	pthis.lsn.EPollListenerWrite(fd, msgpacket.ProtoPacketToBin(uint16(msgType), protoMsg))
}

// ConstructMsgQueSrv <addr> example 127.0.0.1:8888
func ConstructMsgQueSrv(msgqueCenterAddr string, addrBind string, addrOut string, epollCoroutineCount int) *MsgQueSrv {
	mqMgr := &MsgQueSrv{
		addrOut : addrOut,
		queSrvID : server_common.SRV_ID_INVALID,
		smgr : ConstructorSrvMgr(),
	}

	lsn, err := lin_common.ConstructorEPollListener(mqMgr, addrBind, epollCoroutineCount,
		lin_common.ParamEPollListener{
			ParamET: true,
			ParamEpollWaitTimeoutMills: 180 * 1000,
			ParamIdleClose: 600*1000,
			ParamNeedTick: true,
		})
	if err != nil {
		lin_common.LogErr("constructor epoll listener err:", err)
		return nil
	}
	mqMgr.lsn = lsn

	//connect to msg que center, when connect success send reg pack
	mqMgr.msgqueCenterAddr = msgqueCenterAddr
	mqMgr.fdCenter, err = lsn.EPollListenerDial(msgqueCenterAddr, &tcpAttachDataMsgQueCenterDial{}, false)
	if err != nil {
		lin_common.LogErr("dial to msg que center err:", err)
	}
	lin_common.LogInfo("connect end~~~~~~~~~~~~~~~~~~~~~")

	return mqMgr
}

func (pthis*MsgQueSrv)Wait() {
	pthis.lsn.EPollListenerWait()
}

func (pthis*MsgQueSrv)Dump(bDetail bool) (str string) {

	str += "\r\naddr out:" + pthis.addrOut + "\r\n"
	str += "msg que center:" + pthis.msgqueCenterAddr + pthis.fdCenter.String()
	str += pthis.queSrvID.String() + "\r\n\r\n"

	str += pthis.otherMgr.Dump()
	str += pthis.smgr.Dump()

	return
}
