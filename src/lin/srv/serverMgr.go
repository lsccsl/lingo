package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"lin/log"
	"lin/msgpacket"
	"sync"
)

type MAP_CLIENT map[int64/*client id*/]*Client
type MAP_SERVER map[int64/*server id*/]*Server
type interProtoMsg struct {
	msgType msgpacket.MSG_TYPE
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

	heartbeatIntervalSec int
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

	switch msgpacket.MSG_TYPE(packType) {
	case msgpacket.MSG_TYPE__MSG_LOGIN:
		t, ok := protoMsg.(*msgpacket.MSG_LOGIN)
		if ok && t != nil {
			pthis.processClientLogin(t.Id, tcpConn)
		}

	case msgpacket.MSG_TYPE__MSG_SRV_REPORT:
		t, ok := protoMsg.(*msgpacket.MSG_SRV_REPORT)
		if ok && t != nil {
			if tcpConn.IsAccept {
				pthis.processSrvReport(tcpConn, t.SrvId)
			}
		}

	case msgpacket.MSG_TYPE__MSG_RPC:
		t, ok := protoMsg.(*msgpacket.MSG_RPC)
		if ok && t != nil {
			pthis.processRPCReq(tcpConn, t)
		}

	case msgpacket.MSG_TYPE__MSG_RPC_RES:
		t, ok := protoMsg.(*msgpacket.MSG_RPC_RES)
		if ok && t != nil {
			pthis.processRPCRes(tcpConn, t)
		}

	default:
		pthis.processMsg(tcpConn, msgpacket.MSG_TYPE(packType), protoMsg)
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

func ConstructServerMgr(srvID int64, heartbeatIntervalSec int) *ServerMgr {
	srvMgr := &ServerMgr{srvID: srvID}
	srvMgr.mapClient = make(MAP_CLIENT)
	srvMgr.mapServer = make(MAP_SERVER)
	srvMgr.heartbeatIntervalSec = heartbeatIntervalSec
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

	msgRes := &msgpacket.MSG_LOGIN_RES{}
	msgRes.Id = clientID
	msgRes.ConnectId = int64(tcpConn.TcpConnectionID())
	tcpConn.TcpConnectSendProtoMsg(msgpacket.MSG_TYPE__MSG_LOGIN_RES, msgRes)
}

func (pthis*ServerMgr)processMsg(tcpConn * TcpConnection, msgType msgpacket.MSG_TYPE, protoMsg proto.Message) {
	if tcpConn.SrvID != 0 {
		srv := pthis.getServer(tcpConn.SrvID)
		if srv != nil {
			srv.PushProtoMsg(msgType, protoMsg)
			return
		}
		pthis.tcpMgr.TcpMgrCloseConn(tcpConn.TcpConnectionID())
		return
	} else {
		cli := pthis.getClient(tcpConn.ClientID)
		if cli != nil {
			cli.PushClientMsg(msgType, protoMsg)
			return
		}
		pthis.tcpMgr.TcpMgrCloseConn(tcpConn.TcpConnectionID())
		return
	}
}



func (pthis*ServerMgr)ClientWriteProtoMsg(clientID int64, msgType msgpacket.MSG_TYPE, protoMsg proto.Message) {
	oldC := pthis.getClient(clientID)
	if oldC == nil {
		return
	}
	conn := oldC.ClientGetConnection()
	if  conn == nil {
		return
	}
	conn.TcpConnectSendProtoMsg(msgType, protoMsg)
}

func (pthis*ServerMgr)processSrvReport(tcpAccept * TcpConnection, srvID int64){
	tcpAccept.SrvID = srvID

	srv := pthis.getServer(srvID)
	if srv != nil {
		srv.PushInterMsg(&interMsgSrvReport{tcpAccept})
		return
	} else {
		srv = ConstructServer(pthis, srvID, pthis.heartbeatIntervalSec)
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
		srv = ConstructServer(pthis, srvID, pthis.heartbeatIntervalSec)
		pthis.addServer(srv)
		srv.PushInterMsg(&interMsgConnDial{tcpDial})
	}

	msgR := &msgpacket.MSG_SRV_REPORT{}
	msgR.SrvId = pthis.srvID
	tcpDial.TcpConnectSendProtoMsg(msgpacket.MSG_TYPE__MSG_SRV_REPORT, msgR)
}

func (pthis*ServerMgr)processRPCReq(tcpConn * TcpConnection, msg *msgpacket.MSG_RPC) {
	var msgRes proto.Message
	msgRPC := ParseProtoMsg(msg.MsgBin, msg.MsgType)
	if tcpConn.SrvID != 0 {
		srv := pthis.getServer(tcpConn.SrvID)
		if srv != nil {
			msgRes = srv.processRPC(tcpConn, msgRPC)
		} else {
			pthis.tcpMgr.TcpMgrCloseConn(tcpConn.TcpConnectionID())
			return
		}
	} else {
		cli := pthis.getClient(tcpConn.ClientID)
		if cli != nil {
			msgRes = cli.processRPC(tcpConn, msgRPC)
		} else {
			pthis.tcpMgr.TcpMgrCloseConn(tcpConn.TcpConnectionID())
			return
		}
	}

	msgRPCRes := &msgpacket.MSG_RPC_RES{MsgId:msg.MsgId, MsgType:msg.MsgType}
	if msgRes != nil {
		var err error
		msgRPCRes.MsgBin, err = proto.Marshal(msgRes)
		if err != nil {
			log.LogErr(err)
		}
	}
	tcpConn.TcpConnectSendProtoMsg(msgpacket.MSG_TYPE__MSG_RPC_RES, msgRPCRes)
}
func (pthis*ServerMgr)processRPCRes(tcpConn * TcpConnection, msg *msgpacket.MSG_RPC_RES) {
	msgRPC := ParseProtoMsg(msg.MsgBin, msg.MsgType)
	if tcpConn.SrvID != 0 {
		srv := pthis.getServer(tcpConn.SrvID)
		if srv != nil {
			srv.processRPCRes(tcpConn, msgRPC)
		}
	} else {
		cli := pthis.getClient(tcpConn.ClientID)
		if cli != nil {
			cli.processRPCRes(tcpConn, msgRPC)
		}
	}
}

func (pthis*ServerMgr)Dump() string {
	var str string
	str += "\r\nclient:\r\n"
	func(){
		pthis.ClientMapMgr.mapClientMutex.Lock()
		defer pthis.ClientMapMgr.mapClientMutex.Unlock()
		for _, val := range pthis.ClientMapMgr.mapClient {
			var connID TCP_CONNECTION_ID
			if val.tcpConn != nil {
				connID = val.tcpConn.TcpConnectionID()
			}
			str += fmt.Sprintf("\r\n client id:%v id:%v", val.clientID, connID)
		}
	}()

	str += "\r\nserver:\r\n"
	func(){
		pthis.ServerMapMgr.mapServerMutex.Lock()
		defer pthis.ServerMapMgr.mapServerMutex.Unlock()
		for _, val := range pthis.ServerMapMgr.mapServer {
			var acptID TCP_CONNECTION_ID
			var connID TCP_CONNECTION_ID
			if val.connAcpt != nil {
				acptID = val.connAcpt.TcpConnectionID()
			}
			if val.connDial != nil {
				connID = val.connDial.TcpConnectionID()
			}
			str += fmt.Sprintf("\r\n server id:%v acpt:%v dial:%v", val.srvID, acptID, connID)
		}
	}()

	str += "\r\ntcp connect:\r\n"
	func(){
		pthis.tcpMgr.mapConnMutex.Lock()
		defer pthis.tcpMgr.mapConnMutex.Unlock()
		for _, val := range pthis.tcpMgr.mapConn {
			str += fmt.Sprintf(" \r\n connection:%v remote:[%v] local:[%v] IsAccept:%v SrvID:%v ClientID:%v",
				val.TcpConnectionID(), val.netConn.RemoteAddr(), val.netConn.LocalAddr(), val.IsAccept, val.SrvID, val.ClientID)
		}
	}()

	str += "\r\ntcp dial data\r\n"
	func(){
		pthis.tcpMgr.TcpDialMgr.mapDialDataMutex.Lock()
		pthis.tcpMgr.TcpDialMgr.mapDialDataMutex.Unlock()
		for _, val := range pthis.tcpMgr.TcpDialMgr.mapDialData {
			var connID TCP_CONNECTION_ID
			if val.tcpConn != nil {
				connID = val.tcpConn.TcpConnectionID()
			}
			str += fmt.Sprintf("\r\n srvID:%v connection:%v [%v:%v]", val.srvID, connID, val.ip, val.port)
		}
	}()

	return str
}

