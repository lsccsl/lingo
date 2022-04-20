package main

import (
	"lin/lin_common"
	"net"
)

type TcpClient struct {
	fd lin_common.FD_DEF
	addr *net.Addr

	clientID int64
}

func ConstructorTcpClient(fd lin_common.FD_DEF) *TcpClient {
	tc := &TcpClient{
		fd : fd,
	}

	return tc
}