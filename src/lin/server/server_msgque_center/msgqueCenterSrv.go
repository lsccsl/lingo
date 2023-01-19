package main

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"lin/lin_common"
	"lin/msgpacket"
	"lin/server/server_common"
	"net"
)

// MsgQueCenterSrv this struct is access by multi go coroutine, all member must be 'thread safe'
type MsgQueCenterSrv struct {
	lsn *lin_common.EPollListener

	queSrvMgr *MsgQueSrvMgr

	//srvIDSeed atomic.Int64
}




type tcpAttachDataMsgQueSrv struct {
	queSrvID server_common.SRV_ID
}

type tcpAttachDataSrv struct {
	queSrvID server_common.SRV_ID
}

func (pthis*MsgQueCenterSrv)TcpAcceptConnection(fd lin_common.FD_DEF, addr net.Addr, inAttachData interface{})(outAttachData interface{}) {
	lin_common.LogDebug(addr)
	return nil
}

func (pthis*MsgQueCenterSrv)TcpDialConnection(fd lin_common.FD_DEF, addr net.Addr, inAttachData interface{})(outAttachData interface{}) {
	lin_common.LogDebug(addr)
	return nil
}

func (pthis*MsgQueCenterSrv)TcpData(fd lin_common.FD_DEF, readBuf *bytes.Buffer, inAttachData interface{})(bytesProcess int, outAttachData interface{}) {
	packType, bytesProcess, protoMsg := msgpacket.ProtoUnPacketFromBin(readBuf)
	if protoMsg == nil {
		return
	}

	outAttachData = nil
	switch msgpacket.PB_MSG_INTER_TYPE(packType) {
	case msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUECENTER_REGISTER:
		{
			outAttachData = pthis.process_PB_MSG_INTER_QUECENTER_REGISTER(fd, protoMsg)
		}

	case msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUECENTER_HEARTBEAT:
		{
			lin_common.LogDebug(fd, "PB_MSG_INTER_QUECENTER_HEARTBEAT:", protoMsg, " attach:", inAttachData)
			pthis.process_PB_MSG_INTER_QUESRV_HEARTBEAT(fd, protoMsg)
		}

	case msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_CLISRV_REG_MSGQUE_CENTER:
		{
			pthis.process_PB_MSG_INTER_CLISRV_REG_MSGQUE_CENTER(fd, protoMsg)
		}

	case msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_REPORT_TO_OTHER_QUE:
		{
			lin_common.LogDebug(fd, "PB_MSG_INTER_QUESRV_REPORT_TO_OTHER_QUE:", protoMsg, " attach:", inAttachData)
			pbReport, ok := protoMsg.(*msgpacket.PB_MSG_INTER_QUESRV_REPORT_TO_OTHER_QUE)
			count := 0
			if ok {
				if pbReport.AllSrv != nil {
					if pbReport.AllSrv.ArraySrv != nil {
						count = len(pbReport.AllSrv.ArraySrv)
					}
				}
			}
			pthis.queSrvMgr.ResetQueSrvChooseCount(server_common.SRV_ID(pbReport.QueSrvId), count)
		}

	case msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_MSG:
		{
			lin_common.LogDebug(fd, "PB_MSG_INTER_SRV_MSG:", protoMsg, " attach:", inAttachData)
			pthis.process_PB_MSG_INTER_MSG(fd, protoMsg, inAttachData)
		}

	default:
		{
			lin_common.LogInfo("packType:", packType, " bytesProcess:", bytesProcess)
		}
	}

	return
}

func (pthis*MsgQueCenterSrv)TcpClose(fd lin_common.FD_DEF, closeReason lin_common.EN_TCP_CLOSE_REASON, inAttachData interface{}) {
	lin_common.LogInfo(fd, " closeReason:", closeReason, " inAttachData:", inAttachData)

	switch t := inAttachData.(type) {
	case *tcpAttachDataMsgQueSrv:
		{
			pthis.process_TcpClose_MsgQueSrv(fd, t)
		}
	}
}

func (pthis*MsgQueCenterSrv)TcpOutBandData(fd lin_common.FD_DEF, data interface{}, inAttachData interface{}) {
	lin_common.LogDebug(fd)
}

func (pthis*MsgQueCenterSrv)TcpTick(fd lin_common.FD_DEF, tNowMill int64, inAttachData interface{}){
	lin_common.LogDebug(fd, " tNowMill:", tNowMill, " inAttachData:", inAttachData)
}



