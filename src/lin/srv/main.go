package main

import (
	"flag"
	"fmt"
	"lin/lin_common"
	_ "lin/msgpacket"
	"net"
	"net/http"
	"os"
)

var srvMgr *ServerMgr

// --path="../cfg/cfg.yml" --id=1
func main() {
	lin_common.ProfileInit()
	fmt.Println(os.Args)

	var pathCfg string
	var id string
	flag.StringVar(&pathCfg, "path", "cfg.yml", "config path")
	flag.StringVar(&id, "id", "123", "server id")
	flag.Parse()
	ReadCfg(pathCfg)
	srvCfg := GetSrvCfgByID(id)

	srvMgr = ConstructServerMgr(srvCfg.SrvID, 30, 10)

	httpAddr, err := net.ResolveTCPAddr("tcp", srvCfg.HttpAddr)
	if err != nil {
		lin_common.LogErr(err)
		return
	}
	httpSrv, err := StartHttpSrvMgr(httpAddr.IP.String(), httpAddr.Port)
	if err != nil {
		lin_common.LogErr(err)
	}

	httpSrv.HttpSrvAddCallback("/test", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, request.URL.Path, " ", request.Form)
	})

	tcpAddr, err := net.ResolveTCPAddr("tcp", srvCfg.BindAddr)
	if err != nil {
		lin_common.LogErr(err)
		return
	}
	tcpMgr, err := StartTcpManager(tcpAddr.IP.String(), tcpAddr.Port, srvMgr, 600)
	if err != nil {
		lin_common.LogErr("addr:", tcpAddr, " err:", err)
		return
	}
	lin_common.LogDebug(tcpMgr)

	srvMgr.tcpMgr = tcpMgr
	srvMgr.httpSrv = httpSrv

	if len(Global_ServerCfg.MapServer) > 1 {
		for _, val := range Global_ServerCfg.MapServer {
			if val.Cluster != srvCfg.Cluster {
				continue
			}
			if val.BindAddr == srvCfg.BindAddr {
				continue
			}
			dialAddr, err := net.ResolveTCPAddr("tcp", val.BindAddr)
			if err != nil {
				lin_common.LogErr(err)
				return
			}
			tcpMgr.TcpDialMgrDial(val.SrvID, dialAddr.IP.String(), dialAddr.Port, 180, 15, true, 10)
			lin_common.LogDebug(val)
		}
	}

	AddCmd("dump", "dump", func(argStr []string){
		bDetail := false
		if len(argStr) > 0 {
			bDetail = true
		}
		str := srvMgr.Dump(bDetail)
		lin_common.LogDebug(str)
	})
	commandLineInit()

	ParseCmd()
	tcpMgr.TcpMgrWait()
}

// aoi path finding
