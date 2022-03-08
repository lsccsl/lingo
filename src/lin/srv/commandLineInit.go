package main

import (
	"lin/log"
	"lin/msgpacket"
	"strconv"
)

func testrpc(argStr []string){
	var srvID int64 = 1
	if len(argStr) >= 1 {
		srvID, _ = strconv.ParseInt(argStr[0], 10, 64)
	}

	log.LogDebug(srvID)
	msg := srvMgr.SendRPC_Async(srvID, msgpacket.MSG_TYPE__MSG_TEST, &msgpacket.MSG_TEST{Id:567}, 10 * 1000)

	log.LogDebug(msg)
}

func commandLineInit(){
	AddCmd("help", "help", DumpAllCmd)
	AddCmd("testrpc", "testrpc srvID", testrpc)
}
