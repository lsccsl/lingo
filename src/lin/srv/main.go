package main

import (
	"fmt"
	"lin/log"
	"net"
	"net/http"
	"os"
)


func main() {
	fmt.Println(os.Args)

	var pathCfg string
	if len(os.Args) >= 2 {
		pathCfg = os.Args[1]
	}
	ReadCfg(pathCfg)

	InitMsgParseVirtualTable()
	server := ConstructServerMgr(Global_ServerCfg.SrvID)

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
	tcpMgr, err := StartTcpManager(tcpAddr.IP.String(), tcpAddr.Port, server, 180)
	if err != nil {
		log.LogErr(err)
		return
	}
	log.LogDebug(tcpMgr)

	server.tcpMgr = tcpMgr
	server.httpSrv = httpSrv

	if len(Global_ServerCfg.MapCluster) > 1 {
		for key, val := range Global_ServerCfg.MapCluster {
			if val == Global_ServerCfg.BindAddr {
				continue
			}
			dialAddr, err := net.ResolveTCPAddr("tcp", val)
			if err != nil {
				log.LogErr(err)
				return
			}
			tcpMgr.TcpDialMgrDial(int64(key), dialAddr.IP.String(), dialAddr.Port, 180, 15, true, 10)
			log.LogDebug(val)
		}
	}

	tcpMgr.TcpMgrWait()
}

