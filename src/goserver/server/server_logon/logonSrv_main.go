package main

import (
	"flag"
	"goserver/common"
	"goserver/msgpacket"
	"goserver/server/server_common"
	"strconv"
)

func main()  {
	common.InitLog("./srv.log", "./srv_err.log", true, true, false)

	msgpacket.InitMsgParseVirtualTable("../cfg")

	var pathCfg string
	var id string
	flag.StringVar(&pathCfg, "cfg", "../cfg/srvcfg.yml", "config path")
	flag.StringVar(&id, "id", "1", "que srv id")
	flag.Parse()
	server_common.ReadCfg(pathCfg)

	gs := ConstructLogonSrv(id)

	common.AddCmd("dump", "dump", func(argStr []string)string{
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
		str := gs.mqClient.Dump(bDetail)
		if bLog {
			common.LogDebug(str)
		}
		return str
	})
	common.ParseCmd()

	gs.Wait()
}
