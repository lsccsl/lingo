package main

import (
	"lin/lin_common"
	"lin/log"
	"net"
	"strconv"
	"sync"
)

type MAP_TCPCONN map[int64]*TcpConnection

type TcpAccept struct {
	tcpLsn       net.Listener
	wg           sync.WaitGroup
	cbConnection InterfaceTcpConnection
	closeExpireSec int

	mapConnMutex sync.Mutex
	mapConn MAP_TCPCONN
}

func (pthis * TcpAccept)go_tcpAccept() {
	for {
		conn, err := pthis.tcpLsn.Accept()
		if err != nil {
			log.LogErr("tcp accept err", err)
		}

		tcpConn, err := StartTcpConnection(pthis, pthis.TcpAcceptGenClientID(), conn, pthis.closeExpireSec, pthis.cbConnection)
		if err != nil {
			log.LogErr("start accept tcp connect err", err)
		}
		pthis.TcpAcceptAddTcpConn(tcpConn)
	}

	pthis.wg.Done()
}

func StartTcpListener(ip string, port int, CBConnection InterfaceTcpConnection,  closeExpireSec int) (*TcpAccept, error) {
	ts := &TcpAccept{}

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

func (pthis * TcpAccept) TcpAcceptGenClientID() int64 {
	return lin_common.GenUUID64_V4()
}

func (pthis * TcpAccept) TcpAcceptWait() {
	log.LogDebug("begin wait")
	pthis.wg.Wait()
	log.LogDebug("end wait")
}

func (pthis * TcpAccept) TcpAcceptAddTcpConn(tcpConn *TcpConnection) {
	pthis.mapConnMutex.Lock()
	pthis.mapConn[tcpConn.TcpClientID()] = tcpConn
	pthis.mapConnMutex.Unlock()
}

func (pthis * TcpAccept) TcpAcceptDelTcpConn(id int64) {
	pthis.mapConnMutex.Lock()
	delete(pthis.mapConn, id)
	pthis.mapConnMutex.Unlock()
}
