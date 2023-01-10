package main

import (
	"lin/lin_common"
	"lin/msgpacket"
)

func main() {
	lin_common.InitLog("./srv.log", "./srv_err.log", true, true)

	msgpacket.InitMsgParseVirtualTable("../cfg")

	mqMgr := ConstructMsgqueMgr("0.0.0.0:11000", 10)
	mqMgr.MsgQueCenterMgrWait()
}
