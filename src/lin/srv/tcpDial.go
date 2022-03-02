package main

import (
	"bytes"
	"encoding/binary"
	"lin/lin_common"
	"lin/log"
	"lin/msg"
	"sync"
	"time"
)

type TcpDialMgr struct {
	wg sync.WaitGroup
	closeExpireSec int
	chQuit chan interface{}

	mapConnMutex sync.Mutex
	mapConn MAP_TCPCONN
}


func StartTcpDial(closeExpireSec int) (*TcpDialMgr, error) {
	tm := &TcpDialMgr{
		closeExpireSec:closeExpireSec,
		chQuit: make(chan interface{}),
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
	chTimeout := time.After(time.Second * time.Duration(30))
	CHECK_LOOP:
	for {
		select {
		case <-chTimeout:
		case <-pthis.chQuit:
			break CHECK_LOOP
		}
	}

	pthis.wg.Done()
}

func (pthis * TcpDialMgr)TcpDialMgrWait() {
	pthis.wg.Wait()
}

func (pthis * TcpDialMgr) TcpDialMgrDial(ip string, port int, closeExpireSec int, dialTimeoutSec int) (*TcpConnection, error) {
	tcpConn, err := startTcpDial(pthis, ip, port, closeExpireSec, dialTimeoutSec)
	if err != nil {
		return nil, err
	}
	return tcpConn, nil
}