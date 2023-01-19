package main

import (
	"lin/lin_common"
	"lin/msgpacket"
	"lin/server/server_msg_que_client"
	"strconv"
)

func main() {
	lin_common.InitLog("./srv.log", "./srv_err.log", true, true)

	msgpacket.InitMsgParseVirtualTable("../cfg")

	mqCli := msg_que_client.ConstructMgrQueClient("117.78.3.242:10000", 123)

	lin_common.AddCmd("dump", "dump", func(argStr []string)string{
		bDetail := false
		bLog := true
		if len(argStr) >= 1 {
			detail, _ := strconv.Atoi(argStr[0])
			bDetail = (detail != 0)
		}
		if len(argStr) >= 2 {
			needLog, _ := strconv.Atoi(argStr[1])
			bLog = (needLog != 0)
		}
		str := mqCli.Dump(bDetail)
		if bLog {
			lin_common.LogDebug(str)
		}
		return str
	})
	lin_common.ParseCmd()
	mqCli.Wait()
}
