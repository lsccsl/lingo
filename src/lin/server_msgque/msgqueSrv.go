package main

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"lin/lin_common"
	"lin/msgpacket"
	"net"
)

type MsgQueSrv struct {
	lsn *lin_common.EPollListener

	addrOut string
}


func (pthis*MsgQueSrv)TcpAcceptConnection(fd lin_common.FD_DEF, addr net.Addr, inAttachData interface{})(outAttachData interface{}) {
	lin_common.LogDebug(addr)
	return nil
}

func (pthis*MsgQueSrv)TcpDialConnection(fd lin_common.FD_DEF, addr net.Addr, inAttachData interface{})(outAttachData interface{}) {
	lin_common.LogDebug(addr)

	tcpAddr, err := net.ResolveTCPAddr("tcp", pthis.addrOut)
	if err != nil {
		lin_common.LogErr(err)
	}
	pbMsgReg := &msgpacket.PB_MSG_INTER_QUESRV_REGISTER{
		Ip: tcpAddr.IP.String(),
		Port: int32(tcpAddr.Port),
	}
	pthis.SendProtoMsg(fd, msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_REGISTER, pbMsgReg)

	return nil
}

func (pthis*MsgQueSrv)TcpData(fd lin_common.FD_DEF, readBuf *bytes.Buffer, inAttachData interface{})(bytesProcess int, outAttachData interface{}) {
	lin_common.LogDebug(fd)
	packType, bytesProcess, protoMsg := msgpacket.ProtoUnPacketFromBin(readBuf)
	if protoMsg == nil {
		return
	}
	lin_common.LogInfo("packType:", packType, " bytesProcess:", bytesProcess, " proto msg", protoMsg, "")

	switch msgpacket.PB_MSG_INTER_TYPE(packType) {
	case msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_REGISTER_RES:
		{
		}
	}

	return bytesProcess,nil
}

func (pthis*MsgQueSrv)TcpClose(fd lin_common.FD_DEF, closeReason lin_common.EN_TCP_CLOSE_REASON, inAttachData interface{}) {
	lin_common.LogDebug(fd)

}

func (pthis*MsgQueSrv)TcpOutBandData(fd lin_common.FD_DEF, data interface{}, inAttachData interface{}) {
	lin_common.LogDebug(fd)
}



func (pthis*MsgQueSrv)SendProtoMsg(fd lin_common.FD_DEF, msgType msgpacket.PB_MSG_INTER_TYPE, protoMsg proto.Message){
	pthis.lsn.EPollListenerWrite(fd, msgpacket.ProtoPacketToBin(uint16(msgType), protoMsg))
}



func (pthis*MsgQueSrv)processRegRes() {

}

// ConstructMsgQueSrv <addr> example 127.0.0.1:8888
func ConstructMsgQueSrv(msgqueCenterAddr string, addrBind string, addrOut string, epollCoroutineCount int) *MsgQueSrv{
	mqMgr := &MsgQueSrv{
		addrOut : addrOut,
	}

	lsn, err := lin_common.ConstructorEPollListener(mqMgr, addrBind, epollCoroutineCount, lin_common.ParamEPollListener{ParamET: true})
	if err != nil {
		lin_common.LogErr("constructor epoll listener err:", err)
		return nil
	}
	mqMgr.lsn = lsn

	//连接入msg que center,连接成功后请求数据包
	lsn.EPollListenerDial(msgqueCenterAddr, nil)

	return mqMgr
}

func (pthis*MsgQueSrv)MsgQueSrvWait() {
	pthis.lsn.EPollListenerWait()
}