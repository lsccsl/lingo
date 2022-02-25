package main

import (
	"lin/lin_common"
	"lin/log"
	"net"
	"strconv"
	"sync"
)

type MAP_TCPCONN map[int64]*TcpConnection

type TcpSrv struct {
	tcpLsn       net.Listener
	wg           sync.WaitGroup
	cbConnection InterfaceTcpConnection
	closeExpireSec int

	mapConnMutex sync.Mutex
	mapConn MAP_TCPCONN
}

func (pthis * TcpSrv)go_tcpAccept() {
	for {
		conn, err := pthis.tcpLsn.Accept()
		if err != nil {
			log.LogErr("tcp accept err", err)
		}

		tcpConn, err := StartTcpAcceptClient(pthis, pthis.TcpSrvGenClientID(), conn, pthis.closeExpireSec, pthis.cbConnection)
		if err != nil {
			log.LogErr("start accept tcp connect err", err)
		}
		pthis.TcpSrvAddTcpConn(tcpConn)
	}

	pthis.wg.Done()
}

func StartTcpListener(ip string, port int, CBConnection InterfaceTcpConnection,  closeExpireSec int) (*TcpSrv, error) {
	ts := &TcpSrv{}

	addr := ip + ":" + strconv.Itoa(port)
	lsn, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	ts.tcpLsn = lsn
	ts.cbConnection = CBConnection
	ts.closeExpireSec = closeExpireSec
	ts.mapConn = make(MAP_TCPCONN)

	ts.wg.Add(1)
	go ts.go_tcpAccept()

	return ts, nil
}

func (pthis * TcpSrv) TcpSrvGenClientID() int64 {
	return lin_common.GenUUID64_V4()
}

func (pthis * TcpSrv) TcpSrvWait() {
	log.LogDebug("begin wait")
	pthis.wg.Wait()
	log.LogDebug("end wait")
}

func (pthis * TcpSrv) TcpSrvAddTcpConn(tcpConn *TcpConnection) {
	pthis.mapConnMutex.Lock()
	pthis.mapConn[tcpConn.TcpClientID()] = tcpConn
	pthis.mapConnMutex.Unlock()
}

func (pthis * TcpSrv) TcpSrvDelTcpConn(id int64) {
	pthis.mapConnMutex.Lock()
	delete(pthis.mapConn, id)
	pthis.mapConnMutex.Unlock()
}
