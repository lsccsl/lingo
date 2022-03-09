package main

import (
	"flag"
	"fmt"
	"lin/log"
	"lin/msgpacket"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"time"
)

var srvMgr *ServerMgr

func goMemstatic() {
	for {
		//go vet ./…
		//GODEBUG=gctrace=1 ./xxx
		//GODEBUG=gctrace=1 ./xxx 2> gc.log
		runtime.GC()
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		fmt.Printf("Alloc:%d(bytes) HeapIdle:%d(bytes) HeapReleased:%d(bytes)\r\n", ms.Alloc, ms.HeapIdle, ms.HeapReleased)

		time.Sleep(time.Second * time.Duration(10))
	}
}

// --path="../cfg/cfg.yml" --id=1
func main() {
	runtime.SetMutexProfileFraction(1)
	runtime.SetBlockProfileRate(1)

	go func() {
		http.ListenAndServe("0.0.0.0:6060", nil) //go tool pprof --text http://127.0.0.1:6060/debug/pprof/heap
	}()
	go goMemstatic()

	fmt.Println(os.Args)

	var pathCfg string
	var id string
	flag.StringVar(&pathCfg, "path", "cfg.yml", "config path")
	flag.StringVar(&id, "id", "123", "server id")

	flag.Parse()

	ReadCfg(pathCfg)

	srvCfg := GetSrvCfgByID(id)

	msgpacket.InitMsgParseVirtualTable()
	srvMgr = ConstructServerMgr(srvCfg.SrvID, 90, 10)

	httpAddr, err := net.ResolveTCPAddr("tcp", srvCfg.HttpAddr)
	if err != nil {
		log.LogErr(err)
		return
	}
	httpSrv, err := StartHttpSrvMgr(httpAddr.IP.String(), httpAddr.Port)
	if err != nil {
		log.LogErr(err)
	}

	httpSrv.HttpSrvAddCallback("/test", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, request.URL.Path, " ", request.Form)
	})

	tcpAddr, err := net.ResolveTCPAddr("tcp", srvCfg.BindAddr)
	if err != nil {
		log.LogErr(err)
		return
	}
	tcpMgr, err := StartTcpManager(tcpAddr.IP.String(), tcpAddr.Port, srvMgr, 180)
	if err != nil {
		log.LogErr("addr:", tcpAddr, " err:", err)
		return
	}
	log.LogDebug(tcpMgr)

	srvMgr.tcpMgr = tcpMgr
	srvMgr.httpSrv = httpSrv

	if len(Global_ServerCfg.MapServer) > 1 {
		for _, val := range Global_ServerCfg.MapServer {
			if val.Cluster != srvCfg.Cluster {
				continue
			}
			if val.BindAddr == srvCfg.BindAddr {
				continue
			}
			dialAddr, err := net.ResolveTCPAddr("tcp", val.BindAddr)
			if err != nil {
				log.LogErr(err)
				return
			}
			tcpMgr.TcpDialMgrDial(val.SrvID, dialAddr.IP.String(), dialAddr.Port, 180, 15, true, 10)
			log.LogDebug(val)
		}
	}

	AddCmd("dump", "dump", func(argStr []string){
		str := srvMgr.Dump()
		fmt.Println(str)
	})
	commandLineInit()

	ParseCmd()
	tcpMgr.TcpMgrWait()
}

