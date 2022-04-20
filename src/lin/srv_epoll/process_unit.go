package main

import (
	"lin/lin_common"
	"lin/msgpacket"
)

type MAP_CLIENT map[int]*TcpClient

type eSrvMgrProcessUnit struct {
	chMsg chan interface{}
	eSrvMgr *EpollServerMgr
	mapClient MAP_CLIENT
}


func (pthis*eSrvMgrProcessUnit)ProcessProtoMsg(msg *msgProto){
	switch t := msg.protoMsg.(type) {
	case *msgpacket.MSG_LOGIN:
		pthis.Process_MSG_LOGIN(msg.fd, t)
	case *msgpacket.MSG_SRV_REPORT:
	case *msgpacket.MSG_RPC:
	case *msgpacket.MSG_RPC_RES:
	default:
	}
}

func (pthis*eSrvMgrProcessUnit)getClient(fd lin_common.FD_DEF) *TcpClient {
	oldC, _ := pthis.mapClient[fd.FD]
	return oldC
}
func (pthis*eSrvMgrProcessUnit)addClient(c *TcpClient) {
	pthis.mapClient[c.fd.FD] = c
}


func (pthis*eSrvMgrProcessUnit)Process_MSG_LOGIN(fd lin_common.FD_DEF, msg *msgpacket.MSG_LOGIN){
	lin_common.LogDebug("login:", fd.String(), " clientid:", msg.Id)

	oldC := pthis.getClient(fd)
	if oldC != nil {
		if !oldC.fd.IsSame(&fd){
			c := ConstructorTcpClient(fd)
			pthis.addClient(c)
		}
	} else {
		c := ConstructorTcpClient(fd)
		pthis.addClient(c)
	}

	msgRes := &msgpacket.MSG_LOGIN_RES{}
	msgRes.Id = msg.Id
	msgRes.ConnectId = int64(fd.Magic)
	msgRes.Fd = int64(fd.FD)

	pthis.eSrvMgr.lsn.EPollListenerWrite(fd, msgpacket.ProtoPacketToBin(msgpacket.MSG_TYPE__MSG_LOGIN_RES, msgRes))
}

func (pthis*eSrvMgrProcessUnit)_go_Process_unit(){
	for {
		msg := <- pthis.chMsg
		switch t := msg.(type) {
		case *msgProto:
			pthis.ProcessProtoMsg(t)
		}
	}
}
