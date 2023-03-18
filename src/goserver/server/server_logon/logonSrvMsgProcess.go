package main

import (
	"github.com/golang/protobuf/proto"
	"goserver/common"
	"goserver/msgpacket"
	"goserver/server/server_common"
	"goserver/server/server_linux_common"
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
	pthis.GetCenterSrvUUID()
}


func (pthis*LogonSrv)Go_CallBackNtf(ntf * msgpacket.PB_MSG_INTER_QUESRV_NTF) {
	common.LogDebug(ntf)
}


func (pthis*LogonSrv)process_PB_MSG_LOGON(fd common.FD_DEF, pbMsg * msgpacket.PB_MSG_LOGON) {

	// todo call auth client here

	// 先查询是否存在
	msgRead := &msgpacket.PB_MSG_DBSERVER_READ{
		DatabaseAppName:"user",
		TableName:"DBUserMain",
	}
	dbKey := &msgpacket.DBUserMain{
		UserId:pbMsg.ClientId,
	}
	msgRead.Key, _ = proto.Marshal(dbKey)

	// pthis.mqClient.SendMsgToSrvUUID()
	// 不存在添加

	pbReq := &msgpacket.PB_MSG_LOGONSRV_CENTERSRV_LOGON{ClientId: pbMsg.ClientId}
	pbRes, err := pthis.mqClient.SendMsgToSrvUUID(pthis.centerSrvUUID,
		msgpacket.PB_MSG_TYPE__PB_MSG_LOGONSRV_CENTERSRV_LOGON, pbReq)
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

	server_linux_common.SendProtoMsg(pthis.lsn, fd, msgpacket.PB_MSG_TYPE__PB_MSG_LOGON_RES, pbClientRes)
}
