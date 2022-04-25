package main

import (
	"lin/lin_common"
	"lin/msgpacket"
)

type MAP_CLIENT map[int]*TcpClient

type eSrvMgrProcessUnitStatic struct {
	clientCount int
	totalRecv int64
}
type eSrvMgrProcessUnit struct {
	chMsg chan interface{}
	eSrvMgr *EpollServerMgr
	mapClient MAP_CLIENT

	eSrvMgrProcessUnitStatic
}


func (pthis*eSrvMgrProcessUnit)ProcessProtoMsg(msg *msgProto){
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

func (pthis*eSrvMgrProcessUnit)getClient(fd lin_common.FD_DEF) *TcpClient {
	oldC, _ := pthis.mapClient[fd.FD]
	return oldC
}
func (pthis*eSrvMgrProcessUnit)addClient(c *TcpClient) {
	pthis.mapClient[c.fd.FD] = c

	pthis.clientCount = len(pthis.mapClient)
}
func (pthis*eSrvMgrProcessUnit)delClient(fd lin_common.FD_DEF) {
	delete(pthis.mapClient, fd.FD)

	pthis.clientCount = len(pthis.mapClient)
}

func (pthis*eSrvMgrProcessUnit)Process_msgTcpClose(msg *msgTcpClose) {
	c := pthis.getClient(msg.fd)
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


func (pthis*eSrvMgrProcessUnit)Process_MSG_LOGIN(fd lin_common.FD_DEF, msg *msgpacket.MSG_LOGIN){
	lin_common.LogDebug("login:", fd.String(), " clientid:", msg.Id)

	oldC := pthis.getClient(fd)
	if oldC != nil {
		if !oldC.fd.IsSame(&fd){
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

	pthis.eSrvMgr.lsn.EPollListenerWrite(fd, msgpacket.ProtoPacketToBin(msgpacket.MSG_TYPE__MSG_LOGIN_RES, msgRes))
	//lin_common.TMP_tcpWrite(fd, msgpacket.ProtoPacketToBin(msgpacket.MSG_TYPE__MSG_LOGIN_RES, msgRes))
}

func (pthis*eSrvMgrProcessUnit)Process_protoMsg(msg *msgProto) {
	pthis.totalRecv ++

	c := pthis.getClient(msg.fd)
	if c == nil {
		return
	}

	c.Process_protoMsg(msg)
}


func (pthis*eSrvMgrProcessUnit)_go_Process_unit(){
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
