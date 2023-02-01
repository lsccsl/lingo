package msgque_client

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"goserver/common"
	"goserver/msgpacket"
	"goserver/server/server_common"
	"net"
	"strconv"
	"sync"
)

type MgrQueClient struct {
	lsn *common.EPollListener
	srvUUID server_common.SRV_ID
	sryType server_common.SRV_TYPE

	fdCenter         common.FD_DEF
	msgqueCenterAddr string // "x.x.x.x:n"

	fdQueSrv   common.FD_DEF
	queSrvAddr string
	queDialWait sync.WaitGroup

	msgMgr *ClientSrvMsgMgr
	otherClisrvMgr *OtherClientSrvMgr

	cb MsgProcessCB
}

type MsgProcessCB interface {
	// new coroutine
	Go_CallBackMsg(pbMsg proto.Message, pbMsgType int32,
		srvUUIDFrom server_common.SRV_ID,
		srvType server_common.SRV_TYPE,
		timeOutMills int) (msgType int32, protoMsg proto.Message)

	// new coroutine
	Go_CallBackReport(pbLocal *msgpacket.PB_SRV_INFO_ALL)

	// new coroutine
	Go_CallBackNtf(ntf *msgpacket.PB_MSG_INTER_QUESRV_NTF)
}


type tcpAttachDataMsgQueCenter struct {
}
type tcpAttachDataMsgQueDial struct {

}

func (pthis*MgrQueClient)TcpAcceptConnection(fd common.FD_DEF, addr net.Addr, inAttachData interface{})(outAttachData interface{}) {
	return nil
}

func (pthis*MgrQueClient)TcpDialConnection(fd common.FD_DEF, addr net.Addr, inAttachData interface{})(outAttachData interface{}) {

	switch t := inAttachData.(type) {
	case *tcpAttachDataMsgQueCenter:
		{
			pthis.process_TcpDialToMsgCenter(fd)
		}

	default:
		{
			common.LogDebug(fd, addr, t)
		}
	}

	return nil
}

func (pthis*MgrQueClient)TcpData(fd common.FD_DEF, readBuf *bytes.Buffer, inAttachData interface{})(bytesProcess int, outAttachData interface{}) {
	packType, bytesProcess, protoMsg := msgpacket.ProtoUnPacketFromBin(readBuf)

	switch msgpacket.PB_MSG_TYPE(packType) {
	case msgpacket.PB_MSG_TYPE__PB_MSG_INTER_CLISRV_REG_MSGQUE_CENTER_RES:
		{
			pthis.process_PB_MSG_INTER_CLISRV_REG_MSGQUE_CENTER_RES(fd, protoMsg)
		}

	case msgpacket.PB_MSG_TYPE__PB_MSG_INTER_CLISRV_REG_TO_QUE_RES:
		{
			pthis.process_PB_MSG_INTER_CLISRV_REG_TO_QUE_RES(protoMsg)
		}

	case msgpacket.PB_MSG_TYPE__PB_MSG_INTER_QUESRV_REPORT_BROADCAST:
		{
			pthis.process_PB_MSG_INTER_QUESRV_REPORT_BROADCAST(protoMsg)
		}
	case msgpacket.PB_MSG_TYPE__PB_MSG_INTER_QUESRV_NTF:
		{
			pthis.process_PB_MSG_INTER_QUESRV_NTF(protoMsg)
		}

	case msgpacket.PB_MSG_TYPE__PB_MSG_INTER_MSG:
		{
			pthis.process_PB_MSG_INTER_MSG(fd, protoMsg, inAttachData)
		}

	case msgpacket.PB_MSG_TYPE__PB_MSG_INTER_MSG_RES:
		{
			pthis.process_PB_MSG_INTER_MSG_RES(fd, protoMsg, inAttachData)
		}

	case msgpacket.PB_MSG_TYPE__PB_MSG_INTER_CLISRV_HEARTBEAT_RES:
		{
			common.LogDebug(fd, "PB_MSG_INTER_CLISRV_HEARTBEAT_RES", " attach:", inAttachData)
		}

	default:
		{
			common.LogDebug(fd, "packType:", packType, " bytesProcess:", bytesProcess, " proto msg", protoMsg, " attach:", inAttachData)
		}
	}

	return
}

