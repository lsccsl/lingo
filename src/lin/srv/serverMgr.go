package main

import (
	"bytes"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"lin/log"
	"lin/msg"
	"sync"
)

type MAP_CLIENT map[int64/*client id*/]*Client
type MAP_SERVER map[int64/*server id*/]*Server
type interMsg struct {
	msgType msg.MSG_TYPE
	protoMsg proto.Message
}

type ClientMapMgr struct {
	mapClientMutex sync.Mutex
	mapClient MAP_CLIENT
}
type ServerMgr struct {
	ClientMapMgr
	mapServer MAP_SERVER
	tcpMgr *TcpMgr
	httpSrv *HttpSrvMgr
}

func (pthis*ServerMgr)CBReadProcess(tcpConn * TcpConnection, recvBuf * bytes.Buffer) (bytesProcess int) {
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

	case msg.MSG_TYPE__MSG_SRV_REPORT:
		t, ok := protoMsg.(*msg.MSG_SRV_REPORT)
		if ok && t != nil {
			if tcpConn.IsAccept {
				pthis.processSrvReport(tcpConn, t.SrvId)
			}
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

	default:
		pthis.processClient(tcpConn, msg.MSG_TYPE(packType), protoMsg)
	}

	return int(packLen)
}

func (pthis*ServerMgr)CBConnectAccept(tcpConn * TcpConnection, err error) {
	if err != nil {
		log.LogErr(err)
	}
	if tcpConn == nil {
		return
	}
	log.LogDebug(tcpConn.TcpGetConn().LocalAddr(), tcpConn.TcpGetConn().RemoteAddr(), tcpConn.TcpConnectionID())
}
func (pthis*ServerMgr)CBConnectDial(tcpConn * TcpConnection, err error) {
	if err != nil {
		log.LogErr(err)
	}
	if tcpConn == nil {
		return
	}
	log.LogDebug(tcpConn.TcpGetConn().LocalAddr(), tcpConn.TcpGetConn().RemoteAddr(), tcpConn.TcpConnectionID())
}

func (pthis*ServerMgr)CBConnectClose(tcpConn * TcpConnection) {
	log.LogDebug("id:", tcpConn.TcpConnectionID())
	if !tcpConn.IsAccept {
		pthis.tcpMgr.TcpDialMgrCheckReDial(tcpConn.SrvID)
	}
}

func ConstructServerMgr() *ServerMgr {
	srvMgr := &ServerMgr{
		mapServer: make(MAP_SERVER),
	}
	srvMgr.mapClient = make(MAP_CLIENT)
	return srvMgr
}

func (pthis*ServerMgr)addClient(clientID int64, tcpConn * TcpConnection) {
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

func (pthis*ServerMgr)processClient(tcpConn * TcpConnection, msgType msg.MSG_TYPE, protoMsg proto.Message) {
	oldC, ok := pthis.mapClient[tcpConn.ClientID]
	if ok && oldC != nil {
		oldC.PushClientMsg(msgType, protoMsg)
		return
	}
	pthis.tcpMgr.TcpMgrCloseConn(tcpConn.TcpConnectionID())
}

func (pthis*ServerMgr)ClientWriteProtoMsg(clientID int64, msgType msg.MSG_TYPE, protoMsg proto.Message) {
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

func (pthis*ServerMgr)processSrvReport(tcpAccept * TcpConnection, srvID int64){
	tcpDial := pthis.tcpMgr.TcpDialGetSrvConn(srvID)
	if tcpDial == nil {
		return
	}

	srv := ConstructServer(tcpDial, tcpAccept)

	pthis.mapServer[srvID] = srv
}
/*
func (pthis*ServerMgr)processSrvReport(tcpAccept * TcpConnection, protoMsg *msg.MSG_SRV_REPORT){

}*/
