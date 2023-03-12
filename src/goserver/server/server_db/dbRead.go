package main

import (
	"github.com/golang/protobuf/proto"
	"goserver/msgpacket"
)

func (pthis*DBSrv)process_Go_CallBackMsg_PB_MSG_DBSERVER_READ(msg * msgpacket.PB_MSG_DBSERVER_READ) (msgType int32, protoMsg proto.Message) {
	msgRet := &msgpacket.PB_MSG_DBSERVER_READ_RES{}
	msgType = int32(msgpacket.PB_MSG_TYPE__PB_MSG_DBSERVER_READ_RES)
	protoMsg = msgRet

	// get db type by table name
	// read from db


	return
}
