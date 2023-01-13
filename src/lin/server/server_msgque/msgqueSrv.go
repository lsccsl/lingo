package main

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"lin/lin_common"
	"lin/msgpacket"
	"lin/server/server_common"
	"net"
	"strconv"
	"sync"
)

type MsgQueSrv struct {
	lsn *lin_common.EPollListener

	addrOut string
	queSrvID server_common.MSGQUE_SRV_ID

	mapOtherMsgQueSrv sync.Map // server_common.MSGQUE_SRV_ID - otherMsgQueSrvInfo
}

type otherMsgQueSrvInfo struct {
	fdDial lin_common.FD_DEF
	fdConn lin_common.FD_DEF
	ip string
	port int32
	queSrvID server_common.MSGQUE_SRV_ID
}

func (pthis *otherMsgQueSrvInfo)String() (str string) {
	str += "que srv id:" + pthis.queSrvID.ToString() +
		"[" + pthis.ip + ":" + strconv.FormatInt(int64(pthis.port), 10) + "]" +
		" fdDial:" + pthis.fdDial.String() +
		" fdConn:" + pthis.fdConn.String()

	return
}

type tcpAttachDataMsgQueSrvDial struct{
	queSrvID server_common.MSGQUE_SRV_ID
}
type tcpAttachDataMsgQueSrvConn struct {
	queSrvID server_common.MSGQUE_SRV_ID
}
type tcpAttachDataMsgQueCenter struct {
}

func (pthis*MsgQueSrv)TcpAcceptConnection(fd lin_common.FD_DEF, addr net.Addr, inAttachData interface{})(outAttachData interface{}) {
	lin_common.LogDebug(fd, addr, inAttachData)
	return nil
}

func (pthis*MsgQueSrv)TcpDialConnection(fd lin_common.FD_DEF, addr net.Addr, inAttachData interface{})(outAttachData interface{}) {
	lin_common.LogDebug(fd, addr, inAttachData)

	switch inAttachData.(type) {
	case *tcpAttachDataMsgQueCenter: // dial to msg que center tcp connection ok
		{
			// send reg msg to msq que center
			tcpAddr, err := net.ResolveTCPAddr("tcp", pthis.addrOut)
			if err != nil {
				lin_common.LogErr(err)
			}
			pbMsgReg := &msgpacket.PB_MSG_INTER_QUESRV_REGISTER{
				Ip: tcpAddr.IP.String(),
				Port: int32(tcpAddr.Port),
			}
			pthis.SendProtoMsg(fd, msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_REGISTER, pbMsgReg)

			return nil
		}

	case *tcpAttachDataMsgQueSrvDial: //  dial to other msg que srv tcp connection ok
		{
			// send conn msg to other msg que
			pbMsgConn := &msgpacket.PB_MSG_INTER_QUESRV_CONNECT{
				QueSrvId:int64(pthis.queSrvID),
			}
			pthis.SendProtoMsg(fd, msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_CONNECT, pbMsgConn)
		}
	}

	return
}

func (pthis*MsgQueSrv)TcpClose(fd lin_common.FD_DEF, closeReason lin_common.EN_TCP_CLOSE_REASON, inAttachData interface{}) {
	lin_common.LogDebug(fd, inAttachData, " closeReason:", closeReason)
}

func (pthis*MsgQueSrv)TcpOutBandData(fd lin_common.FD_DEF, data interface{}, inAttachData interface{}) {
	lin_common.LogDebug(fd, data, inAttachData)
}

func (pthis*MsgQueSrv)TcpData(fd lin_common.FD_DEF, readBuf *bytes.Buffer, inAttachData interface{})(bytesProcess int, outAttachData interface{}) {
	lin_common.LogDebug(fd)
	packType, bytesProcess, protoMsg := msgpacket.ProtoUnPacketFromBin(readBuf)
	if protoMsg == nil {
		return
	}
	lin_common.LogInfo("packType:", packType, " bytesProcess:", bytesProcess, " proto msg", protoMsg, "")

	switch msgpacket.PB_MSG_INTER_TYPE(packType) {
	case msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_REGISTER_RES:
		{
			pthis.process_PB_MSG_INTER_QUESRV_REGISTER_RES(protoMsg)
		}

	case msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_ONLINE_NTF:
		{
			pthis.process_PB_MSG_INTER_QUESRV_ONLINE_NTF(protoMsg)
		}

	case msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_OFFLINE_NTF:
		{
			pthis.process_PB_MSG_INTER_QUESRV_OFFLINE_NTF(protoMsg)
		}

	case msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_CONNECT:
		{
			outAttachData = pthis.process_PB_MSG_INTER_QUESRV_CONNECT(fd, protoMsg)
		}
	}

	return
}


