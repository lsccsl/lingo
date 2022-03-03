package main

import (
	"lin/log"
	"sync/atomic"
)

type Server struct {
	connDial *TcpConnection
	connAcpt *TcpConnection
	chServerMsg chan *interMsg

	isStopProcess int32
}

func ConstructServer(connDial *TcpConnection, connAcpt *TcpConnection)*Server {
	s := &Server{
		connDial:connDial,
		connAcpt:connAcpt,
		chServerMsg:make(chan *interMsg),
		isStopProcess:0,
	}
	go s.go_serverProcess()
	return s
}

func (pthis*Server) go_serverProcess() {
	defer func() {
		err := recover()
		if err != nil {
			log.LogErr(err)
		}
	}()

MSG_LOOP:
	for {
		select {
		case clientMsg := <- pthis.chServerMsg:
			if clientMsg == nil {
				break MSG_LOOP
			}
			//pthis.processServerMsg(clientMsg)
		}
	}

	atomic.StoreInt32(&pthis.isStopProcess, 1)
	close(pthis.chServerMsg)
}