func (pthis*MsgQueCenterSrv)process_PB_MSG_INTER_CLISRV_REG_MSGQUE_CENTER(fd lin_common.FD_DEF, pbMsg proto.Message) {
	lin_common.LogDebug(fd, "PB_MSG_INTER_CLISRV_REG_MSGQUE_CENTER:", pbMsg)
	// choose a que srv
	pbReg := pbMsg.(*msgpacket.PB_MSG_INTER_CLISRV_REG_MSGQUE_CENTER)

	pbRes := &msgpacket.PB_MSG_INTER_CLISRV_REG_MSGQUE_CENTER_RES{
		Res:msgpacket.PB_RESPONSE_CODE_PB_RESPONSE_CODE_OK,
	}

	defer pthis.SendProtoMsg(fd, msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_CLISRV_REG_MSGQUE_CENTER_RES, pbRes)

	// choose a que srv
	qsi, ok := pthis.queSrvMgr.ChooseMostIdleQueSrv()
	if !ok {
		pbRes.Res = msgpacket.PB_RESPONSE_CODE_PB_RESPONSE_CODE_no_que_srv
		return
	}

	srvUUID := server_common.SRV_ID(pbReg.SrvUuid)
	if srvUUID == server_common.SRV_ID_INVALID {
		srvUUID = pthis.genSrvID()
	}
	pbRes.SrvUuid = int64(srvUUID)
	pbRes.QueSrvId = int64(qsi.queSrvID)
	pbRes.QueSrvIp = qsi.ip
	pbRes.QueSrvPort = qsi.port

	// and client redirect to msg que srv that been chosen
}

func (pthis*MsgQueCenterSrv)process_PB_MSG_INTER_QUESRV_HEARTBEAT(fd lin_common.FD_DEF, pbMsg proto.Message) {
	pbHB := pbMsg.(*msgpacket.PB_MSG_INTER_QUECENTER_HEARTBEAT)
	if pbHB == nil {
		return
	}
	// send heartbeat back
	pbHBRes := &msgpacket.PB_MSG_INTER_QUECENTER_HEARTBEAT_RES{}
	pbHBRes.QueSrvId = pbHB.QueSrvId
	lin_common.LogDebug("receive heartbeat ", server_common.SRV_ID(pbHB.QueSrvId).String())

	pthis.SendProtoMsg(fd, msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUECENTER_HEARTBEAT_RES, pbHBRes)
}

func (pthis*MsgQueCenterSrv)process_PB_MSG_INTER_QUECENTER_REGISTER(fd lin_common.FD_DEF, pbMsg proto.Message) interface{}{
	lin_common.LogDebug(fd, "PB_MSG_INTER_QUESRV_REGISTER:", pbMsg, " attach:")
	regMsg, ok := pbMsg.(*msgpacket.PB_MSG_INTER_QUECENTER_REGISTER)
	if !ok || regMsg == nil {
		lin_common.LogErr("err msg", pbMsg)
		return nil
	}

	queSrvID := server_common.SRV_ID(regMsg.QueSrvId)
	if queSrvID == server_common.SRV_ID_INVALID {
		queSrvID = pthis.genSrvID()
	}

	//assign id
	qsiReg := MsgQueSrvInfo{
		queSrvID : queSrvID,
		fd :       fd,
		ip :       regMsg.Ip,
		port:      regMsg.Port,
	}
	//add msg que server list
	pthis.queSrvMgr.StoreQueSrvInfo(&qsiReg)
	if regMsg.AllSrv != nil {
		if regMsg.AllSrv.ArraySrv != nil {
			pthis.queSrvMgr.ResetQueSrvChooseCount(queSrvID, len(regMsg.AllSrv.ArraySrv))
		}
	}

	//response
	regRet := &msgpacket.PB_MSG_INTER_QUECENTER_REGISTER_RES{}
	regRet.QueSrvId = int64(qsiReg.queSrvID)

	pthis.queSrvMgr.RangeQueSrvInfo(func(key, value any) bool{

		qsi, ok := value.(MsgQueSrvInfo)
		if !ok {
			return true
		}
		queSrvInfo := &msgpacket.PB_MSG_INTER_QUESRV_INFO {
			QueSrvId: int64(qsi.queSrvID),
			Ip:qsi.ip,
			Port: qsi.port,
		}
		regRet.QueSrvInfo = append(regRet.QueSrvInfo, queSrvInfo)

		return true
	})

	pthis.SendProtoMsg(fd, msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUECENTER_REGISTER_RES, regRet)

	pthis.queSrvMgr.RangeQueSrvInfo(func(key, value any) bool{
		qsi, ok := value.(MsgQueSrvInfo)
		if !ok {
			return true
		}
		if qsi.queSrvID == qsiReg.queSrvID {
			return true
		}
		//notify other msg que srv online
		ntf := &msgpacket.PB_MSG_INTER_QUECENTER_ONLINE_NTF{
			QueSrvInfo : &msgpacket.PB_MSG_INTER_QUESRV_INFO{
				QueSrvId:int64(qsiReg.queSrvID),
				Ip: qsiReg.ip,
				Port: qsiReg.port,
			},
		}
		pthis.SendProtoMsg(qsi.fd, msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUECENTER_ONLINE_NTF, ntf)

		return true
	})

	return &tcpAttachDataMsgQueSrv{
		qsiReg.queSrvID,
	}
}

