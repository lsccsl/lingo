package main

import (
	"bytes"
	"lin/lin_common"
	"net"
)

type MsgQueCenterMgr struct {
	lsn *lin_common.EPollListener
}

func (pthis*MsgQueCenterMgr)TcpAcceptConnection(fd lin_common.FD_DEF, addr net.Addr, inAttachData interface{})(outAttachData interface{}) {
	lin_common.LogDebug(addr)
	return nil
}

func (pthis*MsgQueCenterMgr)TcpDialConnection(fd lin_common.FD_DEF, addr net.Addr, inAttachData interface{})(outAttachData interface{}) {
	lin_common.LogDebug(addr)
	return nil
}

func (pthis*MsgQueCenterMgr)TcpData(fd lin_common.FD_DEF, readBuf *bytes.Buffer, inAttachData interface{})(bytesProcess int, outAttachData interface{}) {
	lin_common.LogDebug(fd)
	return 0,nil
}

func (pthis*MsgQueCenterMgr)TcpClose(fd lin_common.FD_DEF, closeReason lin_common.EN_TCP_CLOSE_REASON, inAttachData interface{}) {
	lin_common.LogDebug(fd)

}

func (pthis*MsgQueCenterMgr)TcpOutBandData(fd lin_common.FD_DEF, data interface{}, inAttachData interface{}) {
	lin_common.LogDebug(fd)

}

// ConstructMsgQueCenterMgr <addr> example 127.0.0.1:8888
func ConstructMsgQueCenterMgr(addr string, epollCoroutineCount int) *MsgQueCenterMgr{
	mqMgr := &MsgQueCenterMgr{}

	lsn, err := lin_common.ConstructorEPollListener(mqMgr, addr, epollCoroutineCount, lin_common.ParamEPollListener{ParamET: true})
	if err != nil {
		lin_common.LogErr("constructor epoll listener err:", err)
		return nil
	}
	mqMgr.lsn = lsn

	return mqMgr
}

func (pthis*MsgQueCenterMgr)MsgQueCenterMgrWait() {
	pthis.lsn.EPollListenerWait()
}
