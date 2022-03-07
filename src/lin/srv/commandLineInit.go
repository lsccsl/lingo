package main

import (
	"lin/msgpacket"
	"strconv"
)

func testrpc(argStr []string){
	var srvID int64 = 1
	if len(argStr) >= 1 {
		srvID, _ = strconv.ParseInt(argStr[0], 10, 64)
	}

	srvMgr.SendRPC_Async(srvID, msgpacket.MSG_TYPE__MSG_TEST, &msgpacket.MSG_TEST{Id:567}, 10)
}

func commandLineInit(){
	AddCmd("help", "help", DumpAllCmd)
	AddCmd("testrpc", "testrpc srvID", testrpc)
}
