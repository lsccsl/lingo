package main

import (
	"github.com/golang/protobuf/proto"
	"lin/lin_common"
	"lin/msgpacket"
	"lin/server/server_common"
	"lin/server/server_msg_que_client"
)

type CenterSrv struct {
	mqClient *msgque_client.MgrQueClient
	gmgr * GameSrvMgr
}

func (pthis*CenterSrv)Wait() {
	pthis.mqClient.Wait()
}

func (pthis*CenterSrv)Go_CallBackProcessMsg(pbMsg proto.Message, pbMsgType int32,
	srvUUIDFrom server_common.SRV_ID,
	srvType server_common.SRV_TYPE,
	timeOutMills int) (msgType int32, protoMsg proto.Message) {

	lin_common.LogDebug("msg:", pbMsg, " pbMsgType:", pbMsgType, srvUUIDFrom, srvType, "timeOutMills:", timeOutMills)

	msgType = int32(msgpacket.PB_MSG_TYPE__PB_MSG_INTER_NULL)
	protoMsg = nil

	switch msgpacket.PB_MSG_TYPE(pbMsgType) {

	}

	return
}

func (pthis*CenterSrv)Go_CallBackProcessReport() {
	// record all game server
	lin_common.LogDebug("~~~~~~ Go_CallBackProcessReport ~~~~~~~~~~~~~")
}

func ConstructCenterSrv()*CenterSrv {
	cs := &CenterSrv{
		gmgr:ConstructGameSrvMgr(),
	}

	cs.mqClient = msgque_client.ConstructMgrQueClient(server_common.Global_ServerCfg.MsgQueCent.OutAddr, server_common.SRV_TYPE_center_server, cs)

	return cs
}

func (pthis*CenterSrv)Dump(bDetail bool) string {
	str := pthis.mqClient.Dump(bDetail)
	return str
}