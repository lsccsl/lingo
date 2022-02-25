package main

import (
	"bytes"
	"fmt"
	"net"
)

const G_MTU int = 1536
const MAX_PACK_LEN int = 65535

type InterfaceTcpConnection interface {
	CBReadProcess(recvBuf * bytes.Buffer)(bytesProcess int)
	CBConnect(tcpConn * TcpConnection)
}

type interMsgTcpWrite struct {
	bin []byte
}

type TcpConnection struct {
	clientConn net.Conn
	CBTcpConnection InterfaceTcpConnection
	chMsgWrite chan *interMsgTcpWrite
}

func StartAcceptTcpConnect(conn net.Conn, CBTcpConnection InterfaceTcpConnection) (*TcpConnection, error) {
	tcpConn := &TcpConnection{
		clientConn:conn,
		CBTcpConnection:CBTcpConnection,
	}

	if tcpConn.CBTcpConnection != nil {
		tcpConn.CBTcpConnection.CBConnect(tcpConn)
	}

	go tcpConn.go_tcpConnRead()
	go tcpConn.go_tcpConnWrite()

	return tcpConn, nil
}

func (pthis * TcpConnection)go_tcpConnRead() {
	TmpBuf := make([]byte, G_MTU)
	recvBuf := bytes.NewBuffer(make([]byte, 0, MAX_PACK_LEN))

	READ_LOOP:
	for {
		readSize, err := pthis.clientConn.Read(TmpBuf)
		if err != nil {
			break READ_LOOP
		}
		recvBuf.Write(TmpBuf[0:readSize])

		if pthis.CBTcpConnection == nil {
			recvBuf.Next(readSize)
			continue
		}

		bytesProcess := pthis.CBTcpConnection.CBReadProcess(recvBuf)
		if bytesProcess < 0 {
			break READ_LOOP
		} else if bytesProcess > 0 {
			recvBuf.Next(bytesProcess)
		}
	}
}

func (pthis * TcpConnection)go_tcpConnWrite() {
	WRITE_LOOP:
	for {
		select {
		case tcpW := <- pthis.chMsgWrite:
			if tcpW == nil {
				break WRITE_LOOP
			}
			pthis.clientConn.Write(tcpW.bin)
		}
	}
}

func (pthis * TcpConnection)TcpWrite(bin []byte) {
	tcpW := &interMsgTcpWrite{
		make([]byte,0,len(bin)),
	}
	copy(tcpW.bin, bin)
	fmt.Println(&tcpW.bin, &bin)
	pthis.chMsgWrite <- tcpW
	//tcpW.bin = append(tcpW.bin, bin...)
}
