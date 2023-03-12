package main

import (
	"flag"
	"goserver/common"
	"goserver/msgpacket"
	"goserver/server/server_common"
	"strconv"
)

func main() {
	common.InitLog("./srv.log", "./srv_err.log", true, true, false)

	msgpacket.InitMsgParseVirtualTable("../cfg")

	// read config
	var pathCfg string
	var id string
	flag.StringVar(&pathCfg, "cfg", "../cfg/srvcfg.yml", "config path")
	flag.StringVar(&id, "id", "1", "que srv id")
	flag.Parse()
	server_common.ReadCfg(pathCfg)
	qCfg := server_common.GetMsgQueSrvCfg(id)
	common.LogInfo(qCfg)

	// begin epoll listener
	mqMgr := ConstructMsgQueSrv(server_common.Global_ServerCfg.MsgQueCent.OutAddr,
		qCfg.BindAddr, qCfg.OutAddr,10)

	// cmd line
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
		str := mqMgr.Dump(bDetail)
		if bLog {
			common.LogDebug(str)
		}
		return str
	})
	common.ParseCmd()

	mqMgr.Wait()
}

// todo send_async msg / get srv id by srvtype
// todo client srv reconnect