func (pthis*MgrQueClient)TcpClose(fd common.FD_DEF, closeReason common.EN_TCP_CLOSE_REASON, inAttachData interface{}) {
	// redial to que center
	switch inAttachData.(type) {
	case *tcpAttachDataMsgQueDial:
		{
			var err error
			pthis.fdCenter, err = pthis.lsn.EPollListenerDial(pthis.msgqueCenterAddr, &tcpAttachDataMsgQueCenter{}, false)
			if err != nil {
				common.LogErr("dial to msg que center err:", err)
			}
		}
	}
}

func (pthis*MgrQueClient)TcpOutBandData(fd common.FD_DEF, data interface{}, inAttachData interface{}) {
}

func (pthis*MgrQueClient)TcpTick(fd common.FD_DEF, tNowMill int64, inAttachData interface{}) {
	//lin_common.LogDebug(fd, " tNowMill:", tNowMill, " inAttachData:", inAttachData)
	pbHB := &msgpacket.PB_MSG_INTER_CLISRV_HEARTBEAT{SrvUuid: int64(pthis.srvUUID)}
	pthis.SendProtoMsg(fd, msgpacket.PB_MSG_TYPE__PB_MSG_INTER_CLISRV_HEARTBEAT, pbHB)
}



func (pthis*MgrQueClient)process_PB_MSG_INTER_CLISRV_REG_TO_QUE_RES(pbMsg proto.Message) {
	pthis.queDialWait.Done()

	pbRes, ok := pbMsg.(*msgpacket.PB_MSG_INTER_CLISRV_REG_TO_QUE_RES)
	if !ok  || nil == pbRes {
		return
	}
	common.LogDebug(pbRes)

	pthis.otherClisrvMgr.SetOtherClientSrvFromPB(pbRes.AllSrv)
	if pthis.cb != nil {
		go pthis.cb.Go_CallBackReport(pbRes.LocalAllSrv)
	}
}

func (pthis*MgrQueClient)process_PB_MSG_INTER_QUESRV_REPORT_BROADCAST(pbMsg proto.Message) {
	pbReport, ok := pbMsg.(*msgpacket.PB_MSG_INTER_QUESRV_REPORT_BROADCAST)
	if !ok  || nil == pbReport {
		return
	}
	common.LogDebug(pbReport)

	pthis.otherClisrvMgr.SetOtherClientSrvFromPB(pbReport.AllSrv)
	if pthis.cb != nil {
		go pthis.cb.Go_CallBackReport(pbReport.LocalAllSrv)
	}
}

func (pthis*MgrQueClient)process_PB_MSG_INTER_QUESRV_NTF(pbMsg proto.Message) {
	pbNtf, ok := pbMsg.(*msgpacket.PB_MSG_INTER_QUESRV_NTF)
	if !ok  || nil == pbNtf {
		return
	}
	common.LogDebug(pbNtf)

	if pthis.cb != nil {
		go pthis.cb.Go_CallBackNtf(pbNtf)
	}
}

func (pthis*MgrQueClient)process_PB_MSG_INTER_CLISRV_REG_MSGQUE_CENTER_RES(fd common.FD_DEF, pbMsg proto.Message) {
	common.LogDebug(fd, pbMsg)
	// dial to msg que srv
	pbRes := pbMsg.(*msgpacket.PB_MSG_INTER_CLISRV_REG_MSGQUE_CENTER_RES)

	if pbRes.Res != msgpacket.PB_RESPONSE_CODE_PB_RESPONSE_CODE_OK {
		common.LogErr("reg to msg que center err:", pbRes.Res)
	}

	pthis.srvUUID = server_common.SRV_ID(pbRes.SrvUuid)
	pthis.queSrvAddr = pbRes.QueSrvIp + ":" + strconv.FormatInt(int64(pbRes.QueSrvPort), 10)

	var err error
	pthis.fdQueSrv, err = pthis.lsn.EPollListenerDial(pthis.queSrvAddr, &tcpAttachDataMsgQueDial{}, true)
	if err != nil {
		common.LogErr("dial to msg que center err:", err)
	}

	// send reg to msq que srv
	{
		pbReg := &msgpacket.PB_MSG_INTER_CLISRV_REG_TO_QUE{
			SrvUuid: int64(pthis.srvUUID),
			SrvType: int32(pthis.sryType),
		}
		pthis.SendProtoMsg(pthis.fdQueSrv, msgpacket.PB_MSG_TYPE__PB_MSG_INTER_CLISRV_REG_TO_QUE, pbReg)
	}

	pthis.lsn.EPollListenerCloseTcp(fd, server_common.EN_TCP_CLOSE_REASON_srv_reg_ok)
}

