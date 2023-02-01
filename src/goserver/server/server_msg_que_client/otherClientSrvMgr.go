package msgque_client

import (
	"goserver/msgpacket"
	"goserver/server/server_common"
	"sync"
)

type MAP_OTHER_CLIENTSRV_INFO map[server_common.SRV_ID]*OtherClientSrvInfo
type OtherClientSrvMgr struct {
	mapOtherCliSrvRWLock sync.RWMutex
	mapOtherCliSrv       MAP_OTHER_CLIENTSRV_INFO
}

type OtherClientSrvInfo struct {
	srvUUID server_common.SRV_ID
	srvType server_common.SRV_TYPE
}

func (pthis*OtherClientSrvMgr)SetOtherClientSrvFromPB(allSrv * msgpacket.PB_SRV_INFO_ALL) {
	if nil == allSrv {
		return
	}
	if nil == allSrv.ArraySrv {
		return
	}

	pthis.mapOtherCliSrvRWLock.Lock()
	defer pthis.mapOtherCliSrvRWLock.Unlock()

	pthis.mapOtherCliSrv = make(MAP_OTHER_CLIENTSRV_INFO)

	for _, v := range allSrv.ArraySrv {
		pthis.mapOtherCliSrv[server_common.SRV_ID(v.SrvUuid)] = &OtherClientSrvInfo{
			srvUUID: server_common.SRV_ID(v.SrvUuid),
			srvType: server_common.SRV_TYPE(v.SrvType),
		}
	}
}

func (pthis*OtherClientSrvMgr)GetCliSrvByType(srvType server_common.SRV_TYPE) (arraySrv []*server_common.SrvBaseInfo) {
	pthis.mapOtherCliSrvRWLock.RLock()
	defer pthis.mapOtherCliSrvRWLock.RUnlock()

	for _, v := range pthis.mapOtherCliSrv {
		if srvType != v.srvType && srvType != server_common.SRV_TYPE_none {
			continue
		}
		arraySrv = append(arraySrv, &server_common.SrvBaseInfo{SrvUUID: v.srvUUID, SrvType: v.srvType})
	}

	return
}

func ConstructOtherClientSrvMgr()*OtherClientSrvMgr {
	omgr := &OtherClientSrvMgr{
		mapOtherCliSrv:make(MAP_OTHER_CLIENTSRV_INFO),
	}

	return omgr
}

func (pthis*OtherClientSrvMgr)Dump() string {
	pthis.mapOtherCliSrvRWLock.RLock()
	defer pthis.mapOtherCliSrvRWLock.RUnlock()

	str := "\r\n\r\nclient srv:\r\n"
	for _, v := range pthis.mapOtherCliSrv {
		str += v.srvUUID.String() + v.srvType.String() + "\r\n"
	}
	return str
}