package main

import (
	"fmt"
	"lin/log"
	"lin/msgpacket"
	"net"
	"net/http"
	"os"
)

var srvMgr *ServerMgr

func main() {
	fmt.Println(os.Args)

	var pathCfg string
	if len(os.Args) >= 2 {
		pathCfg = os.Args[1]
	}
	ReadCfg(pathCfg)

	msgpacket.InitMsgParseVirtualTable()
	srvMgr = ConstructServerMgr(Global_ServerCfg.SrvID, 90, 10)

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
	tcpMgr, err := StartTcpManager(tcpAddr.IP.String(), tcpAddr.Port, srvMgr, 180)
	if err != nil {
		log.LogErr("addr:", tcpAddr, " err:", err)
		return
	}
	log.LogDebug(tcpMgr)

	srvMgr.tcpMgr = tcpMgr
	srvMgr.httpSrv = httpSrv

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

	AddCmd("dump", "dump", func(argStr []string){
		str := srvMgr.Dump()
		fmt.Println(str)
	})
	commandLineInit()

	ParseCmd()
	tcpMgr.TcpMgrWait()
}

