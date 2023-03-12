package main

import (
	"fmt"
	. "goserver/msgpacket"
	"strconv"
	"sync"
	"time"
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

	MultiLogin(count, idbase)
}
func MultiLogin(count int, idbase int) {
	for i := 0; i < count; i ++ {
		StartClient(int64(idbase + i), Global_testCfg.addr)
	}
}

func CommandStatic(argStr []string) {
	id := 123
	if len(argStr) >= 1 {
		id,_ = strconv.Atoi(argStr[0])
	}

	c := Global_cliMgr.ClientMgrGet(int64(id))
	if c == nil {
		return
	}

	msg := &MSG_TCP_STATIC{}
	msg.Seq = c.GetSeq()
	c.TcpSend(msg)
}

func CommandLoopTest(argStr []string) {
	count := 900000000
	if len(argStr) >= 1 {
		count, _ = strconv.Atoi(argStr[0])
	}
	smallcount := 1
	if len(argStr) >= 2 {
		smallcount, _ = strconv.Atoi(argStr[1])
	}
	for _, val := range Global_cliMgr.mapClient {
		val.TcpSendLoop(count, smallcount)
	}
}



func CommandDump(argStr []string) {
	Global_cliMgr.total = 0
	var id int64 = 0
	if len(argStr) >= 1 {
		id, _ = strconv.ParseInt(argStr[0], 10, 64)
	}
	c := Global_cliMgr.ClientMgrGet(id)
	if c != nil {
		fmt.Println(c.ClientDump())
	} else {
		for _, val := range Global_cliMgr.mapClient {
			fmt.Println(val.ClientDump())
			Global_cliMgr.total += val.testCountTotal
		}
	}

	totalDiff := Global_cliMgr.total - Global_cliMgr.totalLast
	tnow := float64(time.Now().UnixMilli())
	tdiff := (tnow - Global_cliMgr.timestamp)/float64(1000)
	aver := float64(totalDiff) / tdiff
	fmt.Println(" client count:", len(Global_cliMgr.mapClient), " total:", Global_cliMgr.total, " last:", Global_cliMgr.totalLast,
		" totalDiff:", totalDiff, " tdiff:", tdiff, " aver:", aver)
	Global_cliMgr.timestamp = tnow
	Global_cliMgr.totalLast = Global_cliMgr.total
}

func AddAllCmd(){
	InitCmd()
	AddCmd("test", "test", CommandTest)
	AddCmd("login", "login", CommandLogin)
	AddCmd("mt", "mult test", CommandMultTest)
	AddCmd("mlogin", "loginMult", CommandMultLogin)
	AddCmd("static", "static", CommandStatic)
	AddCmd("lt", "loop test", CommandLoopTest)
	AddCmd("dump", "dump id", CommandDump)
}