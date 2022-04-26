package main

import "lin/lin_common"

type srvEvt_addremote struct {
	fd lin_common.FD_DEF
	srvID int64
	addr string
	closeExpireSec int
}

type MAP_TCPSRV map[int64/*srv id*/]*TcpSrv

type TcpSrvMgrProcessUnit struct {
	chSrv chan interface{}
	tcpSrvMgr *TcpSrvMgr
}
type TcpSrvMgr struct {
	eSrvMgr *EpollServerMgr
	mapSrv MAP_TCPSRV

	processUnit []*TcpSrvMgrProcessUnit
}

func (pthis*TcpSrvMgrProcessUnit)_go_Process_unit(){
	for {
		msg := <- pthis.chSrv
		switch t := msg.(type) {
		case *srvEvt_addremote:
			pthis.process_srvEvt_addremote(t)
		}
	}
}

func (pthis*TcpSrvMgrProcessUnit)process_srvEvt_addremote(evt * srvEvt_addremote){
	// todo : add srv
}

func (pthis*TcpSrvMgr)getProcessUnit(fd lin_common.FD_DEF)*TcpSrvMgrProcessUnit{
	processUnitCount := len(pthis.processUnit)
	idx := fd.FD % processUnitCount
	if idx >= processUnitCount {
		return nil
	}
	pu := pthis.processUnit[idx]
	if pu == nil {
		return nil
	}
	return pu
}

func (pthis*TcpSrvMgr)TcpSrvMgrAddRemoteSrv(srvID int64, addr string, closeExpireSec int){
	fd, err := pthis.eSrvMgr.lsn.EPollListenerDial(addr)
	if err != nil {
		lin_common.LogErr("srv:", srvID, " dial err")
	}
	lin_common.LogDebug("srv:", srvID, " fd:", fd.String())

	pu := pthis.getProcessUnit(fd)
	if pu != nil {
		pu.chSrv <- &srvEvt_addremote{
			fd : fd,
			srvID : srvID,
			addr : addr,
			closeExpireSec : closeExpireSec,
		}
	}
}

func ConstructorTcpSrvMgr(eSrvMgr *EpollServerMgr, srvProcessUnitCount int) *TcpSrvMgr {
	tcpSrvMgr := &TcpSrvMgr{
		eSrvMgr : eSrvMgr,
		mapSrv : make(MAP_TCPSRV),
	}

/*	for i := 0; i < srvProcessUnitCount; i ++ {
		processUnit := &TcpSrvMgrProcessUnit{
			chSrv : make(chan interface{}, 100),
			tcpSrvMgr : tcpSrvMgr,
		}
		processUnit._go_Process_unit()
	}*/

	return tcpSrvMgr
}

