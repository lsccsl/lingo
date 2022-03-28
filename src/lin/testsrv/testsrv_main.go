package main

import (
	"lin/lin_common"
	"lin/msgpacket"
	"sync"
)

var Global_TestSrvMgr = &TestSrvMgr{
	mapSrv : make(MAP_TEST_SERVER),
}
type TestCfg struct {
	ip string
	port int
	httpAddr string

	local_ip string
	local_port_start int
}
var Global_testCfg = &TestCfg {
	ip : "192.168.2.129",port : 2003,
	//ip : "10.0.14.48",	port : 2002,

	httpAddr : "http://192.168.2.129:8803/addserver",
	//httpAddr : "http://10.0.14.48:8802/addserver",

	local_ip : "10.0.14.48",
	local_port_start : 3000,
}
var Global_wg sync.WaitGroup
func main() {
	lin_common.InitLog("./testsrv.log", true)
/*	d := net.Dialer{Timeout: time.Second * time.Duration(30)}
	ctx, canelfun := context.WithCancel(context.Background())
	go func() {
		_, err := d.DialContext(ctx, "tcp", "192.168.2.129:2005")
		lin_common.LogDebug("err string:", err.Error())
		switch t:=err.(type) {
		case *net.OpError:
			switch t1 := t.Err.(type) {
			case *os.SyscallError:
				lin_common.LogDebug(t1)
			default:
				tyerr := reflect.TypeOf(t.Err)
				lin_common.LogDebug(t1, " type kind:", tyerr.Kind(),
					" PkgPath:", tyerr.PkgPath(), " name:", tyerr.Name(), " string:", tyerr.String())
			}
		default:
			lin_common.LogDebug(t)
		}
	}()

	fmt.Println(canelfun)
	canelfun()*/


	lin_common.ProfileInit()

	commandLineInit()

	msgpacket.InitMsgParseVirtualTable()

	lin_common.ParseCmd()

	Global_wg.Wait()
}

// todo 会卡