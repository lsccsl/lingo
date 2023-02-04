package main

import (
	"github.com/golang/protobuf/proto"
	"goserver/common"
	"goserver/msgpacket"
	"goserver/server/server_common"
)

/*// 像声明接口一样声明
type MyInt interface {
	int | int8 | int16 | int32 | int64
}

// T的类型为声明的MyInt
func GetMaxNum[T MyInt](a, b T) T {
	if a > b {
		return a
	}

	return b
}*/

func (pthis*CenterSrv)Go_CallBackReport(pbLocal *msgpacket.PB_SRV_INFO_ALL) {

	// record all game server
	arrayGS := pthis.mqClient.GetCliSrvByType(server_common.SRV_TYPE_game_server)
	common.LogDebug("gs:", arrayGS)
	pthis.gmgr.SetGameSrv(arrayGS)

	//var pbAddr *msgpacket.PB_MSG_CENTERSRV_GAMESRV_GETINFO_RES
	for _, v := range arrayGS {
		pbRes, err := pthis.mqClient.SendMsgToSrvUUID(v.SrvUUID,
			msgpacket.PB_MSG_TYPE__PB_MSG_CENTERSRV_GAMESRV_GETINFO, &msgpacket.PB_MSG_CENTERSRV_GAMESRV_GETINFO{})
		if err != nil {
			common.LogErr("get game srv out addr err:", err)
			continue
		}
		pbAddr, ok := pbRes.(*msgpacket.PB_MSG_CENTERSRV_GAMESRV_GETINFO_RES)
		if nil == pbAddr || !ok {
			common.LogErr("get game srv out addr response nil")
			continue
		}
		pthis.gmgr.SetGameSrvOutAddr(v.SrvUUID, pbAddr.OutIp, int(pbAddr.OutPort))
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

	case *msgpacket.PB_MSG_LOGONSRV_CENTERSRV_LOGON:
		{
			//choose a game server
			gi := pthis.gmgr.GetGameSrv()
			msgType = int32(msgpacket.PB_MSG_TYPE__PB_MSG_LOGONSRV_CENTERSRV_LOGON_RES)
			protoMsg = &msgpacket.PB_MSG_LOGONSRV_CENTERSRV_LOGON_RES{
				ClientId: t.ClientId,
				GameSrvUuid:int64(gi.SrvUUID),
				Ip:gi.outIP,
				Port:int32(gi.outPort),
			}
		}

	default:
		{
			common.LogDebug("msg:", t, " pbMsgType:", pbMsgType, srvUUIDFrom, srvType, "timeOutMills:", timeOutMills)
		}
	}

	return
}

