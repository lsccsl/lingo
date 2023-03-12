package main

import "C"
import (
	"encoding/json"
	"flag"
	"fmt"
	"goserver/common"
	"goserver/msgpacket"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"
)

type ServerFromHttp struct {
	SrvID int64
	IP string
	Port int
}

const TCP_READ_CLOSE_EXPIRE = 600

func main() {
	rand.Seed(time.Now().Unix())
	fd := common.FD_DEF{}
	fmt.Println("fd:", fd.String())

	// read config
	var pathCfg string
	var id string
	flag.StringVar(&pathCfg, "path", "cfg.yml", "config path")
	flag.StringVar(&id, "id", "123", "server id")
	flag.Parse()
	ReadCfg(pathCfg)
	srvCfg := GetSrvCfgByID(id)
	if srvCfg == nil {
		fmt.Println("read cfg err", pathCfg)
	}

	// log and profile
	common.InitLog("./epollsrv.log", "./epollsrv_err.log", srvCfg.LogEnableConsolePrint, true, false)
	common.ProfileInit(true, 6060)

	common.LogErr("test err log")

	fmt.Println("begin epoll listen, ip:", srvCfg.BindAddr)

	msgpacket.InitMsgParseVirtualTable(Global_ServerCfg.Msgdef)

	// epoll mgr
	eSrvMgr, err := ConstructorServerMgr(srvCfg.BindAddr,
		20, 20, 8,
		600,900,
		true)
	if err != nil {
		fmt.Println(err)
		common.LogDebug(err)
		return
	}

	fmt.Println("end epoll listen")

	// http interface
	httpAddr, err := net.ResolveTCPAddr("tcp", srvCfg.HttpAddr)
	if err != nil {
		common.LogErr(err)
		return
	}
	httpSrv, err := common.StartHttpSrvMgr(httpAddr.IP.String(), httpAddr.Port)
	if err != nil {
		common.LogErr(err)
	}
	//http://192.168.0.104:8803/cmd?cmd=dump
	httpSrv.HttpSrvAddCallback("/cmd", func(writer http.ResponseWriter, request *http.Request) {
		cmd , _ := request.Form["cmd"]
		fmt.Println(cmd)
		if cmd != nil {
			fmt.Fprint(writer, common.DoCmd(cmd, len(cmd)))
		}
	}, "")
	httpSrv.HttpSrvAddCallback("/addserver", func(writer http.ResponseWriter, request *http.Request) {
		bin := make([]byte, request.ContentLength, request.ContentLength)
		request.Body.Read(bin)
/*		lin_common.LogDebug("~~~~~recv addserver", string(bin), " len:", request.ContentLength,
			" body len:", len(bin), " post form:", request.PostForm)*/
		sh := &ServerFromHttp{}
		json.Unmarshal(bin, sh)
		common.LogDebug("add srv:", sh.SrvID, " addr:", sh.IP, ":", sh.Port)
		eSrvMgr.AddRemoteSrv(sh.SrvID, sh.IP + ":" + strconv.Itoa(sh.Port), TCP_READ_CLOSE_EXPIRE)
		writer.Write(bin)
	}, "")

	// command line
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
		str := eSrvMgr.Dump(bDetail)
		if bLog {
			common.LogDebug(str)
		}
		return str
	})
	common.ParseCmd()

	eSrvMgr.lsn.EPollListenerWait()
}