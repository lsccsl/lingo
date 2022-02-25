package main

import (
	"bytes"
	"net"
)

const G_MTU int = 1536
const MAX_PACK_LEN int = 65535

type TcpConnection struct {
	clientConn net.Conn
}

func StartAcceptTcpConnect(conn net.Conn) (*TcpConnection, error) {
	tcpConn := &TcpConnection{
		clientConn:conn,
	}
	tcpConn.go_tcpConnRead()

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
	}
}
