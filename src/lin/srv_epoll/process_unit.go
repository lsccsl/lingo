package main

import (
	"github.com/golang/protobuf/proto"
	"lin/lin_common"
	"lin/msgpacket"
)

/* begin process unit msg define */
type CLIENT_MSG_TYPE int
const(
	CLIENT_TCP_CLOSE CLIENT_MSG_TYPE = 1
	CLIENT_LOGIN CLIENT_MSG_TYPE = 2
	CLIENT_PROTO CLIENT_MSG_TYPE = 3
)
type msgClient struct {
	clientID int64
	fd lin_common.FD_DEF
	msgType CLIENT_MSG_TYPE
	msg interface{}
}
/* end process unit msg define */


type MAP_CLIENT map[int64/* client id */]*TcpClient

type ClientProcessUnitStatic struct {
	clientCount int
	totalRecv int64
}
type EPollProcessUnit struct {
	_chMsg chan *msgClient
	eSrvMgr *EpollServerMgr
	mapClient MAP_CLIENT

	ClientProcessUnitStatic
}


func (pthis*EPollProcessUnit)getClient(clientID int64) *TcpClient {
	oldC, ok := pthis.mapClient[clientID]
	if !ok {
		return nil
	}
	return oldC
}
func (pthis*EPollProcessUnit)addClient(c *TcpClient) {
	pthis.mapClient[c.clientID] = c

	pthis.clientCount = len(pthis.mapClient)
}
func (pthis*EPollProcessUnit)delClient(cliID int64) {
	delete(pthis.mapClient, cliID)

	pthis.clientCount = len(pthis.mapClient)
}

func (pthis*EPollProcessUnit)Process_TcpClose(c *TcpClient, fd lin_common.FD_DEF) {
	if c == nil {
		return
	}
	if !c.fd.IsSame(&fd) {
		return
	}
	lin_common.LogDebug(fd.String(), " clientid:", c.clientID)
	pthis.delClient(c.clientID)
	c.Destructor()
}


func (pthis*EPollProcessUnit)Process_LOGIN(cliID int64, fd lin_common.FD_DEF){
	lin_common.LogDebug("login:", fd.String(), " clientid:", cliID)

	oldC := pthis.getClient(cliID)
	if oldC != nil {
		if !oldC.fd.IsSame(&fd){
			if oldC.fd.FD != fd.FD {
				pthis.delClient(oldC.clientID)
				pthis.eSrvMgr.lsn.EPollListenerCloseTcp(oldC.fd)
			}

			c := ConstructorTcpClient(pthis, fd, cliID)
			pthis.addClient(c)
		}
	} else {
		c := ConstructorTcpClient(pthis, fd, cliID)
		pthis.addClient(c)
	}

	msgRes := &msgpacket.MSG_LOGIN_RES{}
	msgRes.Id = cliID
	msgRes.ConnectId = int64(fd.Magic)
	msgRes.Fd = int64(fd.FD)

	pthis.eSrvMgr.SendProtoMsg(fd, msgpacket.MSG_TYPE__MSG_LOGIN_RES, msgRes)
}



func (pthis*EPollProcessUnit)_go_Process_unit(){
	for {
		msg := <- pthis._chMsg
		c := pthis.getClient(msg.clientID)
		if c == nil {
			if CLIENT_LOGIN == msg.msgType {
				pthis.Process_LOGIN(msg.clientID, msg.fd)
				continue
			}
		}
		switch msg.msgType {
		case CLIENT_PROTO:
			{
				pthis.totalRecv ++
				c.Process_protoMsg(msg)
			}
		case CLIENT_TCP_CLOSE:
			pthis.Process_TcpClose(c, msg.fd)
		}
	}
}

func (pthis*EPollProcessUnit)PushTcpLoginMsg(cliID int64, fd lin_common.FD_DEF){
	pthis._chMsg <- &msgClient{clientID: cliID, fd:fd, msgType: CLIENT_LOGIN}
}

func (pthis*EPollProcessUnit)PushTcpCloseMsg(cliID int64, fd lin_common.FD_DEF){
	pthis._chMsg <- &msgClient{clientID: cliID, fd:fd, msgType: CLIENT_TCP_CLOSE}
}

func (pthis*EPollProcessUnit)PushProtoMsg(cliID int64, fd lin_common.FD_DEF, msg proto.Message){
	pthis._chMsg <- &msgClient{clientID: cliID, fd:fd, msgType: CLIENT_PROTO, msg: msg}
}

func ConstructorEPollProcessUnit(eSrvMgr *EpollServerMgr) *EPollProcessUnit {
	return &EPollProcessUnit{
		_chMsg : make(chan *msgClient, 100),
		eSrvMgr : eSrvMgr,
		mapClient : make(MAP_CLIENT),
	}
}