func (pthis*MsgQueCenterSrv)process_TcpClose_MsgQueSrv(fd lin_common.FD_DEF, attachData *tcpAttachDataMsgQueSrv) {
	lin_common.LogInfo(attachData.queSrvID.String())

	qsi, ok := pthis.queSrvMgr.LoadQueSrvInfo(attachData.queSrvID)
	if !ok {
		lin_common.LogErr(attachData.queSrvID.String(), " can't find")
		return
	}

	if !qsi.fd.IsSame(&fd) {
		lin_common.LogErr(attachData.queSrvID.String(), " fd is not same, current:", qsi.fd, " close:", fd)
		return
	}

	pthis.queSrvMgr.DeleteQueSrvInfo(attachData.queSrvID)

	pthis.queSrvMgr.RangeQueSrvInfo(func(key, value any) bool{
		qsi, ok := value.(MsgQueSrvInfo)
		if !ok {
			return true
		}
		//notify other msg que srv offline
		ntf := &msgpacket.PB_MSG_INTER_QUECENTER_OFFLINE_NTF{
			QueSrvId:int64(attachData.queSrvID),
		}
		pthis.SendProtoMsg(qsi.fd, msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUECENTER_OFFLINE_NTF, ntf)
		return true
	})
}

func (pthis*MsgQueCenterSrv)genSrvID() server_common.SRV_ID {
	return server_common.SRV_ID(lin_common.GenUUID64_V4())
	//return server_common.MSGQUE_SRV_ID(pthis.srvIDSeed.Add(1))
}

func (pthis*MsgQueCenterSrv)SendProtoMsg(fd lin_common.FD_DEF, msgType msgpacket.PB_MSG_INTER_TYPE, protoMsg proto.Message){
	pthis.lsn.EPollListenerWrite(fd, msgpacket.ProtoPacketToBin(uint16(msgType), protoMsg))
}

func (pthis*MsgQueCenterSrv)Wait() {
	pthis.lsn.EPollListenerWait()
}

// ConstructMsgQueCenterSrv <addr> example 127.0.0.1:8888
func ConstructMsgQueCenterSrv(addr string, epollCoroutineCount int) *MsgQueCenterSrv {
	mqMgr := &MsgQueCenterSrv{
		queSrvMgr: ConstructMsgQueSrvMgr(),
	}
	//mqMgr.srvIDSeed.Store(1)

	lsn, err := lin_common.ConstructorEPollListener(mqMgr, addr, epollCoroutineCount,
		lin_common.ParamEPollListener{ParamET: true,
			ParamEpollWaitTimeoutMills:180*1000,
			ParamIdleClose:600 * 1000,
		})
	if err != nil {
		lin_common.LogErr("constructor epoll listener err:", err)
		return nil
	}
	mqMgr.lsn = lsn

	return mqMgr
}

func (pthis*MsgQueCenterSrv)Dump(bDetail bool) (str string) {

	//str = "\r\nque srv id seed:" + strconv.FormatInt(pthis.srvIDSeed.Load(), 10) + "\r\n"

	str += pthis.lsn.EPollListenerDump()

	str += pthis.queSrvMgr.Dump()

	return
}
