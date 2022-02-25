package main

import (
	"lin/log"
)


func main() {
	InitMsgParseVirtualTable()
	srvMgr := ConstructSrvManager()

	tcpLsn, err := StartTcpListener("0.0.0.0", 1123, srvMgr, 30)
	if err != nil {
		log.LogErr(err)
		return
	}
	log.LogDebug(tcpLsn)

	tcpLsn.TcpSrvWait()
}

