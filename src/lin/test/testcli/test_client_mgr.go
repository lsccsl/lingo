package main

import "sync"

type MAP_CLIENT map[int64]*ClientTcpInfo

type ClientMgr struct{
	mapClientMutex sync.Mutex
	mapClient      MAP_CLIENT

	timestamp float64
	total int64
	totalLast int64
}

func (pthis *ClientMgr)ClientMgrAdd(c *ClientTcpInfo) {
	pthis.mapClientMutex.Lock()
	defer pthis.mapClientMutex.Unlock()

	pthis.mapClient[c.id] = c
}

func (pthis *ClientMgr)ClientMgrDel(id int64) {
	pthis.mapClientMutex.Lock()
	defer pthis.mapClientMutex.Unlock()

	delete(pthis.mapClient, id)
}

func (pthis *ClientMgr)ClientMgrGet(id int64) *ClientTcpInfo {
	pthis.mapClientMutex.Lock()
	defer pthis.mapClientMutex.Unlock()

	c, _ := pthis.mapClient[id]
	return c
}
