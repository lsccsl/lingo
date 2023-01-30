package main

import (
	"lin/lin_common"
	"lin/server/server_common"
	"lin/server/server_msg_que_client"
)

type CenterSrv struct {
	mqClient *msgque_client.MgrQueClient
	gmgr * GameSrvMgr
}

func (pthis*CenterSrv)Wait() {
	pthis.mqClient.WaitEpoll()
}



func (pthis*CenterSrv)Go_CallBackProcessReport() {
	// record all game server
	arrayGS := pthis.mqClient.GetCliSrvByType(server_common.SRV_TYPE_game_server)
	lin_common.LogDebug("gs:", arrayGS)
	pthis.gmgr.SetGameSrv(arrayGS)
}

func ConstructCenterSrv()*CenterSrv {
	cs := &CenterSrv{
		gmgr:ConstructGameSrvMgr(),
	}

	cs.mqClient = msgque_client.ConstructMgrQueClient(server_common.Global_ServerCfg.MsgQueCent.OutAddr, server_common.SRV_TYPE_center_server, cs)
	cs.mqClient.DialToQueSrv()

	return cs
}

func (pthis*CenterSrv)Dump(bDetail bool) string {
	str := pthis.mqClient.Dump(bDetail)
	str += pthis.gmgr.Dump()
	return str
}