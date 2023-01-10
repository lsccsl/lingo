package main

import (
	"bytes"
	"lin/lin_common"
	"lin/msgpacket"
	"net"
)

type MsgqueMgr struct {
	lsn *lin_common.EPollListener

}


func (pthis*MsgqueMgr)TcpAcceptConnection(fd lin_common.FD_DEF, addr net.Addr, inAttachData interface{})(outAttachData interface{}) {
	lin_common.LogDebug(addr)
	return nil
}

func (pthis*MsgqueMgr)TcpDialConnection(fd lin_common.FD_DEF, addr net.Addr, inAttachData interface{})(outAttachData interface{}) {
	lin_common.LogDebug(addr)
	return nil
}

func (pthis*MsgqueMgr)TcpData(fd lin_common.FD_DEF, readBuf *bytes.Buffer, inAttachData interface{})(bytesProcess int, outAttachData interface{}) {
	lin_common.LogDebug(fd)
	packType, bytesProcess, protoMsg := msgpacket.ProtoUnPacketFromBin(readBuf)
	if protoMsg == nil {
		return
	}
	lin_common.LogInof("packType", packType)
	return 0,nil
}

func (pthis*MsgqueMgr)TcpClose(fd lin_common.FD_DEF, closeReason lin_common.EN_TCP_CLOSE_REASON, inAttachData interface{}) {
	lin_common.LogDebug(fd)

}

func (pthis*MsgqueMgr)TcpOutBandData(fd lin_common.FD_DEF, data interface{}, inAttachData interface{}) {
	lin_common.LogDebug(fd)

}


// ConstructMsgQueCenterMgr <addr> example 127.0.0.1:8888
func ConstructMsgqueMgr(addr string, epollCoroutineCount int) *MsgqueMgr{
	mqMgr := &MsgqueMgr{}

	lsn, err := lin_common.ConstructorEPollListener(mqMgr, addr, epollCoroutineCount, lin_common.ParamEPollListener{ParamET: true})
	if err != nil {
		lin_common.LogErr("constructor epoll listener err:", err)
		return nil
	}
	mqMgr.lsn = lsn

	return mqMgr
}

func (pthis*MsgqueMgr)MsgQueCenterMgrWait() {
	pthis.lsn.EPollListenerWait()
}