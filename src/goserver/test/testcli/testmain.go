package main

import (
	"goserver/common"
	"goserver/msgpacket"
	"os"
	"strconv"
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
	//addr : "192.168.2.129:3001",
	addr : "192.168.2.129:2003",
	//addr : "10.0.14.48:2001",
}

// testcli 500 100 | lt
func main() {
	common.InitLog("./testcli.log", true, false, false)
	//lin_common.ProfileInit()
	AddAllCmd()
	msgpacket.InitMsgParseVirtualTable()

	if len(os.Args) >= 3 {
		count, _ := strconv.Atoi(os.Args[1])
		idbase, _ := strconv.Atoi(os.Args[2])
		common.LogDebug("auto login:", count, " idbase:", idbase)
		MultiLogin(count, idbase)
	}

	ParseCmd()
	Global_wg.Wait()
}

// todo : del go chan put, set tcp read time out