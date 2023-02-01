package msgque_client

import (
	"github.com/golang/protobuf/proto"
	"goserver/msgpacket"
	"goserver/server/server_common"
	"sync"
	"sync/atomic"
)

type MAP_MSG_REQ map[server_common.MSG_ID]*MsgReq
type ClientSrvMsgMgr struct {
	mapMsgReqMutex sync.Mutex
	mapMsgReq      MAP_MSG_REQ

	seq atomic.Int64
}
type MsgReq struct {
	MsgID server_common.MSG_ID
	chNtf chan *MsgRes
}
type MsgRes struct {
	PBMsg proto.Message
	Res msgpacket.PB_RESPONSE_CODE
}

func ConstructClientSrvMsgMgr() *ClientSrvMsgMgr {
	rmgr := &ClientSrvMsgMgr{
		mapMsgReq:make(MAP_MSG_REQ),
	}

	return rmgr
}

func (pthis*ClientSrvMsgMgr)ClientSrvMsgMgrAddReq(MsgID server_common.MSG_ID) *MsgReq {
	pthis.mapMsgReqMutex.Lock()
	defer pthis.mapMsgReqMutex.Unlock()

	val, _ := pthis.mapMsgReq[MsgID]
	if val != nil {
		return nil
	}
	rreq := &MsgReq{MsgID: MsgID,
		chNtf: make(chan *MsgRes, 10),
	}
	pthis.mapMsgReq[MsgID] = rreq
	return rreq
}

func (pthis*ClientSrvMsgMgr)ClientSrvMsgMgrFindReq(MsgID server_common.MSG_ID) *MsgReq {
	pthis.mapMsgReqMutex.Lock()
	defer pthis.mapMsgReqMutex.Unlock()
	val, _ := pthis.mapMsgReq[MsgID]
	return val
}

func (pthis*ClientSrvMsgMgr)ClientSrvMsgMgrDelReq(MsgID server_common.MSG_ID) {
	pthis.mapMsgReqMutex.Lock()
	defer pthis.mapMsgReqMutex.Unlock()
	val, _ := pthis.mapMsgReq[MsgID]
	if val != nil {
		if val.chNtf != nil {
			close(val.chNtf)
		}
	}
	delete(pthis.mapMsgReq, MsgID)
}

func (pthis*ClientSrvMsgMgr)Dump() string {
	str := "\r\n\r\n client srv msg mgr:\r\n"
	return str
}