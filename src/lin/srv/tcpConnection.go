package main

import (
	"bytes"
	"fmt"
	"lin/log"
	"net"
	"sync/atomic"
	"time"
)

const G_MTU int = 1536
const MAX_PACK_LEN int = 65535

type InterfaceTcpConnection interface {
	CBReadProcess(pthis * TcpConnection, recvBuf * bytes.Buffer)(bytesProcess int)
	CBConnect(tcpConn * TcpConnection)
	CBConnectClose(id TCP_CONNECTION_ID)
}

type interMsgTcpWrite struct {
	bin []byte
}

type TcpConnection struct {
	connectionID TCP_CONNECTION_ID
	clientAppID int64
	clientConn net.Conn
	cbTcpConnection InterfaceTcpConnection
	chMsgWrite chan *interMsgTcpWrite
	closeExpireSec int
	tcpAccept *TcpAccept

	canWrite int32
}

func StartTcpConnection(tcpAccept *TcpAccept, connectionID TCP_CONNECTION_ID, conn net.Conn, closeExpireSec int, CBTcpConnection InterfaceTcpConnection) (*TcpConnection, error) {
	tcpConn := &TcpConnection{
		connectionID:connectionID,
		clientConn:conn,
		cbTcpConnection:CBTcpConnection,
		closeExpireSec:closeExpireSec,
		tcpAccept:tcpAccept,
		canWrite:0,
	}

	if tcpConn.cbTcpConnection != nil {
		tcpConn.cbTcpConnection.CBConnect(tcpConn)
	}

	go tcpConn.go_tcpConnRead()
	go tcpConn.go_tcpConnWrite()

	return tcpConn, nil
}

func (pthis * TcpConnection)go_tcpConnRead() {
	defer func() {
		pthis.cbTcpConnection.CBConnectClose(pthis.connectionID)
		pthis.tcpAccept.delTcpConn(pthis.connectionID)

		err := recover()
		if err != nil {
			log.LogErr(err)
		}
	}()

	TmpBuf := make([]byte, G_MTU)
	recvBuf := bytes.NewBuffer(make([]byte, 0, MAX_PACK_LEN))

	expireInterval := time.Second * time.Duration(pthis.closeExpireSec)
	var TimerConnClose * time.Timer = nil
	if pthis.closeExpireSec > 0 {
		TimerConnClose = time.AfterFunc(expireInterval, func() {
			pthis.TcpConnectClose()
		})
	}

	READ_LOOP:
	for {
		readSize, err := pthis.clientConn.Read(TmpBuf)
		if err != nil {
			break READ_LOOP
		}

		if pthis.closeExpireSec > 0 {
			TimerConnClose.Reset(expireInterval)
		}

		recvBuf.Write(TmpBuf[0:readSize])

		if pthis.cbTcpConnection == nil {
			recvBuf.Next(readSize)
			continue
		}

		bytesProcess := pthis.cbTcpConnection.CBReadProcess(pthis, recvBuf)
		if bytesProcess < 0 {
			break READ_LOOP
		} else if bytesProcess > 0 {
			recvBuf.Next(bytesProcess)
		}
	}
}

func (pthis * TcpConnection)go_tcpConnWrite() {
	defer func() {
		err := recover()
		if err != nil {
			log.LogErr(err)
		}
	}()

	WRITE_LOOP:
	for {
		select {
		case tcpW := <- pthis.chMsgWrite:
			if tcpW == nil {
				break WRITE_LOOP
			}
			//todo: option wait for more data and combine write to tcp channel
			pthis.clientConn.Write(tcpW.bin)
		}
	}

	atomic.StoreInt32(&pthis.canWrite, 1)
	close(pthis.chMsgWrite)
}

func (pthis * TcpConnection)TcpConnectWrite(bin []byte) {
	if atomic.LoadInt32(&pthis.canWrite) != 0 {
		return
	}

	tcpW := &interMsgTcpWrite{
		make([]byte,0,len(bin)),
	}
	copy(tcpW.bin, bin)
	fmt.Println(&tcpW.bin, &bin)
	pthis.chMsgWrite <- tcpW
	//tcpW.bin = append(tcpW.bin, bin...)
}

func (pthis * TcpConnection)TcpGetConn() net.Conn {
	return pthis.clientConn
}

func (pthis * TcpConnection)TcpConnectClose() {
	if atomic.LoadInt32(&pthis.canWrite) != 0 {
		pthis.chMsgWrite <- nil
	}
	pthis.clientConn.Close()
}

func (pthis * TcpConnection)TcpConnectionID() TCP_CONNECTION_ID {
	return pthis.connectionID
}

func (pthis * TcpConnection)TcpConnectClientAppID() int64 {
	return pthis.clientAppID
}
func (pthis * TcpConnection)TcpConnectSetClientAppID(id int64) {
	pthis.clientAppID = id
}