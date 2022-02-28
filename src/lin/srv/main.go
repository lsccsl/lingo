package main

import (
	"lin/log"
)


func main() {
	InitMsgParseVirtualTable()
	server := ConstructServer()

	tcpAccept, err := StartTcpAccept("0.0.0.0", 1126, server, 30)
	if err != nil {
		log.LogErr(err)
		return
	}
	log.LogDebug(tcpAccept)
	server.accept = tcpAccept

	tcpAccept.TcpAcceptWait()
}

