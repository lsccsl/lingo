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

	return 0, nil
}

// new coroutine
func (pthis*DBSrv)Go_CallBackReport(pbLocal *msgpacket.PB_SRV_INFO_ALL) {

}

// new coroutine
func (pthis*DBSrv)Go_CallBackNtf(ntf *msgpacket.PB_MSG_INTER_QUESRV_NTF) {

}
