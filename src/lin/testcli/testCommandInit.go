package main

import (
	"fmt"
	. "lin/msgpacket"
	"runtime"
	"strconv"
)

// "github.com/golang/protobuf/proto"
// "google.golang.org/protobuf"

func CommandTest(argStr []string) {
	msg := &MSG_TEST{}
	msg.Id = 123
	msg.Str = "666"
	globalTcpInfo.TcpSend(MSG_TYPE__MSG_TEST, msg)
}

func CommandLogin(argStr []string) {
	msg := &MSG_LOGIN{}
	msg.Id = 123
	globalTcpInfo.TcpSend(MSG_TYPE__MSG_LOGIN, msg)
}

func CommandMultTest(argStr []string) {
	cor := 1
	count := 1
	if len(argStr) >= 1 {
		cor, _ = strconv.Atoi(argStr[0])
	}
	if len(argStr) >= 2 {
		count, _ = strconv.Atoi(argStr[1])
	}
	for i := 0; i < cor; i ++ {
		go func() {
			for j := 0; j < count; j ++ {
				for _, val := range Global_cliMgr.mapClient {
					msg := &MSG_TEST{}
					msg.Id = val.id
					msg.Str = fmt.Sprintf("%v_%v_%v", val.id, i, j)
					val.TcpSend(MSG_TYPE__MSG_TEST, msg)
				}
				runtime.Gosched()
			}
		}()
	}
}

func CommandMultLogin(argStr []string) {
	count := 1
	if len(argStr) >= 1 {
		count, _ = strconv.Atoi(argStr[0])
	}
	idbase := 0
	if len(argStr) >= 2 {
		idbase, _ = strconv.Atoi(argStr[1])
	}

	for i := 0; i < count; i ++ {
		StartClient(int64(idbase + i), "192.168.2.129:2003")
	}
}

func AddAllCmd(){
	InitCmd()
	AddCmd("test", "test",CommandTest)
	AddCmd("login", "login",CommandLogin)
	AddCmd("mtest", "mult test",CommandMultTest)
	AddCmd("mlogin", "loginMult",CommandMultLogin)
}