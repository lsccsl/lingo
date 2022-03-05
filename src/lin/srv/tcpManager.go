package main

import (
	"lin/lin_common"
	"lin/log"
	"net"
	"strconv"
	"sync"
)

type TcpMgr struct {
	tcpLsn       net.Listener
	wg           sync.WaitGroup
	cbConnection InterfaceTcpConnection
	closeExpireSec int

	TcpDialMgr

	mapConnMutex sync.Mutex
	mapConn MAP_TCPCONN
}

func (pthis * TcpMgr) CBGenConnectionID() TCP_CONNECTION_ID {
	return TCP_CONNECTION_ID(lin_common.GenUUID64_V4())
}
func (pthis * TcpMgr) CBAddTcpConn(tcpConn *TcpConnection) {
	pthis.mapConnMutex.Lock()
	defer pthis.mapConnMutex.Unlock()

	pthis.mapConn[tcpConn.TcpConnectionID()] = tcpConn
}
func (pthis * TcpMgr) CBGetConnectionCB()InterfaceTcpConnection {
	return pthis.cbConnection
}
func (pthis * TcpMgr) CBDelTcpConn(id TCP_CONNECTION_ID) {
	pthis.mapConnMutex.Lock()
	defer pthis.mapConnMutex.Unlock()

	delete(pthis.mapConn, id)
}

func (pthis * TcpMgr)go_tcpAccept() {
	for {
		conn, err := pthis.tcpLsn.Accept()
		if err != nil {
			log.LogErr("tcp accept err", err)
		}

		log.LogDebug(conn.LocalAddr(), conn.RemoteAddr())

		_, err = startTcpConnection(pthis, conn, pthis.closeExpireSec)
		if err != nil {
			log.LogErr("start accept tcp connect err", err)
		}
	}

	pthis.wg.Done()
}

func StartTcpManager(ip string, port int, CBConnection InterfaceTcpConnection,  closeExpireSec int) (*TcpMgr, error) {
	t := &TcpMgr{}

	addr := ip + ":" + strconv.Itoa(port)
	lsn, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	t.tcpLsn = lsn
	t.cbConnection = CBConnection
	t.closeExpireSec = closeExpireSec
	t.mapConn = make(MAP_TCPCONN)

	t.wg.Add(1)
	go t.go_tcpAccept()

	t.TcpDialMgrStart(t, closeExpireSec)

	return t, nil
}

func (pthis * TcpMgr) TcpMgrWait() {
	log.LogDebug("begin wait")
	pthis.wg.Wait()
	log.LogDebug("end wait")
}


func (pthis * TcpMgr) TcpMgrCloseConn(id TCP_CONNECTION_ID) {
	pthis.mapConnMutex.Lock()
	defer pthis.mapConnMutex.Unlock()

	conn, ok := pthis.mapConn[id]
	if !ok || conn == nil {
		return
	}

	conn.TcpConnectClose()
}
