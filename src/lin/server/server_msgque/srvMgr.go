package main

import (
	"lin/lin_common"
	"lin/msgpacket"
	"lin/server/server_common"
	"net"
	"sync"
)


type MAP_SRV_NET_INFO   map[server_common.SRV_ID]*SrvNetInfo
type MAP_OTHER_SRV_INFO map[server_common.SRV_ID]*OtherSrvInfo

type SrvMgr struct {
	mapSrvRWLock sync.RWMutex
	mapSrvNet MAP_SRV_NET_INFO
	mapOtherSrvInfo MAP_OTHER_SRV_INFO
}


// SrvNetInfo the server reg to this msg que
type SrvNetInfo struct{
	srvUUID server_common.SRV_ID
	srvType server_common.SRV_TYPE
	fd lin_common.FD_DEF
	addr net.Addr
}

func (pthis*SrvNetInfo)String() string {
	str := pthis.srvUUID.String() +
		pthis.srvType.String() +
		pthis.fd.String() +
		" " + pthis.addr.String() +
		"\r\n"

	return str
}


// OtherSrvInfo server in the other msg que
type OtherSrvInfo struct{
	srvUUID server_common.SRV_ID
	srvType server_common.SRV_TYPE

	queSrvID server_common.SRV_ID
}

func (pthis*OtherSrvInfo)String() string {
	str := pthis.srvUUID.String() +
		pthis.srvType.String() +
		" reg to msgque" + pthis.queSrvID.String()

	return str
}


func (pthis*SrvMgr)addSrv(si *SrvNetInfo) {
	// write lock
	pthis.mapSrvRWLock.Lock()
	defer pthis.mapSrvRWLock.Unlock()

	pthis.mapSrvNet[si.srvUUID] = si
}

func (pthis*SrvMgr)delSrv(srvUUID server_common.SRV_ID) {
	// write lock
	pthis.mapSrvRWLock.Lock()
	defer pthis.mapSrvRWLock.Unlock()

	delete(pthis.mapSrvNet, srvUUID)
}

// getAllSrvNetPB get all local accept srv
func (pthis*SrvMgr)getAllSrvNetPB(pb * msgpacket.PB_SRV_INFO_ALL) {
	if pb == nil {
		return
	}

	// read lock
	pthis.mapSrvRWLock.RLock()
	defer pthis.mapSrvRWLock.RUnlock()

	for _, v := range pthis.mapSrvNet {
		pb.ArraySrv = append(pb.ArraySrv,
			&msgpacket.PB_SRV_INFO_ONE{
				SrvUuid:int64(v.srvUUID),
				SrvType :int32(v.srvType),
			})
	}
}

func (pthis*SrvMgr)addOtherQueAllSrvFromPB(queSrvID server_common.SRV_ID, allSrv * msgpacket.PB_SRV_INFO_ALL) {
	// to other srv
	if allSrv == nil {
		return
	}

	pthis.delOtherQueAllSrv(queSrvID)

	// write lock
	pthis.mapSrvRWLock.Lock()
	defer pthis.mapSrvRWLock.Unlock()

	for _, v := range allSrv.ArraySrv {
		soi := OtherSrvInfo{
			srvUUID:  server_common.SRV_ID(v.SrvUuid),
			srvType:  server_common.SRV_TYPE(v.SrvType),
			queSrvID: queSrvID,
		}
		pthis.mapOtherSrvInfo[soi.srvUUID] = &soi
	}
}

func (pthis*SrvMgr)delOtherQueAllSrv(queSrvID server_common.SRV_ID) {
	// write lock
	pthis.mapSrvRWLock.Lock()
	defer pthis.mapSrvRWLock.Unlock()

	arrayID := make([]server_common.SRV_ID, 0)
	for _, v := range pthis.mapOtherSrvInfo {
		if v.queSrvID != queSrvID {
			continue
		}
		arrayID = append(arrayID, v.srvUUID)
	}

	for _, v := range arrayID {
		delete(pthis.mapOtherSrvInfo, v)
	}
}


func ConstructorSrvMgr()*SrvMgr {
	smgr := &SrvMgr{
		mapSrvNet : make(MAP_SRV_NET_INFO),
		mapOtherSrvInfo : make(MAP_OTHER_SRV_INFO),
	}

	return smgr
}

func (pthis*SrvMgr)Dump() string {
	pthis.mapSrvRWLock.RLock()
	defer pthis.mapSrvRWLock.RUnlock()

	str := "\r\n srv in this msg que net info\r\n"
	{
		for _, v := range pthis.mapSrvNet {
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