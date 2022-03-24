package main

import (
	"lin/lin_common"
	"lin/msgpacket"
	"strconv"
)

func testrpc(argStr []string)string{
	var srvID int64 = 1
	if len(argStr) >= 1 {
		srvID, _ = strconv.ParseInt(argStr[0], 10, 64)
	}

	lin_common.LogDebug(srvID)
	msg := srvMgr.SendRPC_Async(srvID, msgpacket.MSG_TYPE__MSG_TEST, &msgpacket.MSG_TEST{Id:567}, 10 * 1000)

	lin_common.LogDebug(msg)
	return ""
}

func commandLineInit(){
	lin_common.AddCmd("help", "help", lin_common.DumpAllCmd)
	lin_common.AddCmd("testrpc", "testrpc srvID", testrpc)
}
