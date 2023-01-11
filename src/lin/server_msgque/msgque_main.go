package main

import (
	"lin/lin_common"
	"lin/msgpacket"
)

func main() {
	lin_common.InitLog("./srv.log", "./srv_err.log", true, true)

	msgpacket.InitMsgParseVirtualTable("../cfg")

	mqMgr := ConstructMsgQueSrv("117.78.3.242:10000", "0.0.0.0:11000", "192.168.15.149:11000",10)
	mqMgr.MsgQueSrvWait()
}
