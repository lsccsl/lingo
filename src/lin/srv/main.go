package main

import (
	"fmt"
	"lin/log"
	"net"
	"net/http"
)


func main() {
	ReadCfg()
	InitMsgParseVirtualTable()
	server := ConstructServer()

	httpAddr, err := net.ResolveTCPAddr("tcp", Global_ServerCfg.HttpAddr)
	if err != nil {
		log.LogErr(err)
		return
	}
	httpSrv, err := StartHttpSrvMgr(httpAddr.IP.String(), httpAddr.Port)
	if err != nil {
		log.LogErr(err)
	}

	httpSrv.HttpSrvAddCallback("/test", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, request.URL.Path, " ", request.Form)
	})

	tcpAddr, err := net.ResolveTCPAddr("tcp", Global_ServerCfg.BindAddr)
	if err != nil {
		log.LogErr(err)
		return
	}
	tcpAccept, err := StartTcpAccept(tcpAddr.IP.String(), tcpAddr.Port, server, 180)
	if err != nil {
		log.LogErr(err)
		return
	}
	log.LogDebug(tcpAccept)

	dialMgr, err := StartTcpDial(180)
	if err != nil {
		log.LogErr(err)
		return
	}

	server.dialMgr = dialMgr
	server.accept = tcpAccept
	server.httpSrv = httpSrv

	if len(Global_ServerCfg.MapCluster) > 1 {
		for key, val := range Global_ServerCfg.MapCluster {
			if val == Global_ServerCfg.BindAddr {
				continue
			}
			dialAddr, err := net.ResolveTCPAddr("tcp", Global_ServerCfg.BindAddr)
			if err != nil {
				log.LogErr(err)
				return
			}
			dialMgr.TcpDialMgrDial(int64(key), dialAddr.IP.String(), dialAddr.Port, 180, 30)
			log.LogDebug(val)
		}
	}

	tcpAccept.TcpAcceptWait()
}

