package main

import (
	"fmt"
	"lin/lin_common"
	"strconv"
	"time"
)

func MultiSrv(count int, idbase int) {
	for i := 0; i < count; i ++ {
		ConstructTestSrv("", Global_testCfg.addr, int64(idbase + i))
	}
}
func CommandNewSrv(argStr []string) string {
	count := 10
	if len(argStr) >= 1 {
		count, _ = strconv.Atoi(argStr[0])
	}
	idbase := 10
	if len(argStr) >= 2 {
		idbase, _ = strconv.Atoi(argStr[1])
	}

	MultiSrv(count, idbase)
	return ""
}

func CommandDump(argStr []string) string {
	Global_TestSrvMgr.total = 0
	for _, val := range Global_TestSrvMgr.mapSrv {
		Global_TestSrvMgr.total += val.totalRpcDial
	}

	totalDiff := Global_TestSrvMgr.total - Global_TestSrvMgr.totalLast
	tnow := float64(time.Now().UnixMilli())
	tdiff := (tnow - Global_TestSrvMgr.timestamp)/float64(1000)
	aver := float64(totalDiff) / tdiff
	fmt.Println(" client count:", len(Global_TestSrvMgr.mapSrv), " total:", Global_TestSrvMgr.total, " last:", Global_TestSrvMgr.totalLast,
		" totalDiff:", totalDiff, " tdiff:", tdiff, " aver:", aver)
	Global_TestSrvMgr.timestamp = tnow
	Global_TestSrvMgr.totalLast = Global_TestSrvMgr.total

	return ""
}
func commandLineInit(){
	lin_common.AddCmd("dump", "dump id",CommandDump)
	lin_common.AddCmd("help", "help", lin_common.DumpAllCmd)
	lin_common.AddCmd("ns", "new server", CommandNewSrv)
}
