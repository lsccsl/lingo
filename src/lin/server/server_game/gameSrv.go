package main

import (
	"bytes"
	"lin/lin_common"
	"lin/msgpacket"
	"lin/server/server_common"
	"lin/server/server_msg_que_client"
	"net"
)

type GameSrv struct {
	mqClient *msgque_client.MgrQueClient

	lsn *lin_common.EPollListener

	outAddr string

	centerSrvUUID server_common.SRV_ID
}

func (pthis*GameSrv)Wait() {
	pthis.mqClient.WaitEpoll()
}

func (pthis*GameSrv)TcpAcceptConnection(fd lin_common.FD_DEF, addr net.Addr, inAttachData interface{})(outAttachData interface{}) {
	return nil
}
func (pthis*GameSrv)TcpDialConnection(fd lin_common.FD_DEF, addr net.Addr, inAttachData interface{})(outAttachData interface{}) {
	return nil
}
func (pthis*GameSrv)TcpData(fd lin_common.FD_DEF, readBuf *bytes.Buffer, inAttachData interface{})(bytesProcess int, outAttachData interface{}) {
	return 0, nil
}
func (pthis*GameSrv)TcpClose(fd lin_common.FD_DEF, closeReason lin_common.EN_TCP_CLOSE_REASON, inAttachData interface{}) {

}
func (pthis*GameSrv)TcpOutBandData(fd lin_common.FD_DEF, data interface{}, inAttachData interface{}) {

}
func (pthis*GameSrv)TcpTick(fd lin_common.FD_DEF, tNowMill int64, inAttachData interface{}) {

}

func (pthis*GameSrv)GetCenterSrvUUID() (centerSrvUUID server_common.SRV_ID) {
	pbMsg := &msgpacket.PB_MSG_INTER_QUESRV_GET_SRVTYPE{}
	pbMsg.SrvType = int32(server_common.SRV_TYPE_center_server)
	pbMsgRes, err := pthis.mqClient.SendMsg(server_common.SRV_ID_INVALID, server_common.SRV_TYPE_msq_que,
		msgpacket.PB_MSG_TYPE__PB_MSG_INTER_QUESRV_GET_SRVTYPE, pbMsg,
		30 * 1000)

	if nil != err {
		lin_common.LogErr("get center server uuid err:", err)
		return server_common.SRV_ID_INVALID
	}

	pbRes, ok := pbMsgRes.(*msgpacket.PB_MSG_INTER_QUESRV_GET_SRVTYPE_RES)
	if !ok {
		lin_common.LogErr("get center server uuid err")
		return server_common.SRV_ID_INVALID
	}
	if nil == pbRes.ArrarySrv {
		lin_common.LogErr("get center server uuid err")
		return server_common.SRV_ID_INVALID
	}
	if 0 == len(pbRes.ArrarySrv) {
		lin_common.LogErr("get center server uuid err")
		return server_common.SRV_ID_INVALID
	}

	centerSrvUUID = server_common.SRV_ID(pbRes.ArrarySrv[0].SrvUuid)
	return
}

func (pthis*GameSrv)SendRegToCenter() {
	pbMsg := &msgpacket.PB_MSG_GAMESRV_CENTERSRV_REG{}
	tcpAddr, err := net.ResolveTCPAddr("tcp", pthis.outAddr)
	pbMsg.OutIp = tcpAddr.IP.String()
	pbMsg.OutPort = int32(tcpAddr.Port)
	_, err = pthis.mqClient.SendMsg(pthis.centerSrvUUID, server_common.SRV_TYPE_center_server,
		msgpacket.PB_MSG_TYPE__PB_MSG_GAMESRV_CENTERSRV_REG, pbMsg,
		30 * 1000)
	if nil != err {
		lin_common.LogErr("reg to center err:", err)
	}
}

func ConstructGameSrv(id string)*GameSrv {
	gs := &GameSrv{}

	gs.mqClient = msgque_client.ConstructMgrQueClient(server_common.Global_ServerCfg.MsgQueCent.OutAddr, server_common.SRV_TYPE_game_server, nil)
	gs.mqClient.DialToQueSrv()

	// get center server uuid
	gs.centerSrvUUID = gs.GetCenterSrvUUID()
	lin_common.LogInfo("center server:", gs.centerSrvUUID)

	gsCfg := server_common.GetGameSrvCfg(id)
	lin_common.LogInfo(gsCfg)

	gs.outAddr = gsCfg.OutAddr

	// create client tcp epoll listener
	var err error
	gs.lsn, err = lin_common.ConstructorEPollListener(gs, gsCfg.BindAddr, 10,
		lin_common.ParamEPollListener{
			ParamET: true,
			ParamEpollWaitTimeoutMills: 180 * 1000,
			ParamIdleClose: 600*1000,
			ParamNeedTick: true,
		})
	if err != nil {
		lin_common.LogErr("create epoll err")
		return nil
	}

	// send reg to center server
	gs.SendRegToCenter()

	return gs
}
