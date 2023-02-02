package main

import (
	"github.com/golang/protobuf/proto"
	"goserver/common"
	"goserver/msgpacket"
	"goserver/server/server_common"
)

func (pthis*LogonSrv)Go_CallBackMsg(pbMsg proto.Message, pbMsgType int32,
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


func (pthis*LogonSrv)Go_CallBackReport(pbLocal *msgpacket.PB_SRV_INFO_ALL) {
	common.LogDebug(pbLocal)

}


func (pthis*LogonSrv)Go_CallBackNtf(ntf * msgpacket.PB_MSG_INTER_QUESRV_NTF) {
	common.LogDebug(ntf)
}


func (pthis*LogonSrv)process_PB_MSG_LOGON(fd common.FD_DEF, pbMsg * msgpacket.PB_MSG_LOGON) {
	pbReq := &msgpacket.PB_MSG_LOGONSRV_CENTERSRV_LOGON{ClientId: pbMsg.ClientId}
	pbRes, err := pthis.mqClient.SendMsg(pthis.centerSrvUUID, server_common.SRV_TYPE_center_server,
		msgpacket.PB_MSG_TYPE__PB_MSG_LOGONSRV_CENTERSRV_LOGON, pbReq, 3*1000)
	if err != nil {
		common.LogErr(err)
		return
	}
	pbAddr := pbRes.(*msgpacket.PB_MSG_LOGONSRV_CENTERSRV_LOGON_RES)

	pbClientRes := &msgpacket.PB_MSG_LOGON_RES{
		ClientId: pbMsg.ClientId,
		Ip:       pbAddr.Ip,
		Port:     pbAddr.Port,
	}

	server_common.SendProtoMsg(pthis.lsn, fd, msgpacket.PB_MSG_TYPE__PB_MSG_LOGONSRV_CENTERSRV_LOGON_RES, pbClientRes)
}
