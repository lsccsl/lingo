package main

import (
	"bytes"
	"goserver/common"
	"goserver/msgpacket"
	"goserver/server/server_common"
	msgque_client "goserver/server/server_msg_que_client"
	"net"
)

type LogonSrv struct {
	mqClient *msgque_client.MgrQueClient

	lsn *common.EPollListener

	outAddr string

	centerSrvUUID server_common.SRV_ID
	dbSrvUUID server_common.SRV_ID
}

func (pthis*LogonSrv)Wait() {
	pthis.mqClient.WaitEpoll()
}

func (pthis*LogonSrv)TcpAcceptConnection(fd common.FD_DEF, addr net.Addr, inAttachData interface{})(outAttachData interface{}) {
	common.LogDebug(fd, addr, inAttachData)
	return nil
}
func (pthis*LogonSrv)TcpDialConnection(fd common.FD_DEF, addr net.Addr, inAttachData interface{})(outAttachData interface{}) {
	common.LogDebug(fd, addr, inAttachData)
	return
}
func (pthis*LogonSrv)TcpData(fd common.FD_DEF, readBuf *bytes.Buffer, inAttachData interface{})(bytesProcess int, outAttachData interface{}) {

	packType, bytesProcess, protoMsg := msgpacket.ProtoUnPacketFromBin(readBuf)
	if protoMsg == nil {
		return
	}
	common.LogDebug("packType:", packType, protoMsg)

	switch t := protoMsg.(type) {
	case *msgpacket.PB_MSG_LOGON:
		pthis.process_PB_MSG_LOGON(fd, t)
	}

	return
}
func (pthis*LogonSrv)TcpClose(fd common.FD_DEF, closeReason common.EN_TCP_CLOSE_REASON, inAttachData interface{}) {
	common.LogDebug(fd, closeReason, inAttachData)
}
func (pthis*LogonSrv)TcpOutBandData(fd common.FD_DEF, data interface{}, inAttachData interface{}) {
	common.LogDebug(fd, data, inAttachData)
}
func (pthis*LogonSrv)TcpTick(fd common.FD_DEF, tNowMill int64, inAttachData interface{}) {
	common.LogDebug(fd, tNowMill, inAttachData)
}

func (pthis*LogonSrv)GetCenterSrvUUID() (centerSrvUUID server_common.SRV_ID) {
	pbMsg := &msgpacket.PB_MSG_INTER_QUESRV_GET_SRVTYPE{}
	pbMsg.SrvType = int32(server_common.SRV_TYPE_center_server)
	pbMsgRes, err := pthis.mqClient.SendMsgToSrvType(server_common.SRV_TYPE_msq_que,
		msgpacket.PB_MSG_TYPE__PB_MSG_INTER_QUESRV_GET_SRVTYPE, pbMsg)

	if nil != err {
		common.LogErr("get center server uuid err:", err)
		return server_common.SRV_ID_INVALID
	}

	pbRes, ok := pbMsgRes.(*msgpacket.PB_MSG_INTER_QUESRV_GET_SRVTYPE_RES)
	if !ok {
		common.LogErr("get center server uuid err")
		return server_common.SRV_ID_INVALID
	}
	if nil == pbRes.ArrarySrv {
		common.LogErr("get center server uuid err")
		return server_common.SRV_ID_INVALID
	}
	if 0 == len(pbRes.ArrarySrv) {
		common.LogErr("get center server uuid err")
		return server_common.SRV_ID_INVALID
	}

	centerSrvUUID = server_common.SRV_ID(pbRes.ArrarySrv[0].SrvUuid)
	common.LogInfo("centerSrvUUID:", centerSrvUUID)
	pthis.centerSrvUUID = centerSrvUUID
	return
}

func (pthis*LogonSrv)GetDBSrvUUID() {
	pbMsg := &msgpacket.PB_MSG_INTER_QUESRV_GET_SRVTYPE{}
	pbMsg.SrvType = int32(server_common.SRV_TYPE_database_server)
	pbMsgRes, err := pthis.mqClient.SendMsgToSrvType(server_common.SRV_TYPE_msq_que,
		msgpacket.PB_MSG_TYPE__PB_MSG_INTER_QUESRV_GET_SRVTYPE, pbMsg)

	if nil != err {
		common.LogErr("get center server uuid err:", err)
		return
	}

	pbRes, ok := pbMsgRes.(*msgpacket.PB_MSG_INTER_QUESRV_GET_SRVTYPE_RES)
	if !ok {
		common.LogErr("get center server uuid err")
		return
	}
	if nil == pbRes.ArrarySrv {
		common.LogErr("get center server uuid err")
		return
	}
	if 0 == len(pbRes.ArrarySrv) {
		common.LogErr("get center server uuid err")
		return
	}

	pthis.dbSrvUUID = server_common.SRV_ID(pbRes.ArrarySrv[0].SrvUuid)
	common.LogInfo("dbSrvUUID:", pthis.dbSrvUUID)
	return
}


func ConstructLogonSrv(id string)*LogonSrv {
	logonCfg := server_common.GetLogonsrvCfg(id)
	common.LogInfo(logonCfg)

	ls := &LogonSrv{}
	ls.outAddr = logonCfg.OutAddr

	// create client tcp epoll listener
	var err error
	ls.lsn, err = common.ConstructorEPollListener(ls, logonCfg.BindAddr, 50,
		common.ParamEPollListener{
			ParamET: true,
			ParamEpollWaitTimeoutMills: 180 * 1000,
			ParamIdleClose: 600*1000,
			ParamNeedTick: true,
		})
	if err != nil {
		common.LogErr("create epoll err")
		return nil
	}

	ls.mqClient = msgque_client.ConstructMgrQueClient(server_common.Global_ServerCfg.MsgQueCent.OutAddr, server_common.SRV_TYPE_logon_server, ls)
	ls.mqClient.DialToQueSrv()

	// get center server uuid
	ls.GetCenterSrvUUID()
	ls.GetDBSrvUUID()
	common.LogInfo("center server:", ls.centerSrvUUID, " db server:", ls.dbSrvUUID)

	return ls
}

