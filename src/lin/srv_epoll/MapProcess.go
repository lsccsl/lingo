package main

import (
	"lin/lin_common"
	"time"
)

type MAP_NAVMAP_INS map[int64]*NavMapIns
type MapProcess struct {
	// nav map search
	navIns *NavMapIns
	// map aoi
	aoi *MapAOI

	chMsg chan *msgMapProcess

	procMgr *MapProcessMgr
}

type MapProcessMgr struct {
	mapProc []*MapProcess

	navMap *NavMap
}


type NavObstacle struct {
	obstacleID uint32
	center Coord3f
	halfExt Coord3f
	yRadian float32
}
// begin map process msg
type msgMapProcess struct {
	msg interface{}
	chRes chan *msgMapProcess
}
type msgNavPathSearch struct {
	src Coord3f
	dst Coord3f
	path []Coord3f
}
type msgNavAddObstacle struct {
	ob NavObstacle
}
type msgNavDelObstacle struct {
	obstacleID uint32
}
type msgNavGetAllObstacle struct {
	ob []*NavObstacle
}
type msgAddAOIObject struct {
	aoiID int

	ntf MapAoiInf
	X float32
	Y float32
	ViewRange float32
}
type msgDelAOIObject struct {
	aoiID int
}
// end map process msg


func (pthis*MapProcess)_go_MapProcess() {
	ticker := time.NewTicker(time.Millisecond * 1000)
	for {
		pthis._go_MapProcess_loop(ticker)
	}
}

func (pthis*MapProcess)_go_MapProcess_loop(ticker *time.Ticker) {
	defer func() {
		err := recover()
		if err != nil {
			lin_common.LogErr(err)
		}
	}()

	select {
	case msg := <-pthis.chMsg:
		switch t := msg.msg.(type) {
		case *msgNavPathSearch:
			pthis.process_msgPathSearch(t)
		case *msgNavAddObstacle:
			pthis.process_msgNavAddObstacle(t)
		case *msgNavDelObstacle:
			pthis.process_msgNavDelObstacle(t)
		case *msgNavGetAllObstacle:
			pthis.process_msgNavGetAllObstacle(t)
		case *msgAddAOIObject:
			pthis.process_msgAddAOIObject(t)
		case *msgDelAOIObject:
			pthis.process_msgDelAOIObject(t)
		}
		msg.chRes <- msg
	case <-ticker.C:
	}
}

func (pthis*MapProcess)process_msgPathSearch(msg *msgNavPathSearch) {
	msg.path = pthis.navIns.path_find(&msg.src, &msg.dst)
}
func (pthis*MapProcess)process_msgNavAddObstacle(msg *msgNavAddObstacle) {
	msg.ob.obstacleID = pthis.navIns.add_obstacle(&Coord3f{msg.ob.center.X,msg.ob.center.Y, msg.ob.center.Z},
		&Coord3f{msg.ob.halfExt.X,msg.ob.halfExt.Y, msg.ob.halfExt.Z},
		msg.ob.yRadian)
}
func (pthis*MapProcess)process_msgNavDelObstacle(msg *msgNavDelObstacle) {
	pthis.navIns.del_obstacle(msg.obstacleID)
}
func (pthis*MapProcess)process_msgNavGetAllObstacle(msg *msgNavGetAllObstacle) {
	map_obstacle := pthis.navIns.get_all_obstacle()

	for k,v := range map_obstacle {
		ob := &NavObstacle{}
		ob.obstacleID = k
		ob.center = v.center
		ob.halfExt = v.half_ext
		ob.yRadian = v.y_radian
		msg.ob = append(msg.ob, ob)
	}
}
func (pthis*MapProcess)process_msgAddAOIObject(msg *msgAddAOIObject){
	msg.aoiID = pthis.aoi.add(msg.X, msg.Y, msg.ViewRange, msg.ntf)
	lin_common.LogDebug("add aoi ", msg.aoiID)
}
func (pthis*MapProcess)process_msgDelAOIObject(msg *msgDelAOIObject){
	lin_common.LogDebug("del aoi ", msg.aoiID)
	pthis.aoi.del(msg.aoiID)
}

func ConstructMapProcess(procMgr *MapProcessMgr) *MapProcess {
	mp := &MapProcess {
		navIns : ConstructNavMapIns(),
		chMsg : make(chan *msgMapProcess),
		procMgr : procMgr,
		aoi : ConstructorMapAOI(),
	}

	mp.navIns.load_from_template(mp.procMgr.navMap)
	return mp
}

func (pthis*MapProcessMgr)addMapProcessMsg(msg interface{}, clientID int64, timeOut time.Duration) {
	mp := pthis.getMapProcess(clientID)
	if mp == nil {
		return
	}

	chRes := make(chan *msgMapProcess)
	mp.chMsg <- &msgMapProcess{msg, chRes}

	select {
	case <- chRes:
	case <- time.After(timeOut):
		lin_common.LogErr("time out clientID:", clientID, " msg:", msg)
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