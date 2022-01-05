package main

import (
	"fmt"
	"net"
)

type GlobalNetInfo struct{
	tcpListerner_ net.Listener
	mapClient_ map[int64]tcpClientInterface
}
var GNetInfo GlobalNetInfo

type tcpClientInterface interface {
	MethodDummpy()
}

type tcpClient struct{
	tcpCon_ net.Conn
}
func (*tcpClient)MethodDummpy(){
}


func go_tcpAccept(){
	for {
		fmt.Println("srv ==== begin accpet")
		var client tcpClient
		var err error
		client.tcpCon_, err = GNetInfo.tcpListerner_.Accept()
		fmt.Println(client.tcpCon_, "srv ==== new con err:", err)

		go client.go_tcpRead()
	}
}

func main(){
	var err error
	GNetInfo.tcpListerner_, err = net.Listen("tcp", "127.0.0.1:6666")
	if nil != err{
		fmt.Println(err)
	}
	go go_tcpAccept()
}

func (pThis*tcpClient)go_tcpRead(){

}
