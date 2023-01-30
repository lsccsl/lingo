package main

import (
	"github.com/golang/protobuf/proto"
	"lin/lin_common"
	"lin/msgpacket"
	"lin/server/server_common"
)

func (pthis*CenterSrv)Go_CallBackProcessMsg(pbMsg proto.Message, pbMsgType int32,
	srvUUIDFrom server_common.SRV_ID,
	srvType server_common.SRV_TYPE,
	timeOutMills int) (msgType int32, protoMsg proto.Message) {


	msgType = int32(msgpacket.PB_MSG_TYPE__PB_MSG_INTER_NULL)
	protoMsg = nil

	switch msgpacket.PB_MSG_TYPE(pbMsgType) {
	case msgpacket.PB_MSG_TYPE__PB_MSG_GAMESRV_CENTERSRV_REG:
		{
			msgType, protoMsg = pthis.process_PB_MSG_TYPE__PB_MSG_GAMESRV_CENTERSRV_REG(srvUUIDFrom, pbMsg)
		}

	default:
		{
			lin_common.LogDebug("msg:", pbMsg, " pbMsgType:", pbMsgType, srvUUIDFrom, srvType, "timeOutMills:", timeOutMills)
		}
	}

	return
}

func (pthis*CenterSrv)process_PB_MSG_TYPE__PB_MSG_GAMESRV_CENTERSRV_REG(srvUUIDFrom server_common.SRV_ID, pbMsg proto.Message) (msgType int32, protoMsg proto.Message) {
	msgType = int32(msgpacket.PB_MSG_TYPE__PB_MSG_GAMESRV_CENTERSRV_REG_RES)
	pbRes := &msgpacket.PB_MSG_GAMESRV_CENTERSRV_REG_RES{}
	protoMsg = pbRes

	pbReg, ok := pbMsg.(*msgpacket.PB_MSG_GAMESRV_CENTERSRV_REG)
	if !ok || nil == pbReg {
		lin_common.LogErr("convert msg err,", pbReg, pbMsg)
		return
	}

	lin_common.LogDebug(srvUUIDFrom, "[", pbReg.OutIp, ":", pbReg.OutPort, "]")

	pthis.gmgr.SetGameSrvOutAddr(srvUUIDFrom, pbReg.OutIp, int(pbReg.OutPort))

	return
}
