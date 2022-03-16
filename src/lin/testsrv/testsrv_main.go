package main

import "sync"

var Global_wg sync.WaitGroup
func main() {
	ConstructTestSrv("10.0.14.48:2002", "192.168.2.129:2003", 2)

	Global_wg.Wait()
}

