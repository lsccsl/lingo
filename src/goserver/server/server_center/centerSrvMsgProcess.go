package main

import (
	"github.com/golang/protobuf/proto"
	"goserver/common"
	"goserver/msgpacket"
	"goserver/server/server_common"
)


func (pthis*CenterSrv)Go_CallBackReport(pbLocal *msgpacket.PB_SRV_INFO_ALL) {
	// record all game server
	arrayGS := pthis.mqClient.GetCliSrvByType(server_common.SRV_TYPE_game_server)
	common.LogDebug("gs:", arrayGS)
	pthis.gmgr.SetGameSrv(arrayGS)

	for _, v := range arrayGS {
		pbRes, err := pthis.mqClient.SendMsg(v.SrvUUID, server_common.SRV_TYPE_game_server,
			msgpacket.PB_MSG_TYPE__PB_MSG_CENTERSRV_GAMESRV_GETINFO, &msgpacket.PB_MSG_CENTERSRV_GAMESRV_GETINFO{},
			3 * 1000)
		if err != nil {
			common.LogErr("get game srv out addr err:", err)
			continue
		}
		pbAddr := pbRes.(*msgpacket.PB_MSG_CENTERSRV_GAMESRV_GETINFO_RES)
		if nil == pbAddr {
			common.LogErr("get game srv out addr response nil")
			continue
		}
		pthis.gmgr.SetGameSrvOutAddr(v.SrvUUID, pbAddr.OutIp, int(pbAddr.OutPort))
		common.LogInfo(pbRes, err)
	}
}

func (pthis*CenterSrv)Go_CallBackNtf(ntf * msgpacket.PB_MSG_INTER_QUESRV_NTF) {
	common.LogDebug(ntf)
}

func (pthis*CenterSrv)Go_CallBackMsg(pbMsg proto.Message, pbMsgType int32,
	srvUUIDFrom server_common.SRV_ID,
	srvType server_common.SRV_TYPE,
	timeOutMills int) (msgType int32, protoMsg proto.Message) {


	msgType = int32(msgpacket.PB_MSG_TYPE__PB_MSG_INTER_NULL)
	protoMsg = nil

	switch t := pbMsg.(type) {
	default:
		{
			common.LogDebug("msg:", t, " pbMsgType:", pbMsgType, srvUUIDFrom, srvType, "timeOutMills:", timeOutMills)
		}
	}

	return
}

