package main

import (
	"lin/lin_common"
	"time"
)

type MAP_NAVMAP_INS map[int64]*NavMapIns
type MapProcess struct {
	navIns *NavMapIns
	chMsg chan interface{}

	procMgr *MapProcessMgr
}

type MapProcessMgr struct {
	mapProc []*MapProcess

	navMap *NavMap
}

type msgPathSearch struct {
	src Coord3f
	dst Coord3f

	chRes chan interface{}
}
type msgPathSearchRes struct {
	path []Coord3f
}

func (pthis*MapProcess)_go_MapProcess() {
	ticker := time.NewTicker(time.Millisecond * 1000)
	for {
		select {
		case msg := <-pthis.chMsg:
			switch t := msg.(type) {
			case *msgPathSearch:
				pthis.process_msgPathSearch(t)
			}
		case <-ticker.C:
		}
	}
}

func (pthis*MapProcess)process_msgPathSearch(msg*msgPathSearch) {

	msgRes := &msgPathSearchRes{}
	msgRes.path = pthis.navIns.path_find(&msg.src, &msg.dst)

	msg.chRes <- msgRes
}

func ConstructMapProcess(procMgr *MapProcessMgr) *MapProcess {
	mp := &MapProcess {
		navIns : ConstructNavMapIns(),
		chMsg : make(chan interface{}),
		procMgr : procMgr,
	}

	mp.navIns.load_from_template(mp.procMgr.navMap)
	return mp
}

func (pthis*MapProcessMgr)pathSearch(src * Coord3f, dst * Coord3f, clientID int64) (path []Coord3f) {
	mp := pthis.getMapProcess(clientID)
	if mp == nil {
		return nil
	}

	chRes := make(chan interface{})
	mp.chMsg <- &msgPathSearch{*src, *dst, chRes}

	select {
	case msg := <- chRes:
		msgRes, ok := msg.(*msgPathSearchRes)
		if ok {
			path = msgRes.path
		}
	case <- time.After(time.Second * 3):
		path = nil
	}
	close(chRes)

	return
}

func (pthis *MapProcessMgr)getMapProcess(clientID int64)*MapProcess{
	countProc := len(pthis.mapProc)
	if countProc <= 0 {
		return nil
	}

	idx := clientID % int64(countProc)
	return pthis.mapProc[idx]
}

func ConstructMapProcessMgr(processCount int) *MapProcessMgr {
	if processCount <= 0 {
		lin_common.LogErr("map process count is 0")
		return nil
	}

	mgr := &MapProcessMgr{
		mapProc : make([]*MapProcess, 0, processCount),
		navMap : ConstructorNavMapMgr("../resource/test_scene.obj"),
	}

	for i := 0; i < processCount; i ++ {
		lin_common.LogDebug("load map process:", i)
		mp := ConstructMapProcess(mgr)
		lin_common.LogDebug("load map process:", i, " done")
		mgr.mapProc = append(mgr.mapProc, mp)
		go mp._go_MapProcess()
	}

	return mgr
}