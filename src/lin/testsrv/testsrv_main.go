package main

import (
	"lin/lin_common"
	"lin/msgpacket"
	"sync"
)

var Global_TestSrvMgr = &TestSrvMgr{
	mapSrv : make(MAP_TEST_SERVER),
}

var Global_wg sync.WaitGroup
func main() {
	commandLineInit()
	lin_common.InitLog("./testsrv.log", true)
	msgpacket.InitMsgParseVirtualTable()
	srv := ConstructTestSrv("10.0.14.48:2001", "192.168.2.129:2003", 1)
	Global_TestSrvMgr.TestSrvMgrAdd(srv)

	lin_common.ParseCmd()

	Global_wg.Wait()
}

