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
	dialMgr *TcpDialMgr
	httpSrv *HttpSrvMgr
}

func (pthis*Server)CBReadProcess(tcpConn * TcpConnection, recvBuf * bytes.Buffer) (bytesProcess int) {
	if recvBuf.Len() < 6 {
		return 0
	}
	binHead := recvBuf.Bytes()[0:6]

	packLen := binary.LittleEndian.Uint32(binHead[0:4])
	packType := binary.LittleEndian.Uint16(binHead[4:6])

	log.LogDebug("packLen:", packLen, " packType:", packType)

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

			msgRes := &msg.MSG_LOGIN_RES{}
			msgRes.Id = t.Id
			msgRes.ConnectId = int64(tcpConn.TcpConnectionID())
			tcpConn.TcpConnectWriteProtoMsg(msg.MSG_TYPE__MSG_LOGIN_RES, msgRes)
		}
	default:
		pthis.processClient(tcpConn, msg.MSG_TYPE(packType), protoMsg)
	}

	return int(packLen)
}
func (pthis*Server)CBConnect(tcpConn * TcpConnection, err error) {
	if err != nil {
		log.LogErr(err)
	}
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
	oldC, ok := pthis.mapClient[clientID]
	if ok && oldC != nil {
		conn := oldC.ClientGetConnection()
		if conn != nil {
			if conn.TcpConnectionID() != tcpConn.TcpConnectionID() {
				pthis.accept.TcpAcceptCloseConn(oldC.ClientGetConnection().TcpConnectionID())
			}
		}
	}

	// todo: add here
	c := ConstructClient(tcpConn, clientID)
	conn := c.ClientGetConnection()
	if conn == nil {
		return
	}
	tcpConn.AppID = clientID

	pthis.mapClient[clientID] = c
}

func (pthis*Server)processClient(tcpConn * TcpConnection, msgType msg.MSG_TYPE, protoMsg proto.Message) {
	oldC, ok := pthis.mapClient[tcpConn.AppID]
	if ok && oldC != nil {
		oldC.PushClientMsg(msgType, protoMsg)
		return
	}
	pthis.accept.TcpAcceptCloseConn(tcpConn.TcpConnectionID())
}

func (pthis*Server)ClientWriteProtoMsg(clientID int64, msgType msg.MSG_TYPE, protoMsg proto.Message) {
	oldC, ok := pthis.mapClient[clientID]
	if !ok || oldC == nil {
		return
	}
	conn := oldC.ClientGetConnection()
	if  conn == nil {
		return
	}
	conn.TcpConnectWriteProtoMsg(msgType, protoMsg)
}
