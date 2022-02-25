package main

import (
	. "lin/msg"
)

// "github.com/golang/protobuf/proto"
// "google.golang.org/protobuf"

func CommandTest(argStr []string) {
	msg := &MSG_TEST{}
	msg.MsgInt = 123
	globalTcpInfo.TcpSend(MSG_TYPE__MSG_TEST, msg)
}

func CommandLogin(argStr []string) {
	msg := &MSG_LOGIN{}
	msg.Id = 123
	globalTcpInfo.TcpSend(MSG_TYPE__MSG_LOGIN, msg)
}

func AddAllCmd(){
	InitCmd()
	AddCmd("test", "test",CommandTest)
	AddCmd("login", "login",CommandLogin)
}