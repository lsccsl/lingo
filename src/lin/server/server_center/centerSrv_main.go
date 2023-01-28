package main

import (
	"flag"
	"lin/lin_common"
	"lin/msgpacket"
	"lin/server/server_common"
	"strconv"
)

func main()  {
	lin_common.InitLog("./srv.log", "./srv_err.log", true, true)

	msgpacket.InitMsgParseVirtualTable("../cfg")

	var pathCfg string
	flag.StringVar(&pathCfg, "cfg", "../cfg/srvcfg.yml", "config path")
	flag.Parse()
	server_common.ReadCfg(pathCfg)

	cs := ConstructCenterSrv()

	lin_common.AddCmd("dump", "dump", func(argStr []string)string{
		bDetail := false
		bLog := true
		if len(argStr) >= 1 {
			detail, _ := strconv.Atoi(argStr[0])
			bDetail = (detail != 0)
		}
		if len(argStr) >= 2 {
			needLog, _ := strconv.Atoi(argStr[1])
			bLog = (needLog != 0)
		}
		str := cs.Dump(bDetail)
		if bLog {
			lin_common.LogDebug(str)
		}
		return str
	})
	lin_common.ParseCmd()
	cs.Wait()
}

