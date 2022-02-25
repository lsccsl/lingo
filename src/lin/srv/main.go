package main

import (
	"bytes"
	"fmt"
	"lin/log"
)
type CallbackTcpConnection struct {
}

func (pthis*CallbackTcpConnection)CBReadProcess(recvBuf * bytes.Buffer)(bytesProcess int){
	log.LogDebug("len", recvBuf.Len())
	return recvBuf.Len()
}
func (pthis*CallbackTcpConnection)CBConnect(tcpConn * TcpConnection){
	log.LogDebug()
}

func main() {
	cb := &CallbackTcpConnection{}

	tcpLsn, err := StartTcpListener("0.0.0.0", 1122, cb)
	if err != nil {
		fmt.Println(err)
	}
	log.LogDebug(tcpLsn)

	tcpLsn.TcpSrvWait()
}

