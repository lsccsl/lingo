package main

import (
	"bytes"
	"lin/lin_common"
	"lin/msgpacket"
	"net"
)

type MsgqueCenterSrv struct {
	lsn *lin_common.EPollListener
}

func (pthis*MsgqueCenterSrv)TcpAcceptConnection(fd lin_common.FD_DEF, addr net.Addr, inAttachData interface{})(outAttachData interface{}) {
	lin_common.LogDebug(addr)
	return nil
}

func (pthis*MsgqueCenterSrv)TcpDialConnection(fd lin_common.FD_DEF, addr net.Addr, inAttachData interface{})(outAttachData interface{}) {
	lin_common.LogDebug(addr)
	return nil
}

func (pthis*MsgqueCenterSrv)TcpData(fd lin_common.FD_DEF, readBuf *bytes.Buffer, inAttachData interface{})(bytesProcess int, outAttachData interface{}) {
	lin_common.LogDebug(fd)
	packType, bytesProcess, protoMsg := msgpacket.ProtoUnPacketFromBin(readBuf)
	if protoMsg == nil {
		return
	}
	lin_common.LogInfo("packType:", packType, " bytesProcess:", bytesProcess)

	return 0,nil
}

func (pthis*MsgqueCenterSrv)TcpClose(fd lin_common.FD_DEF, closeReason lin_common.EN_TCP_CLOSE_REASON, inAttachData interface{}) {
	lin_common.LogDebug(fd)

}

func (pthis*MsgqueCenterSrv)TcpOutBandData(fd lin_common.FD_DEF, data interface{}, inAttachData interface{}) {
	lin_common.LogDebug(fd)

}

// ConstructMsgQueCenterMgr <addr> example 127.0.0.1:8888
func ConstructMsgqueCenterSrv(addr string, epollCoroutineCount int) *MsgqueCenterSrv{
	mqMgr := &MsgqueCenterSrv{}

	lsn, err := lin_common.ConstructorEPollListener(mqMgr, addr, epollCoroutineCount, lin_common.ParamEPollListener{ParamET: true})
	if err != nil {
		lin_common.LogErr("constructor epoll listener err:", err)
		return nil
	}
	mqMgr.lsn = lsn

	return mqMgr
}

func (pthis*MsgqueCenterSrv)MsgqueCenterSrvWait() {
	pthis.lsn.EPollListenerWait()
}
