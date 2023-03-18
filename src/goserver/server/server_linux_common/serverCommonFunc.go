package server_linux_common

import (
	"github.com/golang/protobuf/proto"
	"goserver/common"
	"goserver/msgpacket"
)

func SendProtoMsg(pthis*common.EPollListener, fd common.FD_DEF, msgType msgpacket.PB_MSG_TYPE, protoMsg proto.Message) {
	pthis.EPollListenerWrite(fd, msgpacket.ProtoPacketToBin(uint16(msgType), protoMsg))
}
