package main

import (
	"goserver/common"
	"goserver/server/server_common"
	"sync"
)

type otherMsgQueSrvMgr struct {
	mapOtherMsgQueSrv sync.Map // server_common.MSGQUE_SRV_ID - otherMsgQueSrvInfo
}

func (pthis*otherMsgQueSrvMgr)Clear() {
	pthis.mapOtherMsgQueSrv = sync.Map{}
}

func (pthis*otherMsgQueSrvMgr)updateQueSrvAccept(queSrvID server_common.SRV_ID, fdAccept common.FD_DEF) {
	qsi := otherMsgQueSrvInfo{
		queSrvID: queSrvID,
	}
	val, ok := pthis.mapOtherMsgQueSrv.Load(queSrvID)
	if ok {
		common.LogInfo("find que srv", queSrvID.String())
		qsi, ok = val.(otherMsgQueSrvInfo)
		if !ok {
			common.LogInfo("map value convert err, que srv", queSrvID.String())
			return
		}
	} else {
		common.LogInfo("can't find que srv", queSrvID.String())
	}
	qsi.fdAccept = fdAccept
	pthis.mapOtherMsgQueSrv.Store(qsi.queSrvID, qsi)
}

/*func (pthis*otherMsgQueSrvMgr)updateQueSrvDial(queSrvID server_common.SRV_ID, fdDial lin_common.FD_DEF) {

	val, ok := pthis.mapOtherMsgQueSrv.Load(queSrvID)
	if !ok {
		lin_common.LogInfo("can't find que srv", queSrvID.String())
		return
	}

	qsi := otherMsgQueSrvInfo{
		queSrvID: queSrvID,
	}
	qsi, ok = val.(otherMsgQueSrvInfo)
	if !ok {
		lin_common.LogInfo("map value convert err, que srv", queSrvID.String())
		return
	}

	qsi.fdDial = fdDial
	pthis.mapOtherMsgQueSrv.Store(qsi.queSrvID, qsi)
}*/

func (pthis*otherMsgQueSrvMgr)updateQueSrv(queSrvID server_common.SRV_ID,
	fdDial common.FD_DEF,
	ip string,
	port int32){

	qsi := otherMsgQueSrvInfo{
		fdDial:fdDial,
		ip:ip,
		port:port,
		queSrvID:queSrvID,
	}
	val, ok1 := pthis.mapOtherMsgQueSrv.Load(qsi.queSrvID)
	if ok1 {
		otherQSI, ok2 := val.(otherMsgQueSrvInfo)
		if ok2 {
			qsi.fdAccept = otherQSI.fdAccept
		}
	}
	pthis.mapOtherMsgQueSrv.Store(qsi.queSrvID, qsi)
}

func (pthis*otherMsgQueSrvMgr)LoadAndDelete(queSrvID server_common.SRV_ID, qsi *otherMsgQueSrvInfo) bool {
	if qsi == nil {
		return false
	}
	val, ok := pthis.mapOtherMsgQueSrv.LoadAndDelete(queSrvID)
	if !ok {
		return false
	}
	*qsi, ok = val.(otherMsgQueSrvInfo)
	if !ok {
		return false
	}
	return true
}

func (pthis*otherMsgQueSrvMgr)Load(queSrvID server_common.SRV_ID, qsi *otherMsgQueSrvInfo) bool {
	if qsi == nil {
		return false
	}
	val, ok := pthis.mapOtherMsgQueSrv.Load(queSrvID)
	if !ok {
		return false
	}
	*qsi, ok = val.(otherMsgQueSrvInfo)
	if !ok {
		return false
	}
	return true
}

func (pthis*otherMsgQueSrvMgr)Store(qsi *otherMsgQueSrvInfo) {
	if qsi == nil {
		return
	}
	pthis.mapOtherMsgQueSrv.Store(qsi.queSrvID, *qsi)
}

func (pthis*otherMsgQueSrvMgr)Range(fn func(key, value any) bool) {
	pthis.mapOtherMsgQueSrv.Range(fn)
}

func (pthis*otherMsgQueSrvMgr)Dump() string {
	str := "connect to other msg que srv\r\n"
	pthis.Range(func(key, value any) bool{
		qsi, ok := value.(otherMsgQueSrvInfo)
		if !ok {
			return true
		}

		str += qsi.String() + "\r\n"
		return true
	})

	return str
}