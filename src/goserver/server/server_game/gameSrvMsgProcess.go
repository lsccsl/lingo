package main

import (
	"github.com/golang/protobuf/proto"
	"goserver/common"
	"goserver/msgpacket"
	"goserver/server/server_common"
	"net"
)


func (pthis*GameSrv)Go_CallBackMsg(pbMsg proto.Message, pbMsgType int32,
	srvUUIDFrom server_common.SRV_ID,
	srvType server_common.SRV_TYPE,
	timeOutMills int) (msgType int32, protoMsg proto.Message) {

	msgType = int32(msgpacket.PB_MSG_TYPE__PB_MSG_INTER_NULL)
	protoMsg = nil

	switch t := pbMsg.(type) {
	case *msgpacket.PB_MSG_CENTERSRV_GAMESRV_GETINFO:
		{
			msgType, protoMsg = pthis.process_PB_MSG_CENTERSRV_GAMESRV_GETINFO()
		}

	default:
		{
			common.LogDebug("msg:", t, " pbMsgType:", pbMsgType, srvUUIDFrom, srvType, "timeOutMills:", timeOutMills)
		}
	}

	return
}


func (pthis*GameSrv)Go_CallBackReport(pbLocal *msgpacket.PB_SRV_INFO_ALL) {
	common.LogDebug(pbLocal)

	pthis.GetCenterSrvUUID()
}


func (pthis*GameSrv)Go_CallBackNtf(ntf * msgpacket.PB_MSG_INTER_QUESRV_NTF) {
	common.LogDebug(ntf)
}

func (pthis*GameSrv)process_PB_MSG_CENTERSRV_GAMESRV_GETINFO()(msgType int32, protoMsg proto.Message) {

	msgType = int32(msgpacket.PB_MSG_TYPE__PB_MSG_CENTERSRV_GAMESRV_GETINFO_RES)
	pbRes := &msgpacket.PB_MSG_CENTERSRV_GAMESRV_GETINFO_RES{}
	protoMsg = pbRes

	tcpAddr, err := net.ResolveTCPAddr("tcp", pthis.outAddr)
	if err != nil {
		common.LogErr(err)
	}
	pbRes.OutIp = tcpAddr.IP.String()
	pbRes.OutPort = int32(tcpAddr.Port)

	return
}