package main

import (
	"lin/lin_common"
	"lin/msgpacket"
	"lin/server/server_common"
	"net"
	"sync"
)

type MAP_SRV_NET_INFO   map[server_common.MSGQUE_SRV_ID]*SrvNetInfo
type MAP_OTHER_SRV_INFO map[server_common.MSGQUE_SRV_ID]*OtherSrvInfo

type SrvMgr struct {
	mapSrvRWLock sync.RWMutex
	mapSrvNet MAP_SRV_NET_INFO
	mapOtherSrvInfo MAP_OTHER_SRV_INFO
}

type SrvNetInfo struct{
	srvUUID server_common.MSGQUE_SRV_ID
	srvType int32
	fd lin_common.FD_DEF
	addr net.Addr
}

type OtherSrvInfo struct{
	srvUUID server_common.MSGQUE_SRV_ID
	srvType int32

	queSrvID server_common.MSGQUE_SRV_ID
}

func (pthis*SrvMgr)addSrv(si *SrvNetInfo) {
	pthis.mapSrvRWLock.Lock()
	defer pthis.mapSrvRWLock.Unlock()

	pthis.mapSrvNet[si.srvUUID] = si
}

func (pthis*SrvMgr)addOtherSrv(soi *OtherSrvInfo) {
	pthis.mapSrvRWLock.Lock()
	defer pthis.mapSrvRWLock.Unlock()

	pthis.mapOtherSrvInfo[soi.srvUUID] = soi
}

func (pthis*SrvMgr)getAllSrvPB(pb * msgpacket.PB_MSG_INTER_SRV_REPORT_TO_OTHER_QUE) {
	pthis.mapSrvRWLock.RLock()
	defer pthis.mapSrvRWLock.RUnlock()

	for _, v := range pthis.mapSrvNet {
		pb.Srv = append(pb.Srv,
			&msgpacket.PbSrvDef{
				SrvUuid:int64(v.srvUUID),
				SrvType :v.srvType,
			})
	}
}

func ConstructorSrvMgr()*SrvMgr {
	smgr := &SrvMgr{
		mapSrvNet : make(MAP_SRV_NET_INFO),
		mapOtherSrvInfo : make(MAP_OTHER_SRV_INFO),
	}

	return smgr
}

