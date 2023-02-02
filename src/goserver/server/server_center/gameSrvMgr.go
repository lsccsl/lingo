package main

import (
	"goserver/server/server_common"
	"strconv"
	"sync"
)


type MAP_GAMESRV_INFO map[server_common.SRV_ID]*GameSrvInfo
type GameSrvMgr struct {
	mapGameSrvLock sync.RWMutex
	mapGameSrv     MAP_GAMESRV_INFO
}

type GameSrvInfo struct {
	server_common.SrvBaseInfo
	outIP string
	outPort int
}

func (pthis*GameSrvMgr)SetGameSrv(arraySrv []*server_common.SrvBaseInfo) {
	pthis.mapGameSrvLock.Lock()
	defer pthis.mapGameSrvLock.Unlock()

	pthis.mapGameSrv = make(MAP_GAMESRV_INFO)
	for _, v := range arraySrv {
		pthis.mapGameSrv[v.SrvUUID] = &GameSrvInfo{
			SrvBaseInfo:server_common.SrvBaseInfo{
				SrvUUID: v.SrvUUID,
				SrvType: v.SrvType,
			},
		}
	}
}

func (pthis*GameSrvMgr)SetGameSrvOutAddr(srvUUID server_common.SRV_ID, ip string, port int) {
	pthis.mapGameSrvLock.Lock()
	defer pthis.mapGameSrvLock.Unlock()

	gs, ok := pthis.mapGameSrv[srvUUID]
	if !ok || nil == gs {
		return
	}
	gs.outIP = ip
	gs.outPort = port
}

func (pthis*GameSrvMgr)GetGameSrv() (gi GameSrvInfo) {
	pthis.mapGameSrvLock.RLock()
	defer pthis.mapGameSrvLock.RUnlock()

	if 0 == len(pthis.mapGameSrv) {
		gi = GameSrvInfo{}
		return
	}

	for _, k := range pthis.mapGameSrv {
		gi = *k
		break
	}
	return
}

func ConstructGameSrvMgr()*GameSrvMgr {
	gmgr := &GameSrvMgr{
		mapGameSrv: make(MAP_GAMESRV_INFO),
	}

	return gmgr
}

func (pthis*GameSrvMgr)Dump() string {
	str := "\r\n\r\n game srv:\r\n"

	pthis.mapGameSrvLock.RLock()
	defer pthis.mapGameSrvLock.RUnlock()

	for _, v := range pthis.mapGameSrv {
		str += v.SrvUUID.String() + v.SrvType.String() + "[" + v.outIP + ":" + strconv.FormatInt(int64(v.outPort), 10) + "]" + "\r\n"
	}

	return str
}