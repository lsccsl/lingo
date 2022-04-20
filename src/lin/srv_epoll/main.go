package main

import (
	"flag"
	"fmt"
	"lin/lin_common"
)

func main() {
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
	lin_common.InitLog("./srv.log", srvCfg.LogEnableConsolePrint, true)
	lin_common.ProfileInit(true, 6060)
}
