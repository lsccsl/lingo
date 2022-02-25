package main

import (
	"lin/log"
	"net"
	"strconv"
	"sync"
)

type TcpSrv struct {
	tcpLsn       net.Listener
	wg           sync.WaitGroup
	CBConnection InterfaceTcpConnection
}

func (pthis * TcpSrv)go_tcpAccept() {
	for {
		conn, err := pthis.tcpLsn.Accept()
		if err != nil {
			log.LogErr("tcp accept err", err)
		}

		_, err = StartAcceptTcpConnect(conn, pthis.CBConnection)
		if err != nil {
			log.LogErr("start accept tcp connect err", err)
		}
	}

	pthis.wg.Done()
}

func StartTcpListener(ip string, port int, CBConnection InterfaceTcpConnection) (*TcpSrv, error) {
	ts := &TcpSrv{}

	addr := ip + ":" + strconv.Itoa(port)
	lsn, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	ts.tcpLsn = lsn
	ts.CBConnection = CBConnection

	ts.wg.Add(1)
	go ts.go_tcpAccept()

	return ts, nil
}

func (pthis * TcpSrv) TcpSrvWait() {
	pthis.wg.Wait()
}
