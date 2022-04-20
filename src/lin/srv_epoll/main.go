package main

import (
	"flag"
	"fmt"
	"lin/lin_common"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	fd := lin_common.FD_DEF{}
	fmt.Println("fd:", fd.String())

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
	lin_common.InitLog("./epollsrv.log", /*srvCfg.LogEnableConsolePrint*/true, true)
	lin_common.ProfileInit(true, 6060)

	eSrvMgr, err := ConstructorEpollServerMgr("192.168.2.129:2003", 10)
	lin_common.LogDebug(err)
	eSrvMgr.lsn.EPollListenerWait()
}
