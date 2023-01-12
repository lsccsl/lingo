package main

import "sync"

type MAP_TEST_SERVER map[int64]*TestSrv

type TestSrvMgr struct{
	mapSrvMutex sync.Mutex
	mapSrv      MAP_TEST_SERVER

	timestamp float64
	total int64
	totalLast int64

	totalReqRecv int64
	totalReqRecvLast int64
}

func (pthis *TestSrvMgr)TestSrvMgrAdd(s *TestSrv) {
	pthis.mapSrvMutex.Lock()
	defer pthis.mapSrvMutex.Unlock()

	pthis.mapSrv[s.srvId] = s
}

func (pthis *TestSrvMgr)TestSrvMgrDel(id int64) {
	pthis.mapSrvMutex.Lock()
	defer pthis.mapSrvMutex.Unlock()

	delete(pthis.mapSrv, id)
}

func (pthis *TestSrvMgr)TestSrvMgrGet(id int64) *TestSrv {
	pthis.mapSrvMutex.Lock()
	defer pthis.mapSrvMutex.Unlock()

	c, _ := pthis.mapSrv[id]
	return c
}