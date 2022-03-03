package main

import (
	"sync"
	"time"
)

type dialData struct {
	dialTimeoutSec int
	closeExpireSec int
	tcpConn *TcpConnection
	srvID int64
	ip string
	port int

	needRedial bool
	redialCount int
}
type MAP_DIALDATA map[int64/* srvID */]*dialData

type TcpDialMgr struct {
	wg sync.WaitGroup
	closeExpireSec int
	connMgr InterfaceConnManage

	mapConnMutex sync.Mutex
	mapDialData MAP_DIALDATA
}


func (pthis * TcpDialMgr) TcpDialMgrStart(connMgr InterfaceConnManage, closeExpireSec int){
	pthis.closeExpireSec = closeExpireSec
	pthis.mapDialData = make(MAP_DIALDATA)
	pthis.connMgr = connMgr
	pthis.wg.Add(1)

	go pthis.go_checkRedial()
}

func (pthis * TcpDialMgr) go_checkRedial(){
	chTimer := time.After(time.Second * time.Duration(3))
	for {
		select {
		case <-chTimer:
			{
				chTimer = time.After(time.Second * time.Duration(3))
			}
		}
	}
}

func (pthis * TcpDialMgr)TcpDialMgrWait() {
	pthis.wg.Wait()
}


func (pthis * TcpDialMgr) addDialData(srvID int64, dd *dialData) {
	pthis.mapConnMutex.Lock()
	defer pthis.mapConnMutex.Unlock()

	pthis.mapDialData[srvID] = dd
}
func (pthis * TcpDialMgr) getDialData(srvID int64, ddOut *dialData) bool {
	pthis.mapConnMutex.Lock()
	defer pthis.mapConnMutex.Unlock()

	dd, ok := pthis.mapDialData[srvID]
	if !ok || dd == nil {
		return false
	}
	*ddOut = *dd
	ddOut.tcpConn = nil
	return true
}
func (pthis * TcpDialMgr) delDialData(srvID int64) {
	pthis.mapConnMutex.Lock()
	defer pthis.mapConnMutex.Unlock()

	delete(pthis.mapDialData, srvID)
}
func (pthis * TcpDialMgr) clearDialConn(srvID int64) {
	pthis.mapConnMutex.Lock()
	defer pthis.mapConnMutex.Unlock()
	dd, ok := pthis.mapDialData[srvID]
	if !ok || dd == nil {
		return
	}
	dd.tcpConn = nil
}

func (pthis * TcpDialMgr) TcpDialMgrDial(srvID int64, ip string, port int, closeExpireSec int,
	dialTimeoutSec int,
	needRedial bool, redialCount int) (*TcpConnection, error) {
	tcpConn, err := startTcpDial(pthis.connMgr, ip, port, closeExpireSec, dialTimeoutSec, redialCount)
	if err != nil {
		return nil, err
	}

	tcpConn.SrvID = srvID

	pthis.addDialData(srvID,
		&dialData{
			dialTimeoutSec:dialTimeoutSec,
			closeExpireSec:closeExpireSec,
			tcpConn:tcpConn,
			ip:ip,
			port:port,
			srvID:srvID,
			needRedial:needRedial,
			redialCount:redialCount})

	return tcpConn, nil
}

func (pthis * TcpDialMgr) TcpDialMgrCheckReDial(srvID int64) {
	var dd dialData
	bret := pthis.getDialData(srvID, &dd)
	if !bret {
		return
	}
	if !dd.needRedial {
		pthis.delDialData(srvID)
		return
	}

	pthis.TcpDialMgrDial(dd.srvID, dd.ip, dd.port, dd.closeExpireSec, dd.dialTimeoutSec, dd.needRedial, dd.redialCount)
}

func (pthis * TcpDialMgr) TcpDialDelDialData(srvID int64) {
	pthis.delDialData(srvID)
}
