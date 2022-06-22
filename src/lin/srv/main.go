package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"lin/lin_common"
	"lin/msgpacket"
	_ "lin/msgpacket"
	"lin/tcp"
	"net"
	"net/http"
	"os"
	"strconv"
)

var srvMgr *ServerMgr

type ServerFromHttp struct {
	SrvID int64
	IP string
	Port int
}

const TCP_READ_CLOSE_EXPIRE = 600
// --path="../cfg/cfg.yml" --id=1
// nohup ./srv --path="../cfg/cfg.yml" --id=3 >> srv.out 2>&1 &
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

	lin_common.InitLog("./srv.log", "./srv_err.log", srvCfg.LogEnableConsolePrint, true)
	msgpacket.InitMsgParseVirtualTable(Global_ServerCfg.Msgdef)
	lin_common.ProfileInit(true, 6060)

	srvMgr = ConstructServerMgr(srvCfg.SrvID, 30, 100)

	httpAddr, err := net.ResolveTCPAddr("tcp", srvCfg.HttpAddr)
	if err != nil {
		lin_common.LogErr(err)
		return
	}
	httpSrv, err := lin_common.StartHttpSrvMgr(httpAddr.IP.String(), httpAddr.Port)
	if err != nil {
		lin_common.LogErr(err)
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", srvCfg.BindAddr)
	if err != nil {
		lin_common.LogErr(err)
		return
	}
	tcpMgr, err := tcp.ConstructTcpManager(tcpAddr.IP.String(), tcpAddr.Port, srvMgr, TCP_READ_CLOSE_EXPIRE)
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
			srvMgr.AddRemoteServer(val.SrvID, dialAddr.IP.String(), dialAddr.Port, TCP_READ_CLOSE_EXPIRE, 15, true, 3)
			lin_common.LogDebug(val)
		}
	}

	lin_common.AddCmd("dump", "dump", func(argStr []string)string{
		bDetail := false
		bLog := true
		if len(argStr) >= 1 {
			detail, _ := strconv.Atoi(argStr[0])
			bDetail = (detail != 0)
		}
		if len(argStr) >= 2 {
			needLog, _ := strconv.Atoi(argStr[1])
			bLog = (needLog != 0)
		}
		str := srvMgr.Dump(bDetail)
		if bLog {
			lin_common.LogDebug(str)
		}
		return str
	})
	commandLineInit()

	httpSrv.HttpSrvAddCallback("/test", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, request.URL.Path, " ", request.Form)
	}, "test")
	httpSrv.HttpSrvAddCallback("/cmd", func(writer http.ResponseWriter, request *http.Request) {
		cmd , _ := request.Form["cmd"]
		if cmd != nil {
			fmt.Fprint(writer, lin_common.DoCmd(cmd, len(cmd)))
		}
	}, "cmd")
	httpSrv.HttpSrvAddCallback("/addserver", func(writer http.ResponseWriter, request *http.Request) {
		bin := make([]byte, request.ContentLength, request.ContentLength)
		request.Body.Read(bin)
		//lin_common.LogDebug(string(bin), " len:", request.ContentLength)
		sh := &ServerFromHttp{}
		json.Unmarshal(bin, sh)
		lin_common.LogDebug("add srv:", sh.SrvID, " addr:", sh.IP, ":", sh.Port)
		srvMgr.AddRemoteServer(sh.SrvID, sh.IP, sh.Port, TCP_READ_CLOSE_EXPIRE, 15, true, 3)
		writer.Write(bin)
	}, "addserver")
	httpSrv.HttpSrvAddCallback("/delserver", func(writer http.ResponseWriter, request *http.Request) {
		strSrv, _ := request.Form["srv"]
		if len(strSrv) >= 1 {
			srvID, _ := strconv.ParseInt(strSrv[0], 10, 64)
			srvMgr.delServer(srvID)
		}
	}, "delserver")

	lin_common.ParseCmd()
	tcpMgr.TcpMgrWait()
}

// todo aoi path finding, server tcp connection close process
// https://github.com/panjf2000/gnet
// https://github.com/golang/sys
