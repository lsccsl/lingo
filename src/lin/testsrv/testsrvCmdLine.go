package main

import (
	"fmt"
	"lin/lin_common"
	"time"
)

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
}
