package main

import (
	"lin/lin_common"
	"lin/server/server_common"
	"strconv"
	"sync"
)

type MAP_MSGQUESRV_STATUS map[server_common.SRV_ID]*MsgQueSrvStatus
type MsgQueSrvMgr struct {
	mapMsgQueSrv sync.Map // server_common.MSGQUE_SRV_ID - MsgQueSrvInfo

	mapQueSrvStatusLock sync.Mutex
	mapQueSrvStatus MAP_MSGQUESRV_STATUS
}

type MsgQueSrvInfo struct {
	fd lin_common.FD_DEF
	ip string
	port int32
	queSrvID server_common.SRV_ID
}

func (pthis*MsgQueSrvInfo)String()(str string){
	str = pthis.queSrvID.String() + " fd:" + pthis.fd.String() +
		"[" + pthis.ip + ":" + strconv.FormatInt(int64(pthis.port), 10) + "]"
	return
}

type MsgQueSrvStatus struct {
	queSrvID server_common.SRV_ID
	ChooseConnCount int // msg que connected srv count
}

func (pthis*MsgQueSrvMgr)StoreQueSrvInfo(qsi * MsgQueSrvInfo) {
	pthis.mapMsgQueSrv.Store(qsi.queSrvID, *qsi)

	{
		pthis.mapQueSrvStatusLock.Lock()
		defer pthis.mapQueSrvStatusLock.Unlock()
		v, ok := pthis.mapQueSrvStatus[qsi.queSrvID]
		if !ok {
			pthis.mapQueSrvStatus[qsi.queSrvID] = &MsgQueSrvStatus{queSrvID:qsi.queSrvID}
		} else {
			lin_common.LogInfo(qsi.queSrvID.String(), " status already exsit", v)
		}
	}
}

func (pthis*MsgQueSrvMgr)RangeQueSrvInfo(fn func(key, value any) bool) {
	pthis.mapMsgQueSrv.Range(fn)
}

func  (pthis*MsgQueSrvMgr)LoadQueSrvInfo(queSrvID server_common.SRV_ID) (qsi MsgQueSrvInfo, bRet bool) {
	bRet = false
	val, ok := pthis.mapMsgQueSrv.Load(queSrvID)
	if !ok {
		return
	}

	qsi, ok = val.(MsgQueSrvInfo)
	if !ok {
		return
	}

	bRet = true
	return
}

func (pthis*MsgQueSrvMgr)DeleteQueSrvInfo(queSrvID server_common.SRV_ID) {
	pthis.mapMsgQueSrv.Delete(queSrvID)
}

func (pthis*MsgQueSrvMgr)ChooseMostIdleQueSrv() (qsi MsgQueSrvInfo, bRet bool) {
	bRet = false
	pthis.mapQueSrvStatusLock.Lock()
	defer pthis.mapQueSrvStatusLock.Unlock()

	var status *MsgQueSrvStatus = nil
	minCount := 0
	for _, v := range pthis.mapQueSrvStatus {
		if status == nil {
			minCount = v.ChooseConnCount
			status = v
			continue
		}

		if minCount < v.ChooseConnCount {
			continue
		}

		minCount = v.ChooseConnCount
		status = v
	}

	if status == nil {
		return

	}
	status.ChooseConnCount ++

	lin_common.LogDebug("choose que:", status)

	qsi, bRet = pthis.LoadQueSrvInfo(status.queSrvID)
	return
}

func (pthis*MsgQueSrvMgr)ResetQueSrvChooseCount(queSrvID server_common.SRV_ID, chooseCount int) {
	pthis.mapQueSrvStatusLock.Lock()
	defer pthis.mapQueSrvStatusLock.Unlock()

	v, ok := pthis.mapQueSrvStatus[queSrvID]
	if ok && v != nil {
		v.ChooseConnCount = chooseCount
	}
}

func ConstructMsgQueSrvMgr()*MsgQueSrvMgr {
	mqMgr := &MsgQueSrvMgr{
		mapQueSrvStatus: make(MAP_MSGQUESRV_STATUS),
	}

	return mqMgr
}

func (pthis*MsgQueSrvMgr)Dump() (str string) {
	pthis.mapQueSrvStatusLock.Lock()
	defer pthis.mapQueSrvStatusLock.Unlock()

	str += "msg que srv reg:\r\n"
	pthis.mapMsgQueSrv.Range(func(key, value any) bool{
		qsi, ok := value.(MsgQueSrvInfo)
		if !ok {
			return true
		}

		str += qsi.String() + " "

		status, ok := pthis.mapQueSrvStatus[qsi.queSrvID]
		if ok {
			str += " choose count:" + strconv.FormatInt(int64(status.ChooseConnCount), 10)
		}
		str += "\r\n"
		return true
	})

	return
}