package main

import (
	"bytes"
	"encoding/binary"
	"lin/log"
)

type MAP_CLIENT map[int32]*Client

type SrvManager struct {
	mapClient MAP_CLIENT
}

func (pthis*SrvManager)CBReadProcess(tcpConn * TcpConnection, recvBuf * bytes.Buffer) (bytesProcess int) {
	log.LogDebug("len:", recvBuf.Len())

	binHead := recvBuf.Bytes()[0:6]

	packLen := binary.LittleEndian.Uint32(binHead[0:4])
	packType := binary.LittleEndian.Uint16(binHead[4:6])

	if recvBuf.Len() < int(packLen){
		return 0
	}

	binBody := recvBuf.Bytes()[6:packLen]

	protoMsg := ParseProtoMsg(binBody, int32(packType))
	log.LogDebug("parse protomsg:", protoMsg)

	return int(packLen)
}
func (pthis*SrvManager)CBConnect(tcpConn * TcpConnection) {
	if tcpConn == nil {
		return
	}
	log.LogDebug(tcpConn.TcpGetConn().LocalAddr(), tcpConn.TcpGetConn().RemoteAddr(), tcpConn.TcpClientID())
}

func (pthis*SrvManager)CBConnectClose(id int64) {
	log.LogDebug("id:", id)
}

func ConstructSrvManager() *SrvManager {
	srvMgr := &SrvManager{
		mapClient: make(MAP_CLIENT),
	}
	return srvMgr
}
