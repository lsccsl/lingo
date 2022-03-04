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
type interProtoMsg struct {
	msgType msg.MSG_TYPE
	protoMsg proto.Message
}

type ClientMapMgr struct {
	mapClientMutex sync.Mutex
	mapClient MAP_CLIENT
}
type ServerMapMgr struct {
	mapServerMutex sync.Mutex
	mapServer MAP_SERVER
}
type ServerMgr struct {
	srvID int64
	ClientMapMgr
	ServerMapMgr
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
			pthis.processClientLogin(t.Id, tcpConn)
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

	pthis.processDailConnect(tcpConn)
}

func (pthis*ServerMgr)CBConnectClose(tcpConn * TcpConnection) {
	log.LogDebug("id:", tcpConn.TcpConnectionID(), " is accept:", tcpConn.IsAccept)
	if !tcpConn.IsAccept {
		pthis.delServer(tcpConn.SrvID)
		pthis.tcpMgr.TcpDialMgrCheckReDial(tcpConn.SrvID)
	} else {
		if tcpConn.SrvID != 0 {
			pthis.delServer(tcpConn.SrvID)
		}
	}
}

func ConstructServerMgr(srvID int64) *ServerMgr {
	srvMgr := &ServerMgr{srvID: srvID}
	srvMgr.mapClient = make(MAP_CLIENT)
	srvMgr.mapServer = make(MAP_SERVER)
	return srvMgr
}


func (pthis*ServerMgr)getClient(clientID int64) *Client {
	pthis.ClientMapMgr.mapClientMutex.Lock()
	defer pthis.ClientMapMgr.mapClientMutex.Unlock()

	oldC, _ := pthis.ClientMapMgr.mapClient[clientID]
	return oldC
}
func (pthis*ServerMgr)addClient(c *Client) {
	pthis.ClientMapMgr.mapClientMutex.Lock()
	defer pthis.ClientMapMgr.mapClientMutex.Unlock()

	pthis.ClientMapMgr.mapClient[c.clientID] = c
}
func (pthis*ServerMgr)delClient(clientID int64) {
	pthis.ClientMapMgr.mapClientMutex.Lock()
	defer pthis.ClientMapMgr.mapClientMutex.Unlock()

	oldC, _ := pthis.ClientMapMgr.mapClient[clientID]
	if  oldC != nil {
		oldC.ClientClose()
	}
	delete(pthis.ClientMapMgr.mapClient, clientID)
}


func (pthis*ServerMgr)getServer(srvID int64) *Server {
	pthis.ServerMapMgr.mapServerMutex.Lock()
	defer pthis.ServerMapMgr.mapServerMutex.Unlock()

	oldS, _ := pthis.ServerMapMgr.mapServer[srvID]
	return oldS
}
func (pthis*ServerMgr)addServer(s *Server) {
	pthis.ServerMapMgr.mapServerMutex.Lock()
	defer pthis.ServerMapMgr.mapServerMutex.Unlock()

	pthis.ServerMapMgr.mapServer[s.srvID] = s
}
func (pthis*ServerMgr)delServer(srvID int64) {
	pthis.ServerMapMgr.mapServerMutex.Lock()
	defer pthis.ServerMapMgr.mapServerMutex.Unlock()

	oldS, _ := pthis.ServerMapMgr.mapServer[srvID]
	if oldS != nil {
		oldS.ServerClose()
	}
	delete(pthis.ServerMapMgr.mapServer, srvID)
}

func (pthis*ServerMgr)processClientLogin(clientID int64, tcpConn * TcpConnection) {
	if tcpConn == nil {
		return
	}

	tcpConn.ClientID = clientID

	oldC := pthis.getClient(clientID)
	if oldC != nil {
		conn := oldC.ClientGetConnection()
		if conn != nil {
			if conn.TcpConnectionID() != tcpConn.TcpConnectionID() {
				pthis.delClient(clientID)
			}
		}
	}

	c := ConstructClient(pthis, tcpConn, clientID)
	pthis.addClient(c)

	msgRes := &msg.MSG_LOGIN_RES{}
	msgRes.Id = clientID
	msgRes.ConnectId = int64(tcpConn.TcpConnectionID())
	tcpConn.TcpConnectWriteProtoMsg(msg.MSG_TYPE__MSG_LOGIN_RES, msgRes)
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
	tcpAccept.SrvID = srvID

	srv := pthis.getServer(srvID)
	if srv != nil {
		srv.PushInterMsg(&interMsgSrvReport{tcpAccept})
		return
	} else {
		srv = ConstructServer(pthis, srvID)
		pthis.addServer(srv)
		srv.PushInterMsg(&interMsgSrvReport{tcpAccept})
		return
	}
}

func (pthis*ServerMgr)processDailConnect(tcpDial * TcpConnection){
	srvID := tcpDial.SrvID
	srv := pthis.getServer(srvID)
	if srv != nil {
		srv.PushInterMsg(&interMsgConnDial{tcpDial})
	} else {
		srv = ConstructServer(pthis, srvID)
		pthis.addServer(srv)
		srv.PushInterMsg(&interMsgConnDial{tcpDial})
	}

	msgR := &msg.MSG_SRV_REPORT{}
	msgR.SrvId = pthis.srvID
	tcpDial.TcpConnectWriteProtoMsg(msg.MSG_TYPE__MSG_SRV_REPORT, msgR)
}

// todo:多了一个accept dump all mem data
