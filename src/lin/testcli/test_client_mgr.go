package main

import "sync"

type MAP_CLIENT map[int64]*ClientTcpInfo

type ClientMgr struct{
	mapClientMutex sync.Mutex
	mapClient MAP_CLIENT
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
