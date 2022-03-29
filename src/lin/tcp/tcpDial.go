package tcp

import (
	"context"
	"sync"
)

type dialData struct {
	dialTimeoutSec int
	closeExpireSec int
	tcpConnID TCP_CONNECTION_ID
	srvID int64
	ip string
	port int

	DialCancelFunc context.CancelFunc

	needRedial bool
	redialCount int
}
type MAP_DIALDATA map[int64/* srvID */]*dialData

type TcpDialMgr struct {
	wg sync.WaitGroup
	closeExpireSec int
	connMgr        InterfaceConnManage

	mapDialDataMutex sync.Mutex
	mapDialData      MAP_DIALDATA
}


func (pthis *TcpDialMgr) TcpDialMgrStart(connMgr InterfaceConnManage, closeExpireSec int){
	pthis.closeExpireSec = closeExpireSec
	pthis.mapDialData = make(MAP_DIALDATA)
	pthis.connMgr = connMgr
	pthis.wg.Add(1)

	//go pthis.go_checkRedial()
}

func (pthis *TcpDialMgr)TcpDialMgrWait() {
	pthis.wg.Wait()
}

func (pthis *TcpDialMgr) addDialData(srvID int64, dd *dialData) *dialData {
	pthis.mapDialDataMutex.Lock()
	defer pthis.mapDialDataMutex.Unlock()

	pthis.mapDialData[srvID] = dd
	return dd
}
func (pthis *TcpDialMgr) getDialData(srvID int64) *dialData {
	pthis.mapDialDataMutex.Lock()
	defer pthis.mapDialDataMutex.Unlock()

	dd, ok := pthis.mapDialData[srvID]
	if !ok || dd == nil {
		return nil
	}
	return dd
}

func (pthis *TcpDialMgr) delDialData(srvID int64) {
	pthis.mapDialDataMutex.Lock()
	defer pthis.mapDialDataMutex.Unlock()

	delete(pthis.mapDialData, srvID)
}

func (pthis *TcpDialMgr) TcpDialMgrDial(srvID int64, ip string, port int, closeExpireSec int,
	dialTimeoutSec int,
	needRedial bool, redialCount int) (*TcpConnection, error) {

	oldDial := pthis.getDialData(srvID)
	if oldDial != nil {
		if oldDial.DialCancelFunc != nil {
			oldDial.DialCancelFunc()
		}
	}
	ctx, canelfun := context.WithCancel(context.Background())
	dd := pthis.addDialData(srvID,
		&dialData{
			dialTimeoutSec:dialTimeoutSec,
			closeExpireSec:closeExpireSec,
			ip:ip,
			port:port,
			srvID:srvID,
			needRedial:needRedial,
			redialCount:redialCount,
			DialCancelFunc:canelfun})

	tcpConn, err := startTcpDial(pthis.connMgr, srvID, ip, port, closeExpireSec, dialTimeoutSec, redialCount, ctx)
	if err != nil {
		return nil, err
	}
	dd.tcpConnID = tcpConn.TcpConnectionID()

	return tcpConn, nil
}

/*func (pthis *TcpDialMgr) TcpDialDelDialData(srvID int64) {
	pthis.delDialData(srvID)
}*/
