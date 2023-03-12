package main

import (
	"github.com/golang/protobuf/proto"
	"goserver/msgpacket"
	"goserver/server/server_common"
)

// new coroutine
func (pthis*DBSrv)Go_CallBackMsg(pbMsg proto.Message, pbMsgType int32,
	srvUUIDFrom server_common.SRV_ID,
	srvType server_common.SRV_TYPE,
	timeOutMills int) (msgType int32, protoMsg proto.Message) {

	switch msgpacket.PB_MSG_TYPE(pbMsgType) {
	case msgpacket.PB_MSG_TYPE__PB_MSG_DBSERVER_READ:
		msg, ok := pbMsg.(*msgpacket.PB_MSG_DBSERVER_READ)
		if ok {
			msgType, protoMsg = pthis.process_Go_CallBackMsg_PB_MSG_DBSERVER_READ(msg)
		}
	case msgpacket.PB_MSG_TYPE__PB_MSG_DBSERVER_WRITE:
	}

	return 0, nil
}

// new coroutine
func (pthis*DBSrv)Go_CallBackReport(pbLocal *msgpacket.PB_SRV_INFO_ALL) {

}

// new coroutine
func (pthis*DBSrv)Go_CallBackNtf(ntf *msgpacket.PB_MSG_INTER_QUESRV_NTF) {

}