func (pthis*MgrQueClient)process_TcpDialToMsgCenter(fd common.FD_DEF) {
	common.LogDebug("dial to msg center success")
	pbReg := &msgpacket.PB_MSG_INTER_CLISRV_REG_MSGQUE_CENTER{
		SrvUuid: int64(pthis.srvUUID),
		SrvType: int32(pthis.sryType),
	}

	pthis.SendProtoMsg(fd, msgpacket.PB_MSG_TYPE__PB_MSG_INTER_CLISRV_REG_MSGQUE_CENTER, pbReg)
}

func (pthis*MgrQueClient)SendProtoMsg(fd common.FD_DEF, msgType msgpacket.PB_MSG_TYPE, protoMsg proto.Message){
	pthis.lsn.EPollListenerWrite(fd, msgpacket.ProtoPacketToBin(uint16(msgType), protoMsg))
}


func (pthis*MgrQueClient)MgrQueClientSetCallBack(cb MsgProcessCB) {
	pthis.cb = cb
}

func (pthis*MgrQueClient)GetCliSrvByType(srvType server_common.SRV_TYPE) []*server_common.SrvBaseInfo {
	return pthis.otherClisrvMgr.GetCliSrvByType(srvType)
}

func (pthis*MgrQueClient)WaitEpoll() {
	pthis.lsn.EPollListenerWait()
}
func (pthis*MgrQueClient)DialToQueSrv() {
	pthis.queDialWait.Add(1)
	// dial to msg que center server
	var err error
	pthis.fdCenter, err = pthis.lsn.EPollListenerDial(pthis.msgqueCenterAddr, &tcpAttachDataMsgQueCenter{}, false)
	if err != nil {
		common.LogErr("dial to msg que center err:", err)
	}

	pthis.queDialWait.Wait()
	common.LogInfo("connect end~~~~~~~~~~~~~~~~~~~~~")
}

func ConstructMgrQueClient(msgqueCenterAddr string, srvType server_common.SRV_TYPE, cb MsgProcessCB) *MgrQueClient {
	mqCli := &MgrQueClient{
		srvUUID :       server_common.SRV_ID_INVALID,
		sryType:        srvType,
		msgMgr:         ConstructClientSrvMsgMgr(),
		otherClisrvMgr: ConstructOtherClientSrvMgr(),
		cb:             cb,
	}

	lsn, err := common.ConstructorEPollListener(mqCli, "", 1,
		common.ParamEPollListener{
			ParamET: true,
			ParamEpollWaitTimeoutMills: 180 * 1000,
			ParamIdleClose: 600*1000,
			ParamNeedTick: true,
		})
	if err != nil {
		common.LogErr("constructor epoll listener err:", err)
		return nil
	}
	mqCli.lsn = lsn
	mqCli.msgqueCenterAddr = msgqueCenterAddr

	return mqCli
}

func (pthis*MgrQueClient)Dump(bDetail bool) string {
	str := pthis.lsn.EPollListenerDump()

	str += "fdQueSrv:" + pthis.fdQueSrv.String() + " addr:" + pthis.queSrvAddr

	str += pthis.srvUUID.String() + pthis.sryType.String()

	str += pthis.otherClisrvMgr.Dump()
	str += pthis.msgMgr.Dump()

	return str
}
