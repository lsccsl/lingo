package main

import (
	"bytes"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"lin/log"
	"lin/msg"
)

type MAP_CLIENT map[int64/*client id*/]*Client

type Server struct {
	mapClient MAP_CLIENT
	accept *TcpAccept
}

func (pthis*Server)CBReadProcess(tcpConn * TcpConnection, recvBuf * bytes.Buffer) (bytesProcess int) {
	log.LogDebug("len:", recvBuf.Len())

	binHead := recvBuf.Bytes()[0:6]

	packLen := binary.LittleEndian.Uint32(binHead[0:4])
	packType := binary.LittleEndian.Uint16(binHead[4:6])

	if recvBuf.Len() < int(packLen){
		return 0
	}

	binBody := recvBuf.Bytes()[6:packLen]

	protoMsg := ParseProtoMsg(binBody, int32(packType))
	if protoMsg == nil {
		return int(packLen)
	}
	log.LogDebug("parse protomsg:", protoMsg)

/*	switch t:=protoMsg.(type) {
	case *msg.MSG_LOGIN:
		addClient(t.Id, tcpConn)
	default:
	}*/

	switch msg.MSG_TYPE(packType) {
	case msg.MSG_TYPE__MSG_LOGIN:
		t, ok := protoMsg.(*msg.MSG_LOGIN)
		if ok && t != nil {
			pthis.addClient(t.Id, tcpConn)
		}
	default:
		pthis.processClient(tcpConn, protoMsg)
	}

	return int(packLen)
}
func (pthis*Server)CBConnect(tcpConn * TcpConnection) {
	if tcpConn == nil {
		return
	}
	log.LogDebug(tcpConn.TcpGetConn().LocalAddr(), tcpConn.TcpGetConn().RemoteAddr(), tcpConn.TcpConnectionID())
}

func (pthis*Server)CBConnectClose(id TCP_CONNECTION_ID) {
	log.LogDebug("id:", id)
}

func ConstructServer() *Server {
	server := &Server{
		mapClient: make(MAP_CLIENT),
	}
	return server
}

func (pthis*Server)addClient(clientID int64, tcpConn * TcpConnection) {
	// todo: add auth here
	c := ConstructClient(tcpConn.TcpConnectionID(), clientID)
	tcpConn.TcpConnectSetClientAppID(clientID)

	oldC, ok := pthis.mapClient[clientID]
	if ok && oldC != nil {
		pthis.accept.TcpAcceptCloseConn(oldC.ClientGetConnectID())
	}

	pthis.mapClient[clientID] = c
}

func (pthis*Server)processClient(tcpConn * TcpConnection, protoMsg proto.Message) {
	oldC, ok := pthis.mapClient[tcpConn.TcpConnectClientAppID()]
	if ok && oldC != nil {
		oldC.ClientProcess(protoMsg)
		return
	}
	pthis.accept.TcpAcceptCloseConn(tcpConn.TcpConnectionID())
}
