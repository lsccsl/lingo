package main

import (
	"flag"
	"lin/lin_common"
	"lin/msgpacket"
	"lin/server/server_common"
)

func main() {
	lin_common.InitLog("./srv.log", "./srv_err.log", true, true)

	msgpacket.InitMsgParseVirtualTable("../cfg")

	var pathCfg string
	var id string
	flag.StringVar(&pathCfg, "cfg", "../cfg/srvcfg.yml", "config path")
	flag.StringVar(&id, "id", "1", "que srv id")
	flag.Parse()
	server_common.ReadCfg(pathCfg)

	qCfg := server_common.GetMsgQueSrvCfg(id)
	lin_common.LogInfo(qCfg)

	mqMgr := ConstructMsgQueSrv(server_common.Global_ServerCfg.MsgQueCent.OutAddr,
		qCfg.BindAddr, qCfg.OutAddr,10)
	mqMgr.MsgQueSrvWait()
}
