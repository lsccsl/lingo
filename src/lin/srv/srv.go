package main

import (
	"bytes"
	"lin/log"
)

type MAP_CLIENT map[int64]* TcpConnection

type SrvManager struct {
	mapClient MAP_CLIENT
}

func (pthis*SrvManager)CBReadProcess(tcpConn * TcpConnection, recvBuf * bytes.Buffer)(bytesProcess int){
	log.LogDebug("len:", recvBuf.Len())
	return recvBuf.Len()
}
func (pthis*SrvManager)CBConnect(tcpConn * TcpConnection){
	if tcpConn == nil {
		return
	}
	log.LogDebug(tcpConn.TcpGetConn().LocalAddr(), tcpConn.TcpGetConn().RemoteAddr(), tcpConn.TcpClientID())

	pthis.mapClient[tcpConn.TcpClientID()] = tcpConn
}

func (pthis*SrvManager)CBConnectClose(id int64){
	log.LogDebug("id:", id)

	delete(pthis.mapClient, id)
}

func ConstructSrvManager() *SrvManager {
	srvMgr := &SrvManager{
		mapClient: make(MAP_CLIENT),
	}
	return srvMgr
}
