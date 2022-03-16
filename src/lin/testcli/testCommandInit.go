package main

import (
	"fmt"
	. "lin/msgpacket"
	"strconv"
	"sync"
)

// "github.com/golang/protobuf/proto"
// "google.golang.org/protobuf"

func CommandTest(argStr []string) {
	msg := &MSG_TEST{}
	msg.Id = 123
	msg.Str = "666"
	globalTcpInfo.TcpSend(msg)
}

func CommandLogin(argStr []string) {
	msg := &MSG_LOGIN{}
	msg.Id = 123
	globalTcpInfo.TcpSend(msg)
}

func CommandMultTest(argStr []string) {
	cor := 1
	count := 1000000
	if len(argStr) >= 1 {
		cor, _ = strconv.Atoi(argStr[0])
	}
	if len(argStr) >= 2 {
		count, _ = strconv.Atoi(argStr[1])
	}

	go func() {
		fmt.Println("cor:", cor, " count:", count)
		var wgtmp sync.WaitGroup
		for icor := 0; icor < cor; icor ++ {
			wgtmp.Add(1)
			cori := icor
			go func() {
				for j := 0; j < count; j ++ {
					for _, val := range Global_cliMgr.mapClient {
						msg := &MSG_TEST{}
						msg.Id = val.id
						msg.Str = fmt.Sprintf("%v_%v_%v", val.id, cori, j)
						val.TcpSend(msg)
/*						if (j % 10000 == 0) {
							runtime.Gosched()
						}
*/						if (j % 1000 == 0) {
							fmt.Println("test:", cori, j)
						}
					}
				}
				wgtmp.Done()
				fmt.Println("coroutine done", cori)
			}()
		}
		wgtmp.Wait()
		fmt.Println("\r\n\r\n all coroutine done")
		for _, val := range Global_cliMgr.mapClient {
			msg := &MSG_TCP_STATIC{}
			val.TcpSend(msg)
		}
	}()
}

func CommandMultLogin(argStr []string) {
	count := 2000
	if len(argStr) >= 1 {
		count, _ = strconv.Atoi(argStr[0])
	}
	idbase := 200
	if len(argStr) >= 2 {
		idbase, _ = strconv.Atoi(argStr[1])
	}

	for i := 0; i < count; i ++ {
		StartClient(int64(idbase + i),Global_testCfg.addr)
	}
}

func CommandStatic(argStr []string) {
	id := 123
	if len(argStr) >= 1 {
		id,_ = strconv.Atoi(argStr[0])
	}

	c := Global_cliMgr.ClientMgrAddGet(int64(id))
	if c == nil {
		return
	}

	msg := &MSG_TCP_STATIC{}
	msg.Seq = c.GetSeq()
	c.TcpSend(msg)
}

func CommandLoopTest(argStr []string) {
	count := 1
	if len(argStr) >= 1 {
		count, _ = strconv.Atoi(argStr[0])
	}
	for _, val := range Global_cliMgr.mapClient {
		val.TcpSendLoop(count)
	}
}

func AddAllCmd(){
	InitCmd()
	AddCmd("test", "test",CommandTest)
	AddCmd("login", "login",CommandLogin)
	AddCmd("mt", "mult test",CommandMultTest)
	AddCmd("mlogin", "loginMult",CommandMultLogin)
	AddCmd("static", "static",CommandStatic)
	AddCmd("lt", "loop test",CommandLoopTest)
}
