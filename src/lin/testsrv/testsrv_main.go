package main

import (
	"lin/msgpacket"
	"sync"
)

var Global_wg sync.WaitGroup
func main() {
	msgpacket.InitMsgParseVirtualTable()
	ConstructTestSrv("10.0.14.48:2001", "192.168.2.129:2003", 1)
	Global_wg.Wait()
}

