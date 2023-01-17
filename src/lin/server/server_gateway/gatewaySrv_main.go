package main

import (
	"lin/lin_common"
	"lin/msgpacket"
	"lin/server/server_msg_que_client"
)

func main() {
	lin_common.InitLog("./srv.log", "./srv_err.log", true, true)

	msgpacket.InitMsgParseVirtualTable("../cfg")

	mqCli := server_msg_que_client.ConstructMgrQueClient("117.78.3.242:10000", 123)

	mqCli.Wait()
}
