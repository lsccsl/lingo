package main

import (
	"lin/lin_common"
	"lin/msgpacket"
	"sync"
)

var Global_TestSrvMgr = &TestSrvMgr{
	mapSrv : make(MAP_TEST_SERVER),
}
type TestCfg struct {
	addr string
}
var Global_testCfg = &TestCfg {
	addr : "192.168.2.129:2003",
	//addr : "10.0.14.48:2001",
}
var Global_wg sync.WaitGroup
func main() {
	testhttp()
	commandLineInit()
	lin_common.InitLog("./testsrv.log", true)
	msgpacket.InitMsgParseVirtualTable()
	ConstructTestSrv("10.0.14.48:2001", Global_testCfg.addr, 1)

	lin_common.ParseCmd()

	Global_wg.Wait()
}

