package main

import (
	"bytes"
	"encoding/binary"
	"lin/lin_common"
	"lin/log"
	"lin/msg"
	"sync"
)

type dialData struct {
	dialTimeoutSec int
	closeExpireSec int
	tcpConn *TcpConnection
	srvID int64
	ip string
	port int
}
type MAP_DIALDATA map[int64/* srvID */]*dialData

type TcpDialMgr struct {
	wg sync.WaitGroup
	closeExpireSec int
	chRedial chan int64
	mapDialData MAP_DIALDATA

	mapConnMutex sync.Mutex
	mapConn MAP_TCPCONN
}


func StartTcpDial(closeExpireSec int) (*TcpDialMgr, error) {
	tm := &TcpDialMgr{
		closeExpireSec:closeExpireSec,
		chRedial: make(chan int64),
		mapDialData: make(MAP_DIALDATA),
	}

	tm.wg.Add(1)
	go tm.go_DialCheck()

	return tm, nil
}

func (pthis * TcpDialMgr) CBReadProcess(tcpConn * TcpConnection, recvBuf * bytes.Buffer)(bytesProcess int){
	if recvBuf.Len() < 6 {
		return 0
	}
	binHead := recvBuf.Bytes()[0:6]

	packLen := binary.LittleEndian.Uint32(binHead[0:4])
	packType := binary.LittleEndian.Uint16(binHead[4:6])

	log.LogDebug("packLen:", packLen, " packType:", packType)

	if recvBuf.Len() < int(packLen){
		return 0
	}

	binBody := recvBuf.Bytes()[6:packLen]

	protoMsg := ParseProtoMsg(binBody, int32(packType))
	if protoMsg == nil {
		return int(packLen)
	}
	log.LogDebug("parse protomsg:", protoMsg)


	switch msg.MSG_TYPE(packType) {
	case msg.MSG_TYPE__MSG_RPC:
		t, ok := protoMsg.(*msg.MSG_RPC)
		if ok && t != nil {
		}
	case msg.MSG_TYPE__MSG_RPC_RES:
		t, ok := protoMsg.(*msg.MSG_RPC_RES)
		if ok && t != nil {
		}
	}

	return int(packLen)
}
func (pthis * TcpDialMgr) CBConnect(tcpConn * TcpConnection, err error){
	if err != nil {
		log.LogErr(err)
	}
	if tcpConn == nil {
		return
	}
	log.LogDebug(tcpConn.TcpGetConn().LocalAddr(), tcpConn.TcpGetConn().RemoteAddr(), tcpConn.TcpConnectionID())
}
func (pthis * TcpDialMgr) CBConnectClose(id TCP_CONNECTION_ID){
	log.LogDebug("id:", id)
}

func (pthis * TcpDialMgr) CBGenConnectionID() TCP_CONNECTION_ID {
	return TCP_CONNECTION_ID(lin_common.GenUUID64_V4())
}
func (pthis * TcpDialMgr) CBAddTcpConn(tcpConn *TcpConnection) {
	pthis.mapConnMutex.Lock()
	defer pthis.mapConnMutex.Unlock()

	pthis.mapConn[tcpConn.TcpConnectionID()] = tcpConn
}
func (pthis * TcpDialMgr) CBGetConnectionCB()InterfaceTcpConnection {
	return pthis
}
func (pthis * TcpDialMgr) CBDelTcpConn(id TCP_CONNECTION_ID) {
	pthis.mapConnMutex.Lock()
	defer pthis.mapConnMutex.Unlock()

	delete(pthis.mapConn, id)
}


func (pthis * TcpDialMgr) go_DialCheck() {
	CHECK_LOOP:
	for {
		select {
		case id := <-pthis.chRedial:
			if id < 0 {
				break CHECK_LOOP
			}
			// redial
			pthis.reDialMgrDial(id)
		}
	}

	pthis.wg.Done()
}

func (pthis * TcpDialMgr)TcpDialMgrWait() {
	pthis.wg.Wait()
}

func (pthis * TcpDialMgr) TcpDialMgrDial(srvID int64, ip string, port int, closeExpireSec int, dialTimeoutSec int) (*TcpConnection, error) {
	tcpConn, err := startTcpDial(pthis, ip, port, closeExpireSec, dialTimeoutSec)
	if err != nil {
		return nil, err
	}

	tcpConn.AppID = srvID

	pthis.addDialData(srvID,
		&dialData{
			dialTimeoutSec:dialTimeoutSec,
			closeExpireSec:closeExpireSec,
			tcpConn:tcpConn,
			ip:ip,
			port:port,
			srvID:srvID,})

	return tcpConn, nil
}

func (pthis * TcpDialMgr) addDialData(srvID int64, dd *dialData) {
	pthis.mapConnMutex.Lock()
	defer pthis.mapConnMutex.Unlock()

	pthis.mapDialData[srvID] = dd
}

func (pthis * TcpDialMgr) getDialData(srvID int64, ddOut *dialData) bool {
	dd, ok := pthis.mapDialData[srvID]
	if !ok || dd == nil {
		return false
	}
	*ddOut = *dd
	ddOut.tcpConn = nil
	return true
}

func (pthis * TcpDialMgr) reDialMgrDial(srvID int64) (*TcpConnection, error) {
	var dd dialData
	bret := pthis.getDialData(srvID, &dd)
	if !bret {
		return nil, lin_common.GenErr(lin_common.ERR_no_dialData)
	}

	return pthis.TcpDialMgrDial(dd.srvID, dd.ip, dd.port, dd.closeExpireSec, dd.dialTimeoutSec)
}