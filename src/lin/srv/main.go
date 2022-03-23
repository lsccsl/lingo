package main

import (
	"flag"
	"fmt"
	"lin/lin_common"
	"lin/msgpacket"
	_ "lin/msgpacket"
	"lin/tcp"
	"net"
	"net/http"
	"os"
)

var srvMgr *ServerMgr

// --path="../cfg/cfg.yml" --id=1
func main() {
	fmt.Println(os.Args)

	var pathCfg string
	var id string
	flag.StringVar(&pathCfg, "path", "cfg.yml", "config path")
	flag.StringVar(&id, "id", "123", "server id")
	flag.Parse()
	ReadCfg(pathCfg)
	srvCfg := GetSrvCfgByID(id)
	if srvCfg == nil {
		fmt.Println("read cfg err", pathCfg)
	}

	lin_common.InitLog("./srv.log", srvCfg.LogEnableConsolePrint)
	msgpacket.InitMsgParseVirtualTable()
	lin_common.ProfileInit()

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
	httpSrv.HttpSrvAddCallback("/cmd", func(writer http.ResponseWriter, request *http.Request) {
		cmd , _ := request.Form["cmd"]
		if cmd != nil {
			fmt.Fprint(writer, DoCmd(cmd, len(cmd)))
		}
	})

	tcpAddr, err := net.ResolveTCPAddr("tcp", srvCfg.BindAddr)
	if err != nil {
		lin_common.LogErr(err)
		return
	}
	tcpMgr, err := tcp.StartTcpManager(tcpAddr.IP.String(), tcpAddr.Port, srvMgr, 1200)
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

	AddCmd("dump", "dump", func(argStr []string)string{
		bDetail := false
		if len(argStr) > 0 {
			bDetail = true
		}
		str := srvMgr.Dump(bDetail)
		lin_common.LogDebug(str)
		return str
	})
	commandLineInit()

	ParseCmd()
	tcpMgr.TcpMgrWait()
}

// todo aoi path finding, server tcp connection close process
