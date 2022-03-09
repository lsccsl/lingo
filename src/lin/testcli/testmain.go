package main

import (
	"lin/msgpacket"
	"sync"
)

var Global_wg sync.WaitGroup
var Global_cliMgr *ClientMgr = &ClientMgr{
	mapClient :make(MAP_CLIENT),
}

func main() {
	AddAllCmd()
	msgpacket.InitMsgParseVirtualTable()

	StartClient(123, "10.0.14.48:2001")
	//StartClient(123,"192.168.2.129:2003")

	ParseCmd()
	Global_wg.Wait()
}

