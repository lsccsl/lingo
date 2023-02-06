package main

import (
	"bytes"
	"goserver/common"
	"goserver/msgpacket"
	"goserver/server/server_common"
	msgque_client "goserver/server/server_msg_que_client"
	"net"
)

type GameSrv struct {
	mqClient *msgque_client.MgrQueClient

	lsn *common.EPollListener

	outAddr string

	centerSrvUUID server_common.SRV_ID
}

func (pthis*GameSrv)Wait() {
	pthis.mqClient.WaitEpoll()
}

func (pthis*GameSrv)TcpAcceptConnection(fd common.FD_DEF, addr net.Addr, inAttachData interface{})(outAttachData interface{}) {
	common.LogDebug(fd, addr, inAttachData)
	return nil
}
func (pthis*GameSrv)TcpDialConnection(fd common.FD_DEF, addr net.Addr, inAttachData interface{})(outAttachData interface{}) {
	common.LogDebug(fd, addr, inAttachData)
	return nil
}
func (pthis*GameSrv)TcpData(fd common.FD_DEF, readBuf *bytes.Buffer, inAttachData interface{})(bytesProcess int, outAttachData interface{}) {
	packType, bytesProcess, protoMsg := msgpacket.ProtoUnPacketFromBin(readBuf)
	if protoMsg == nil {
		return
	}
	common.LogDebug("packType:", packType, protoMsg)
	return
}
func (pthis*GameSrv)TcpClose(fd common.FD_DEF, closeReason common.EN_TCP_CLOSE_REASON, inAttachData interface{}) {

}
func (pthis*GameSrv)TcpOutBandData(fd common.FD_DEF, data interface{}, inAttachData interface{}) {

}
func (pthis*GameSrv)TcpTick(fd common.FD_DEF, tNowMill int64, inAttachData interface{}) {

}

func (pthis*GameSrv)GetCenterSrvUUID() (centerSrvUUID server_common.SRV_ID) {
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


func ConstructGameSrv(id string)*GameSrv {
	gsCfg := server_common.GetGameSrvCfg(id)
	common.LogInfo(gsCfg)

	gs := &GameSrv{}
	gs.outAddr = gsCfg.OutAddr

	// create client tcp epoll listener
	var err error
	gs.lsn, err = common.ConstructorEPollListener(gs, gsCfg.BindAddr, 10,
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

	gs.mqClient = msgque_client.ConstructMgrQueClient(server_common.Global_ServerCfg.MsgQueCent.OutAddr, server_common.SRV_TYPE_game_server, gs)
	gs.mqClient.DialToQueSrv()

	// get center server uuid
	gs.GetCenterSrvUUID()
	common.LogInfo("center server:", gs.centerSrvUUID)

	return gs
}
