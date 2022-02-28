package main

import (
	"lin/lin_common"
	"lin/log"
	"net"
	"strconv"
	"sync"
)

type TCP_CONNECTION_ID int64
type MAP_TCPCONN map[TCP_CONNECTION_ID]*TcpConnection

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

		tcpConn, err := StartTcpConnection(pthis, pthis.TcpAcceptGenConnectionID(), conn, pthis.closeExpireSec, pthis.cbConnection)
		if err != nil {
			log.LogErr("start accept tcp connect err", err)
		}
		pthis.TcpAcceptAddTcpConn(tcpConn)
	}

	pthis.wg.Done()
}

func StartTcpAccept(ip string, port int, CBConnection InterfaceTcpConnection,  closeExpireSec int) (*TcpAccept, error) {
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

func (pthis * TcpAccept) TcpAcceptGenConnectionID() TCP_CONNECTION_ID {
	return TCP_CONNECTION_ID(lin_common.GenUUID64_V4())
}

func (pthis * TcpAccept) TcpAcceptWait() {
	log.LogDebug("begin wait")
	pthis.wg.Wait()
	log.LogDebug("end wait")
}

func (pthis * TcpAccept) TcpAcceptAddTcpConn(tcpConn *TcpConnection) {
	pthis.mapConnMutex.Lock()
	defer pthis.mapConnMutex.Unlock()

	pthis.mapConn[tcpConn.TcpConnectionID()] = tcpConn
}

func (pthis * TcpAccept) delTcpConn(id TCP_CONNECTION_ID) {
	pthis.mapConnMutex.Lock()
	defer pthis.mapConnMutex.Unlock()

	delete(pthis.mapConn, id)
}

func (pthis * TcpAccept) TcpAcceptCloseConn(id TCP_CONNECTION_ID) {
	pthis.mapConnMutex.Lock()
	defer pthis.mapConnMutex.Unlock()

	conn, ok := pthis.mapConn[id]
	if !ok || conn == nil {
		return
	}

	conn.TcpConnectClose()
}
