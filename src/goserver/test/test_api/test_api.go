package main

import "goserver/msgpacket"

func main()  {
	test_uuid()
	return
	msgpacket.InitMsgParseVirtualTable("")
	test_go_routine()
	test_map()
	//test_bmp()
	//test_reflect()
	//test_yaml()
}

