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
	var id string
	flag.StringVar(&pathCfg, "cfg", "../cfg/srvcfg.yml", "config path")
	flag.StringVar(&id, "id", "123", "que srv id")
	flag.Parse()
	server_common.ReadCfg(pathCfg)

	mqMgr := ConstructMsgQueSrv("117.78.3.242:10000", "0.0.0.0:11000", "192.168.15.149:11000",10)
	mqMgr.MsgQueSrvWait()
}