func (pthis*MsgQueSrv)process_PB_MSG_INTER_QUESRV_REGISTER_RES(pbMsg proto.Message) {
	regRes, ok := pbMsg.(*msgpacket.PB_MSG_INTER_QUESRV_REGISTER_RES)
	if !ok || regRes == nil {
		return
	}

	pthis.mapOtherMsgQueSrv.Range(func(key, value any) bool{
		//close all current tcp connect to other msg que server
		qsi := value.(otherMsgQueSrvInfo)
		pthis.lsn.EPollListenerCloseTcp(qsi.fdDial, server_common.EN_TCP_CLOSE_REASON_reg_reconnect)
		pthis.lsn.EPollListenerCloseTcp(qsi.fdConn, server_common.EN_TCP_CLOSE_REASON_reg_reconnect)
		return true
	})

	pthis.mapOtherMsgQueSrv = sync.Map{}
	pthis.queSrvID = server_common.MSGQUE_SRV_ID(regRes.QueSrvId)

	//connect to all other msg que srv
	for _, pbqsi := range regRes.QueSrvInfo {
		if pbqsi.QueSrvId == int64(pthis.queSrvID) {
			continue
		}
		addr := pbqsi.Ip + ":" + strconv.FormatInt(int64(pbqsi.Port), 10)
		lin_common.LogInfo("addr", addr, pbqsi)
		attachData := &tcpAttachDataMsgQueSrvDial{
			queSrvID: server_common.MSGQUE_SRV_ID(pbqsi.QueSrvId),
		}

		fdDial, err := pthis.lsn.EPollListenerDial(addr, attachData, false)
		if err != nil {
			lin_common.LogErr("can't connect to:", addr, " ", server_common.MSGQUE_SRV_ID(pbqsi.QueSrvId).ToString())
		}
		qsi := otherMsgQueSrvInfo{
			fdDial:fdDial,
			ip:pbqsi.Ip,
			port:pbqsi.Port,
			queSrvID:server_common.MSGQUE_SRV_ID(pbqsi.QueSrvId),
		}
		pthis.mapOtherMsgQueSrv.Store(qsi.queSrvID, qsi)
	}
}

func (pthis*MsgQueSrv)process_PB_MSG_INTER_QUESRV_ONLINE_NTF(pbMsg proto.Message) {
	//connect to msg que srv
	pbNtf, ok := pbMsg.(*msgpacket.PB_MSG_INTER_QUESRV_ONLINE_NTF)
	if !ok || pbNtf == nil || nil == pbNtf.QueSrvInfo {
		return
	}

	pbqsi := pbNtf.QueSrvInfo
	queSrvID := server_common.MSGQUE_SRV_ID(pbqsi.QueSrvId)
	addr := pbqsi.Ip + ":" + strconv.FormatInt(int64(pbqsi.Port), 10)
	lin_common.LogInfo("addr", addr, " ", pbqsi, " ", queSrvID.ToString())
	attachData := &tcpAttachDataMsgQueSrvDial{
		queSrvID: server_common.MSGQUE_SRV_ID(pbqsi.QueSrvId),
	}

	fdDial, err := pthis.lsn.EPollListenerDial(addr, attachData, false)
	if err != nil {
		lin_common.LogErr("can't connect to:", addr, " ", queSrvID.ToString())
	}
	qsi := otherMsgQueSrvInfo{
		fdDial:fdDial,
		ip:pbqsi.Ip,
		port:pbqsi.Port,
		queSrvID:server_common.MSGQUE_SRV_ID(pbqsi.QueSrvId),
	}
	pthis.mapOtherMsgQueSrv.Store(queSrvID, qsi)
}

