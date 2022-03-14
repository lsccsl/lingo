package main

import (
	"lin/lin_common"
	"lin/msgpacket"
	"sync"
)

var Global_wg sync.WaitGroup
var Global_cliMgr *ClientMgr = &ClientMgr{
	mapClient :make(MAP_CLIENT),
}

type TestCfg struct {
	addr string
}
var Global_testCfg = &TestCfg {
	//addr : "192.168.2.129:2003",
	addr : "10.0.14.48:2001",
}

func main() {
	lin_common.ProfileInit()
	AddAllCmd()
	msgpacket.InitMsgParseVirtualTable()

	StartClient(123, Global_testCfg.addr)

	ParseCmd()
	Global_wg.Wait()
}

// todo : del go chan put, send msg by msg instance name