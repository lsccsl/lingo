package main

import (
	"lin/server/server_common"
	"sync"
)


type MAP_GAMESRV_INFO map[server_common.SRV_ID]*GameSrvInfo
type GameSrvMgr struct {
	mapGameSrvLock sync.RWMutex
	mapGameSrv MAP_GAMESRV_INFO
}

type GameSrvInfo struct {
	server_common.SrvBaseInfo
}

func ConstructGameSrvMgr()*GameSrvMgr {
	gmgr := &GameSrvMgr{
		mapGameSrv: make(MAP_GAMESRV_INFO),
	}

	return gmgr
}

