package main

import (
	"goserver/common"
	"goserver/msgpacket"
	"goserver/server/server_common"
	"net"
	"sync"
)


type MAP_SRV_NET_INFO   map[server_common.SRV_ID]*QueSrvNetInfo
type MAP_OTHER_SRV_INFO map[server_common.SRV_ID]*OtherSrvInfo

type CliSrvMgr struct {
	mapQueSrvRWLock sync.RWMutex
	mapQueSrvNet    MAP_SRV_NET_INFO
	mapOtherSrvInfo MAP_OTHER_SRV_INFO
}



// SrvNetInfo the server reg to this msg que
type QueSrvNetInfo struct{
	server_common.SrvBaseInfo

	fd   common.FD_DEF
	addr net.Addr
}

func (pthis*QueSrvNetInfo)String() string {
	str := pthis.SrvUUID.String() +
		pthis.SrvType.String() +
		pthis.fd.String() +
		" " + pthis.addr.String() +
		"\r\n"

	return str
}


// OtherSrvInfo server in the other msg que
type OtherSrvInfo struct{
	server_common.SrvBaseInfo

	queSrvID server_common.SRV_ID
}

func (pthis*OtherSrvInfo)String() string {
	str := pthis.SrvUUID.String() +
		pthis.SrvType.String() +
		" reg to msgque" + pthis.queSrvID.String() +
		"\r\n"

	return str
}


func (pthis*CliSrvMgr)addQueSrv(si *QueSrvNetInfo) {
	if si == nil {
		return
	}
	// write lock
	pthis.mapQueSrvRWLock.Lock()
	defer pthis.mapQueSrvRWLock.Unlock()

	pthis.mapQueSrvNet[si.SrvUUID] = si
}

func (pthis*CliSrvMgr)delQueSrv(srvUUID server_common.SRV_ID) *QueSrvNetInfo {
	// write lock
	pthis.mapQueSrvRWLock.Lock()
	defer pthis.mapQueSrvRWLock.Unlock()

	si, ok := pthis.mapQueSrvNet[srvUUID]
	if !ok || nil == si {
		return nil
	}
	delete(pthis.mapQueSrvNet, srvUUID)
	return si
}

// getAllSrvNetPB get all local accept srv
func (pthis*CliSrvMgr)getAllSrvNetPB(pb * msgpacket.PB_SRV_INFO_ALL) {
	if pb == nil {
		return
	}

	// read lock
	pthis.mapQueSrvRWLock.RLock()
	defer pthis.mapQueSrvRWLock.RUnlock()

	for _, v := range pthis.mapQueSrvNet {
		pb.ArraySrv = append(pb.ArraySrv,
			&msgpacket.PB_SRV_INFO_ONE{
				SrvUuid:int64(v.SrvUUID),
				SrvType :int32(v.SrvType),
			})
	}
}

func (pthis*CliSrvMgr)RangeAllSrvNet(fn func(server_common.SRV_ID, *QueSrvNetInfo)){
	pthis.mapQueSrvRWLock.RLock()
	defer pthis.mapQueSrvRWLock.RUnlock()

	for _, v := range pthis.mapQueSrvNet {
		fn(v.SrvUUID, v)
	}
}

func (pthis*CliSrvMgr)addOtherQueAllSrvFromPB(queSrvID server_common.SRV_ID, allSrv * msgpacket.PB_SRV_INFO_ALL) {
	// to other srv
	if allSrv == nil {
		return
	}

	pthis.delOtherQueAllSrv(queSrvID)

	// write lock
	pthis.mapQueSrvRWLock.Lock()
	defer pthis.mapQueSrvRWLock.Unlock()

	for _, v := range allSrv.ArraySrv {
		soi := OtherSrvInfo{
			SrvBaseInfo:server_common.SrvBaseInfo{
				SrvUUID:  server_common.SRV_ID(v.SrvUuid),
				SrvType:  server_common.SRV_TYPE(v.SrvType),
			},
			queSrvID: queSrvID,
		}
		pthis.mapOtherSrvInfo[soi.SrvUUID] = &soi
	}
}

