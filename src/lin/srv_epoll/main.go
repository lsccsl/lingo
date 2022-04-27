package main

import (
	"flag"
	"fmt"
	"lin/lin_common"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	fd := lin_common.FD_DEF{}
	fmt.Println("fd:", fd.String())

	// read config
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

	// log and profile
	lin_common.InitLog("./epollsrv.log", srvCfg.LogEnableConsolePrint, true)
	lin_common.ProfileInit(true, 6060)

	// epoll mgr
	eSrvMgr, err := ConstructorEpollServerMgr(srvCfg.BindAddr/*"192.168.2.129:2003"*/,
		1000, 1000, 10,
		900,900,
		true)
	if err != nil {
		lin_common.LogDebug(err)
		return
	}

	// http interface
	httpAddr, err := net.ResolveTCPAddr("tcp", srvCfg.HttpAddr)
	if err != nil {
		lin_common.LogErr(err)
		return
	}
	httpSrv, err := lin_common.StartHttpSrvMgr(httpAddr.IP.String(), httpAddr.Port)
	if err != nil {
		lin_common.LogErr(err)
	}
	httpSrv.HttpSrvAddCallback("/cmd", func(writer http.ResponseWriter, request *http.Request) {
		cmd , _ := request.Form["cmd"]
		if cmd != nil {
			fmt.Fprint(writer, lin_common.DoCmd(cmd, len(cmd)))
		}
	})

	// command line
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
		str := eSrvMgr.Dump(bDetail)
		if bLog {
			lin_common.LogDebug(str)
		}
		return str
	})
	lin_common.ParseCmd()

	eSrvMgr.lsn.EPollListenerWait()
}
