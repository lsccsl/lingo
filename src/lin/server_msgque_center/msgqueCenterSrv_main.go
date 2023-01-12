package main

import (
	"flag"
	"lin/lin_common"
	"lin/msgpacket"
	"lin/server_common"
)

func main() {
	lin_common.InitLog("./srv.log", "./srv_err.log", true, true)

	msgpacket.InitMsgParseVirtualTable("../cfg")

	var pathCfg string
	flag.StringVar(&pathCfg, "cfg", "../cfg/srvcfg.yml", "config path")
	flag.Parse()
	server_common.ReadCfg(pathCfg)

	mqMgr := ConstructMsgQueCenterSrv(server_common.Global_ServerCfg.MsgQueCent.BindAddr, 10)
	mqMgr.Wait()
}