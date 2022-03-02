package main

import (
	"fmt"
	"lin/log"
	"net/http"
)


func main() {
	ReadCfg()
	InitMsgParseVirtualTable()
	server := ConstructServer()

	httpSrv, err := StartHttpSrvMgr("0.0.0.0", 8112)
	if err != nil {
		log.LogErr(err)
	}
	httpSrv.HttpSrvAddCallback("/test", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, request.URL.Path, " ", request.Form)
	})

	tcpAccept, err := StartTcpAccept("0.0.0.0", 1128, server, 180)
	if err != nil {
		log.LogErr(err)
		return
	}
	log.LogDebug(tcpAccept)
	server.accept = tcpAccept
	server.httpSrv = httpSrv

	tcpAccept.TcpAcceptWait()
}

