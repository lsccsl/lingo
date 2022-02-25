package main

import (
	"fmt"
	"lin/log"
)

func main() {

	log.LogDebug("test", "bbb")
	log.LogDebug("test", "aaa")

	tcpLsn, err := StartTcpListener("0.0.0.0", 1122)
	if err != nil {
		fmt.Println(err)
	}
	log.LogDebug(tcpLsn)

	tcpLsn.TcpSrvWait()
}

