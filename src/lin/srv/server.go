package main

import (
	"github.com/golang/protobuf/proto"
	"lin/log"
	msgpacket "lin/msgpacket"
	"sync/atomic"
	"time"
)

type Server struct {
	srvMgr *ServerMgr
	srvID int64
	connDial *TcpConnection
	connAcpt *TcpConnection
	chSrvProtoMsg chan *interProtoMsg
	chInterMsg chan interface{}
	heartbeatIntervalSec int

	isStopProcess int32
}

type interMsgSrvReport struct {
	tcpAccept * TcpConnection
}
type interMsgConnDial struct {
	tcpDial * TcpConnection
}

func ConstructServer(srvMgr *ServerMgr, srvID int64, heartbeatIntervalSec int)*Server {
	s := &Server{
		srvMgr:srvMgr,
		srvID:srvID,
		chSrvProtoMsg:make(chan *interProtoMsg, 100),
		chInterMsg:make(chan interface{}, 100),
		heartbeatIntervalSec:heartbeatIntervalSec,
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

	chTimer := time.After(time.Second * time.Duration(pthis.heartbeatIntervalSec))

	MSG_LOOP:
	for {
		select {
		case ProtoMsg := <- pthis.chSrvProtoMsg:
			{
				if ProtoMsg == nil {
					break MSG_LOOP
				}
				//pthis.processServerMsg(clientMsg)
			}

		case interMsg := <- pthis.chInterMsg:
			{
				switch t:= interMsg.(type) {
				case *interMsgSrvReport:
					pthis.processSrvReport(t.tcpAccept)
				case *interMsgConnDial:
					pthis.processDailConnect(t.tcpDial)
				}
			}

		case <-chTimer:
			{
				chTimer = time.After(time.Second * time.Duration(pthis.heartbeatIntervalSec))
				//send heartbeat
				msgTest := &msgpacket.MSG_TEST{}
				msgTest.Id = pthis.srvID
				pthis.connDial.TcpConnectSendProtoMsg(msgpacket.MSG_TYPE__MSG_TEST, msgTest)
			}
		}
	}

	atomic.StoreInt32(&pthis.isStopProcess, 1)
	close(pthis.chSrvProtoMsg)
	close(pthis.chInterMsg)
}

func (pthis*Server) ServerClose() {
	if pthis.connAcpt != nil {
		pthis.srvMgr.tcpMgr.TcpMgrCloseConn(pthis.connAcpt.TcpConnectionID())
	}
	if pthis.connDial != nil {
		pthis.srvMgr.tcpMgr.TcpMgrCloseConn(pthis.connDial.TcpConnectionID())
	}
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
func (pthis*Server)PushProtoMsg(msgType msgpacket.MSG_TYPE, protoMsg proto.Message){
	if atomic.LoadInt32(&pthis.isStopProcess) == 1 {
		return
	}
	pthis.chSrvProtoMsg <- &interProtoMsg{
		msgType:msgType,
		protoMsg:protoMsg,
	}
}

func (pthis*Server)processRPC(tcpConn * TcpConnection, msg proto.Message) proto.Message {
	switch t:= msg.(type) {
	case *msgpacket.MSG_TEST:
		{
			return pthis.processRPCTest(tcpConn, t)
		}
	}
	return nil
}
func (pthis*Server)processRPCRes(tcpConn * TcpConnection, msg proto.Message) {
}

func (pthis*Server)processRPCTest(tcpDial * TcpConnection, msg *msgpacket.MSG_TEST) *msgpacket.MSG_TEST_RES {
	return &msgpacket.MSG_TEST_RES{Id: 123}
}
