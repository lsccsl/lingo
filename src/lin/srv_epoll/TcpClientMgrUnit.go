package main

import (
	"github.com/golang/protobuf/proto"
	"lin/lin_common"
	"lin/msgpacket"
	"time"
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
	attachData *TcpAttachData
}
/* end process unit msg define */


type MAP_CLIENT map[int64/* client id */]*TcpClient

type TcpClientMgrUnitStatic struct {
	clientCount int
	totalRecv int64
}
type TcpClientMgrUnit struct {
	_chMsg chan *msgClient
	eSrvMgr *ServerMgr
	mapClient MAP_CLIENT

	TcpClientMgrUnitStatic
}


func (pthis*TcpClientMgrUnit)getClient(clientID int64) *TcpClient {
	oldC, ok := pthis.mapClient[clientID]
	if !ok {
		return nil
	}
	return oldC
}
func (pthis*TcpClientMgrUnit)addClient(c *TcpClient) {
	pthis.mapClient[c.clientID] = c

	pthis.clientCount = len(pthis.mapClient)
}
func (pthis*TcpClientMgrUnit)delClient(cliID int64) {
	delete(pthis.mapClient, cliID)

	pthis.clientCount = len(pthis.mapClient)
}

func (pthis*TcpClientMgrUnit)Process_TcpClose(c *TcpClient, fd lin_common.FD_DEF) {
	if c == nil {
		return
	}
	if !c.fd.IsSame(&fd) {
		return
	}
	lin_common.LogDebug(fd.String(), " clientid:", c.clientID, " fd:", fd.String())
	pthis.delClient(c.clientID)

	{
		msgDel := &msgDelAOIObject {
			objID : c.objID,
		}
		pthis.eSrvMgr.mapProcMgr.addMapProcessMsg(msgDel, c.clientID, time.Second * 3)
		c.objID = msgDel.objID
	}

	c.Destructor()
}


func (pthis*TcpClientMgrUnit)Process_LOGIN(msg*msgClient){
	cliID := msg.clientID
	fd := msg.fd
	lin_common.LogDebug("login:", fd.String(), " clientid:", cliID)

	oldC := pthis.getClient(cliID)
	if oldC != nil {
		if !oldC.fd.IsSame(&fd){
			if oldC.fd.FD != fd.FD {
				pthis.delClient(oldC.clientID)
				pthis.eSrvMgr.lsn.EPollListenerCloseTcp(oldC.fd, EN_TCP_CLOSE_REASON_new_conn)
				oldC.Destructor()
			}

			c := ConstructorTcpClient(pthis, fd, cliID)
			pthis.addClient(c)
		}
	} else {
		c := ConstructorTcpClient(pthis, fd, cliID)
		pthis.addClient(c)
	}

	c := pthis.getClient(cliID)
	if c == nil {
		return
	}

	msgRes := &msgpacket.MSG_LOGIN_RES{}

	objID := 0
	{
		msgG := &msgGenAOIID{}
		pthis.eSrvMgr.mapProcMgr.addMapProcessMsg(msgG, cliID, time.Second * 3)
		objID = msgG.objID
		c.objID = objID
	}
	msgRes.ObjId = int64(c.objID)
	msgRes.Id = cliID
	msgRes.ConnectId = int64(fd.Magic)
	msgRes.Fd = int64(fd.FD)
	pthis.eSrvMgr.SendProtoMsg(fd, msgpacket.MSG_TYPE__MSG_LOGIN_RES, msgRes)

	{
		msgAdd := &msgAddAOIObject {
			objID: objID,
			ntf : ConstructMapAOINtfClient(cliID, fd, pthis.eSrvMgr),
			ViewRange : 10,
			fd : fd,
		}
		msgL, ok := msg.msg.(*msgpacket.MSG_LOGIN)
		if ok {
			msgAdd.X = msgL.X
			msgAdd.Y = msgL.Y
			msgAdd.ViewRange = msgL.ViewRange
		}
		pthis.eSrvMgr.mapProcMgr.addMapProcessMsg(msgAdd, cliID, time.Second * 3)
	}
}



func (pthis*TcpClientMgrUnit)_go_Process_unit(){
	for {
		msg := <- pthis._chMsg
		if CLIENT_LOGIN == msg.msgType {
			pthis.Process_LOGIN(msg)
			continue
		}

		c := pthis.getClient(msg.clientID)
		if c == nil {
			lin_common.LogErr("not process fd:", msg.fd.String(), " clientid:", msg.clientID, " msg:", msg.msg,
				" attach data:", msg.attachData)
			continue
		} else {
			switch msg.msgType {
			case CLIENT_PROTO:
				{
					pthis.totalRecv++
					c.Process_protoMsg(msg)
				}
			case CLIENT_TCP_CLOSE:
				pthis.Process_TcpClose(c, msg.fd)
			}
		}
	}
}

func (pthis*TcpClientMgrUnit)PushTcpLoginMsg(cliID int64, msgL *msgpacket.MSG_LOGIN, fd lin_common.FD_DEF){
	pthis._chMsg <- &msgClient{clientID: cliID, fd:fd, msg:msgL, msgType: CLIENT_LOGIN}
}

func (pthis*TcpClientMgrUnit)PushTcpCloseMsg(cliID int64, fd lin_common.FD_DEF){
	pthis._chMsg <- &msgClient{clientID: cliID, fd:fd, msgType: CLIENT_TCP_CLOSE}
}

func (pthis*TcpClientMgrUnit)PushProtoMsg(cliID int64, fd lin_common.FD_DEF, msg proto.Message, attachData *TcpAttachData){
	pthis._chMsg <- &msgClient{clientID: cliID,
		fd : fd,
		msgType : CLIENT_PROTO,
		msg : msg,
		attachData : attachData,
	}
}

func ConstructorTcpClientMgrUnit(eSrvMgr *ServerMgr) *TcpClientMgrUnit {
	return &TcpClientMgrUnit{
		_chMsg : make(chan *msgClient, 100),
		eSrvMgr : eSrvMgr,
		mapClient : make(MAP_CLIENT),
	}
}
