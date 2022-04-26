package main

import (
	"lin/lin_common"
	"lin/msgpacket"
)

type MAP_TCPCLIENT map[int]*TcpClient
type MAP_CLIENTID_TCPFD map[int64]int

type ClientProcessUnitStatic struct {
	clientCount int
	totalRecv int64
}
type EPollProcessUnit struct {
	chMsg chan interface{}
	eSrvMgr *EpollServerMgr
	mapTcpClient MAP_TCPCLIENT
	mapClientIDTcpFD MAP_CLIENTID_TCPFD

	ClientProcessUnitStatic
}


func (pthis*EPollProcessUnit)ProcessProtoMsg(msg *msgProto){
	switch msg.packType {
	case msgpacket.MSG_TYPE__MSG_LOGIN:
		{
			t, ok := msg.protoMsg.(*msgpacket.MSG_LOGIN)
			if t != nil && ok {
				pthis.Process_MSG_LOGIN(msg.fd, t)
			}
		}
	case msgpacket.MSG_TYPE__MSG_SRV_REPORT:
	case msgpacket.MSG_TYPE__MSG_RPC:
	case msgpacket.MSG_TYPE__MSG_RPC_RES:
	default:
		pthis.Process_protoMsg(msg)
	}
}

func (pthis*EPollProcessUnit)getClientByFD(fd int) *TcpClient {
	oldC, _ := pthis.mapTcpClient[fd]
	return oldC
}
func (pthis*EPollProcessUnit)getClientByClientID(clientID int64) *TcpClient {
	oldFD, ok := pthis.mapClientIDTcpFD[clientID]
	if !ok {
		return nil
	}
	oldC, _ := pthis.mapTcpClient[oldFD]
	return oldC
}
func (pthis*EPollProcessUnit)addClient(c *TcpClient) {
	pthis.mapTcpClient[c.fd.FD] = c
	pthis.mapClientIDTcpFD[c.clientID] = c.fd.FD

	pthis.clientCount = len(pthis.mapTcpClient)
}
func (pthis*EPollProcessUnit)delClient(fd lin_common.FD_DEF) {
	oldC, _ := pthis.mapTcpClient[fd.FD]
	if oldC != nil {
		delete(pthis.mapClientIDTcpFD, oldC.clientID)
	}
	delete(pthis.mapTcpClient, fd.FD)

	pthis.clientCount = len(pthis.mapTcpClient)
}

func (pthis*EPollProcessUnit)Process_msgTcpClose(msg *msgTcpClose) {
	c := pthis.getClientByFD(msg.fd.FD)
	if c == nil {
		return
	}
	if !c.fd.IsSame(&msg.fd) {
		return
	}
	lin_common.LogDebug(msg.fd.String(), " clientid:", c.clientID)
	pthis.delClient(msg.fd)
	c.Destructor()
}


func (pthis*EPollProcessUnit)Process_MSG_LOGIN(fd lin_common.FD_DEF, msg *msgpacket.MSG_LOGIN){
	lin_common.LogDebug("login:", fd.String(), " clientid:", msg.Id)

	oldC := pthis.getClientByClientID(msg.Id)
	if oldC != nil {
		if !oldC.fd.IsSame(&fd){
			if oldC.fd.FD != fd.FD {
				pthis.delClient(oldC.fd)
				pthis.eSrvMgr.lsn.EPollListenerCloseTcp(oldC.fd)
			}

			c := ConstructorTcpClient(pthis, fd, msg.Id)
			pthis.addClient(c)
		}
	} else {
		c := ConstructorTcpClient(pthis, fd, msg.Id)
		pthis.addClient(c)
	}

	msgRes := &msgpacket.MSG_LOGIN_RES{}
	msgRes.Id = msg.Id
	msgRes.ConnectId = int64(fd.Magic)
	msgRes.Fd = int64(fd.FD)

	pthis.eSrvMgr.SendProtoMsg(fd, msgpacket.MSG_TYPE__MSG_LOGIN_RES, msgRes)
}

func (pthis*EPollProcessUnit)Process_protoMsg(msg *msgProto) {
	pthis.totalRecv ++

	c := pthis.getClientByFD(msg.fd.FD)
	if c == nil {
		return
	}

	c.Process_protoMsg(msg)
}


func (pthis*EPollProcessUnit)_go_Process_unit(){
	for {
		msg := <- pthis.chMsg
		switch t := msg.(type) {
		case *msgProto:
			pthis.ProcessProtoMsg(t)
		case *msgTcpClose:
			pthis.Process_msgTcpClose(t)
		}
	}
}

func ConstructorEPollProcessUnit(eSrvMgr *EpollServerMgr) *EPollProcessUnit {
	return &EPollProcessUnit{
		chMsg : make(chan interface{}, 100),
		eSrvMgr : eSrvMgr,
		mapTcpClient : make(MAP_TCPCLIENT),
		mapClientIDTcpFD : make(MAP_CLIENTID_TCPFD),
	}
}