func (pthis*CliSrvMgr)delOtherQueAllSrv(queSrvID server_common.SRV_ID) {
	// write lock
	pthis.mapQueSrvRWLock.Lock()
	defer pthis.mapQueSrvRWLock.Unlock()

	arrayID := make([]server_common.SRV_ID, 0)
	for _, v := range pthis.mapOtherSrvInfo {
		if v.queSrvID != queSrvID {
			continue
		}
		arrayID = append(arrayID, v.SrvUUID)
	}

	for _, v := range arrayID {
		delete(pthis.mapOtherSrvInfo, v)
	}
}

func (pthis*CliSrvMgr)getOtherQueAllSrv(pb * msgpacket.PB_SRV_INFO_ALL) {
	if nil == pb {
		return
	}
	for _, v := range pthis.mapOtherSrvInfo {
		pb.ArraySrv = append(pb.ArraySrv,
			&msgpacket.PB_SRV_INFO_ONE{
				SrvUuid:int64(v.SrvUUID),
				SrvType :int32(v.SrvType),
			})
	}
}

func (pthis*CliSrvMgr)findLocalRoute(srvUUID server_common.SRV_ID) (common.FD_DEF, bool) {
	pthis.mapQueSrvRWLock.RLock()
	defer pthis.mapQueSrvRWLock.RUnlock()

	v, ok := pthis.mapQueSrvNet[srvUUID]
	if !ok || nil == v {
		return common.FD_DEF_NIL, false
	}

	return v.fd, true
}

func (pthis*CliSrvMgr)findRemoteRoute(srvUUID server_common.SRV_ID)server_common.SRV_ID {
	pthis.mapQueSrvRWLock.RLock()
	defer pthis.mapQueSrvRWLock.RUnlock()

	v, ok := pthis.mapOtherSrvInfo[srvUUID]
	if !ok || nil == v {
		return server_common.SRV_ID_INVALID
	}

	return v.queSrvID
}

func (pthis*CliSrvMgr)findSrvByType(srvType server_common.SRV_TYPE) (arraySrvID []*server_common.SrvBaseInfo) {
	pthis.mapQueSrvRWLock.RLock()
	defer pthis.mapQueSrvRWLock.RUnlock()

	for _, v := range pthis.mapQueSrvNet {
		if srvType != v.SrvType && srvType != server_common.SRV_TYPE_none {
			continue
		}
		arraySrvID = append(arraySrvID, &server_common.SrvBaseInfo{SrvUUID:v.SrvUUID, SrvType:v.SrvType})
	}
	for _, v := range pthis.mapOtherSrvInfo {
		if srvType != v.SrvType && srvType != server_common.SRV_TYPE_none {
			continue
		}
		arraySrvID = append(arraySrvID, &server_common.SrvBaseInfo{SrvUUID:v.SrvUUID, SrvType:v.SrvType})
	}
	return
}


func ConstructorCliSrvMgr()*CliSrvMgr {
	smgr := &CliSrvMgr{
		mapQueSrvNet : make(MAP_SRV_NET_INFO),
		mapOtherSrvInfo : make(MAP_OTHER_SRV_INFO),
	}

	return smgr
}

func (pthis*CliSrvMgr)Dump() string {
	pthis.mapQueSrvRWLock.RLock()
	defer pthis.mapQueSrvRWLock.RUnlock()

	str := "\r\n srv in this msg que net info\r\n"
	{
		for _, v := range pthis.mapQueSrvNet {
			str += v.String()
		}
	}

	str += "\r\n srv in other msg que info\r\n"
	{
		for _, v:= range pthis.mapOtherSrvInfo {
			str += v.String()
		}
	}

	return str
}