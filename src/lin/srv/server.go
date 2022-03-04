package main

import (
	"lin/log"
	"sync/atomic"
)

type Server struct {
	srvMgr *ServerMgr
	srvID int64
	connDial *TcpConnection
	connAcpt *TcpConnection
	chSrvProtoMsg chan *interProtoMsg
	chInterMsg chan interface{}

	isStopProcess int32
}

type interMsgSrvReport struct {
	tcpAccept * TcpConnection
}
type interMsgConnDial struct {
	tcpDial * TcpConnection
}

func ConstructServer(srvMgr *ServerMgr, srvID int64)*Server {
	s := &Server{
		srvMgr:srvMgr,
		srvID:srvID,
		chSrvProtoMsg:make(chan *interProtoMsg),
		chInterMsg:make(chan interface{}),
		isStopProcess:0,
	}
	go s.go_serverProcess()
	return s
}

func (pthis*Server) go_serverProcess() {
	defer func() {
		atomic.StoreInt32(&pthis.isStopProcess, 1)
		err := recover()
		if err != nil {
			log.LogErr(err)
		}
	}()

MSG_LOOP:
	for {
		select {
		case ProtoMsg := <- pthis.chSrvProtoMsg:
			if ProtoMsg == nil {
				break MSG_LOOP
			}
			//pthis.processServerMsg(clientMsg)

		case interMsg := <- pthis.chInterMsg:
			{
				switch t:= interMsg.(type) {
				case *interMsgSrvReport:
					pthis.processSrvReport(t.tcpAccept)
				case *interMsgConnDial:
					pthis.processDailConnect(t.tcpDial)
				}
			}
		}
	}

	atomic.StoreInt32(&pthis.isStopProcess, 1)
	close(pthis.chSrvProtoMsg)
	close(pthis.chInterMsg)
}

func (pthis*Server) ServerClose() {
	pthis.srvMgr.tcpMgr.TcpMgrCloseConn(pthis.connAcpt.TcpConnectionID())
	pthis.srvMgr.tcpMgr.TcpMgrCloseConn(pthis.connDial.TcpConnectionID())
	pthis.chSrvProtoMsg <- nil
}

func (pthis*Server) ServerCloseAndDelDialData() {
	pthis.srvMgr.tcpMgr.TcpDialDelDialData(pthis.srvID)
	pthis.ServerClose()
}

func (pthis*Server)processSrvReport(tcpAccept * TcpConnection){
	pthis.connAcpt = tcpAccept

	log.LogDebug(pthis.srvID, " ", pthis)
}

func (pthis*Server)processDailConnect(tcpDial * TcpConnection){
	pthis.connDial = tcpDial

	log.LogDebug(pthis.srvID, " ", pthis)
}

func (pthis*Server)PushInterMsg(msg interface{}){
	if atomic.LoadInt32(&pthis.isStopProcess) == 1 {
		return
	}
	pthis.chInterMsg <- msg
}