func (pthis*MsgQueSrv)process_PB_MSG_INTER_QUESRV_OFFLINE_NTF(pbMsg proto.Message) {
	// disconnect to msg que srv and delete from map
	pbNtf, ok := pbMsg.(*msgpacket.PB_MSG_INTER_QUESRV_OFFLINE_NTF)
	if !ok || pbNtf == nil {
		return
	}

	queSrvID := server_common.MSGQUE_SRV_ID(pbNtf.QueSrvId)
	val, ok := pthis.mapOtherMsgQueSrv.LoadAndDelete(queSrvID)
	if !ok {
		return
	}

	qsi, ok := val.(otherMsgQueSrvInfo)
	if !ok {
		return
	}

	pthis.lsn.EPollListenerCloseTcp(qsi.fdDial, server_common.EN_TCP_CLOSE_REASON_recv_ntf_offline)
	pthis.lsn.EPollListenerCloseTcp(qsi.fdConn, server_common.EN_TCP_CLOSE_REASON_recv_ntf_offline)
}

func (pthis*MsgQueSrv)process_PB_MSG_INTER_QUESRV_CONNECT(fd lin_common.FD_DEF, pbMsg proto.Message) interface{} {
	pbConn, ok := pbMsg.(*msgpacket.PB_MSG_INTER_QUESRV_CONNECT)
	if !ok || pbConn == nil {
		return nil
	}

	queSrvID := server_common.MSGQUE_SRV_ID(pbConn.QueSrvId)
	lin_common.LogInfo("que srv connect, fd:", fd, " ", queSrvID.ToString(), " pb msg:", pbConn)

	val, ok := pthis.mapOtherMsgQueSrv.Load(queSrvID)
	if !ok {
		lin_common.LogErr("can't find, que srv", queSrvID.ToString())
		return nil
	}
	qsi, ok := val.(otherMsgQueSrvInfo)
	if !ok {
		lin_common.LogErr("map value convert err, que srv", queSrvID.ToString())
		return nil
	}

	qsi.fdConn = fd
	pthis.mapOtherMsgQueSrv.Store(qsi.queSrvID, qsi)

	// send res pack
	pbRes := &msgpacket.PB_MSG_INTER_QUESRV_CONNECT_RES{QueSrvId: int64(queSrvID)}
	pthis.SendProtoMsg(fd, msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_CONNECT_RES, pbRes)

	return &tcpAttachDataMsgQueSrvConn{queSrvID: server_common.MSGQUE_SRV_ID(queSrvID)}
}

func (pthis*MsgQueSrv)SendProtoMsg(fd lin_common.FD_DEF, msgType msgpacket.PB_MSG_INTER_TYPE, protoMsg proto.Message){
	pthis.lsn.EPollListenerWrite(fd, msgpacket.ProtoPacketToBin(uint16(msgType), protoMsg))
}

// ConstructMsgQueSrv <addr> example 127.0.0.1:8888
func ConstructMsgQueSrv(msgqueCenterAddr string, addrBind string, addrOut string, epollCoroutineCount int) *MsgQueSrv {
	mqMgr := &MsgQueSrv{
		addrOut : addrOut,
	}

	lsn, err := lin_common.ConstructorEPollListener(mqMgr, addrBind, epollCoroutineCount, lin_common.ParamEPollListener{ParamET: true})
	if err != nil {
		lin_common.LogErr("constructor epoll listener err:", err)
		return nil
	}
	mqMgr.lsn = lsn

	//connect to msg que center, when connect success send reg pack
	lsn.EPollListenerDial(msgqueCenterAddr, &tcpAttachDataMsgQueCenter{}, false)

	return mqMgr
}

func (pthis*MsgQueSrv)Wait() {
	pthis.lsn.EPollListenerWait()
}

func (pthis*MsgQueSrv)Dump(bDetail bool) (str string) {

	str += "addr out:" + pthis.addrOut + "\r\n"
	str += "que srv id:" + pthis.queSrvID.ToString() + "\r\n\r\n"

	str += "connect to other msg que srv\r\n"
	pthis.mapOtherMsgQueSrv.Range(func(key, value any) bool{
		qsi, ok := value.(otherMsgQueSrvInfo)
		if !ok {
			return true
		}

		str += qsi.String() + "\r\n"
		return true
	})

	return
}
