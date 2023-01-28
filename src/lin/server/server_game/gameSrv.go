package main

import (
	"lin/msgpacket"
	"lin/server/server_common"
	"lin/server/server_msg_que_client"
)

type GameSrv struct {
	mqClient *msgque_client.MgrQueClient
}

func (pthis*GameSrv)Wait() {
	pthis.mqClient.Wait()
}

func ConstructGameSrv()*GameSrv {
	gs := &GameSrv{}

	gs.mqClient = msgque_client.ConstructMgrQueClient(server_common.Global_ServerCfg.MsgQueCent.OutAddr, server_common.SRV_TYPE_game_server)

	pbMsg := &msgpacket.PB_MSG_INTER_QUESRV_GET_SRVTYPE{}
	pbMsg.SrvType = int32(server_common.SRV_TYPE_center_server)
	gs.mqClient.SendMsgAsyn(server_common.SRV_ID_INVALID, server_common.SRV_TYPE_msq_que,
		msgpacket.PB_MSG_TYPE__PB_MSG_INTER_QUESRV_GET_SRVTYPE, pbMsg,
		30 * 1000)

	return gs
}