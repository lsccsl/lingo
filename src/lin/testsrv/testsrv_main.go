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
	ip string
	port int
	httpAddr string

	local_ip string
	local_port_start int
}
var Global_testCfg = &TestCfg {
	ip : "192.168.2.129",port : 2003,
	//ip : "10.0.14.48",	port : 2002,

	httpAddr : "http://192.168.2.129:8803/addserver",
	//httpAddr : "http://10.0.14.48:8802/addserver",

	local_ip : "10.0.14.48",
	local_port_start : 3000,
}
var Global_wg sync.WaitGroup
func main() {
	//testhttp()
	commandLineInit()
	lin_common.InitLog("./testsrv.log", true)
	msgpacket.InitMsgParseVirtualTable()

	//ConstructTestSrv("10.0.14.48:2001", Global_testCfg.ip + ":" + strconv.Itoa(Global_testCfg.port), 1)

	lin_common.ParseCmd()

	Global_wg.Wait()
}

