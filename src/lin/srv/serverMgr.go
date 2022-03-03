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
	tcpMgr *TcpMgr
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

	case msg.MSG_TYPE__MSG_RPC:
		t, ok := protoMsg.(*msg.MSG_RPC)
		if ok && t != nil {
			//???
		}

	case msg.MSG_TYPE__MSG_RPC_RES:
		t, ok := protoMsg.(*msg.MSG_RPC_RES)
		if ok && t != nil {
			//???
		}

	case msg.MSG_TYPE__MSG_SRV_REPORT:
		t, ok := protoMsg.(*msg.MSG_SRV_REPORT)
		if ok && t != nil {
			//???
		}

	default:
		pthis.processClient(tcpConn, msg.MSG_TYPE(packType), protoMsg)
	}

	return int(packLen)
}

func (pthis*Server)CBConnectAccept(tcpConn * TcpConnection, err error) {
	if err != nil {
		log.LogErr(err)
	}
	if tcpConn == nil {
		return
	}
	log.LogDebug(tcpConn.TcpGetConn().LocalAddr(), tcpConn.TcpGetConn().RemoteAddr(), tcpConn.TcpConnectionID())
}
func (pthis*Server)CBConnectDial(tcpConn * TcpConnection, err error) {
	if err != nil {
		log.LogErr(err)
	}
	if tcpConn == nil {
		return
	}
	log.LogDebug(tcpConn.TcpGetConn().LocalAddr(), tcpConn.TcpGetConn().RemoteAddr(), tcpConn.TcpConnectionID())
}

func (pthis*Server)CBConnectClose(tcpConn * TcpConnection) {
	log.LogDebug("id:", tcpConn.TcpConnectionID())
	if !tcpConn.IsAccept {
		pthis.tcpMgr.TcpDialMgrCheckReDial(tcpConn.SrvID)
	}
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
				pthis.tcpMgr.TcpMgrCloseConn(oldC.ClientGetConnection().TcpConnectionID())
			}
		}
	}

	// todo: add here
	c := ConstructClient(tcpConn, clientID)
	conn := c.ClientGetConnection()
	if conn == nil {
		return
	}
	tcpConn.ClientID = clientID

	pthis.mapClient[clientID] = c
}

func (pthis*Server)processClient(tcpConn * TcpConnection, msgType msg.MSG_TYPE, protoMsg proto.Message) {
	oldC, ok := pthis.mapClient[tcpConn.ClientID]
	if ok && oldC != nil {
		oldC.PushClientMsg(msgType, protoMsg)
		return
	}
	pthis.tcpMgr.TcpMgrCloseConn(tcpConn.TcpConnectionID())
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
