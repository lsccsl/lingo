package main

import "lin/msgpacket"

func main()  {
	msgpacket.InitMsgParseVirtualTable("")
	test_go_routine()
	test_map()
	//test_bmp()
	//test_reflect()
	//test_yaml()
}

