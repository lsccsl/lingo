package main

import (
	msgque_client "goserver/server/server_msg_que_client"
	"goserver/common"
	"goserver/msgpacket"
	"goserver/server/server_common"
)

type CenterSrv struct {
	mqClient *msgque_client.MgrQueClient
	gmgr *GameSrvMgr
}

func (pthis*CenterSrv)Wait() {
	pthis.mqClient.WaitEpoll()
}



func (pthis*CenterSrv)Go_CallBackReport(pbLocal *msgpacket.PB_SRV_INFO_ALL) {
	// record all game server
	arrayGS := pthis.mqClient.GetCliSrvByType(server_common.SRV_TYPE_game_server)
	common.LogDebug("gs:", arrayGS)
	pthis.gmgr.SetGameSrv(arrayGS)

	for _, v := range arrayGS {
		pthis.mqClient.SendMsg(v.SrvUUID, server_common.SRV_TYPE_game_server,
			msgpacket.PB_MSG_TYPE__PB_MSG_CENTERSRV_GAMESRV_GETINFO, nil,
			3 * 1000)
	}
}

func (pthis*CenterSrv)Go_CallBackNtf(ntf * msgpacket.PB_MSG_INTER_QUESRV_NTF) {
	common.LogDebug(ntf)
}

func ConstructCenterSrv()*CenterSrv {
	cs := &CenterSrv{
		gmgr: ConstructGameSrvMgr(),
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