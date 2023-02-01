package main

import (
	"goserver/common"
	"goserver/msgpacket"
	"strconv"
)

func testrpc(argStr []string)string{
	var srvID int64 = 1
	if len(argStr) >= 1 {
		srvID, _ = strconv.ParseInt(argStr[0], 10, 64)
	}

	common.LogDebug(srvID)
	msg, _ := srvMgr.SendRPC_Async(srvID, msgpacket.MSG_TYPE__MSG_TEST, &msgpacket.MSG_TEST{Id: 567}, 10 * 1000)

	common.LogDebug(msg)
	return ""
}

func commandLineInit(){
	common.AddCmd("help", "help", common.DumpAllCmd)
	common.AddCmd("testrpc", "testrpc srvID", testrpc)
}
