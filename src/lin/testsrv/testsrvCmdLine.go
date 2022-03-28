package main

import (
	"fmt"
	"lin/lin_common"
	"lin/msgpacket"
	"strconv"
	"time"
)

func MultiSrv(count int, idbase int) {
	for i := 0; i < count; i ++ {
		srvid := int64(idbase + i)
		port := Global_testCfg.local_port_start + i
		httpAddDial(&ServerFromHttp{
			SrvID: srvid,
			IP: Global_testCfg.local_ip,
			Port: port,
		})
		ConstructTestSrv(Global_testCfg.local_ip + ":" + strconv.Itoa(port), Global_testCfg.ip + ":" + strconv.Itoa(Global_testCfg.port), int64(idbase + i))
	}
}

func CommandMultiSrv(argStr []string) string {
	count := 1
	if len(argStr) >= 1 {
		count, _ = strconv.Atoi(argStr[0])
	}
	idbase := 100
	if len(argStr) >= 2 {
		idbase, _ = strconv.Atoi(argStr[1])
	}

	MultiSrv(count, idbase)
	return ""
}

func CommandTestRPC(argStr []string) string {
	count := 100000000
	if len(argStr) >= 1 {
		count, _ = strconv.Atoi(argStr[0])
	}
	msg := &msgpacket.MSG_TEST_RPC{
		RpcCount: int64(count),
	}
	for _, val := range Global_TestSrvMgr.mapSrv {
		val.tcpAcpt.Write(msgpacket.ProtoPacketToBin(msgpacket.MSG_TYPE__MSG_TEST_RPC, msg))
	}

	return ""
}

func CommandDump(argStr []string) string {
	Global_TestSrvMgr.total = 0
	Global_TestSrvMgr.totalReqRecv = 0
	for _, val := range Global_TestSrvMgr.mapSrv {
		fmt.Println("dial id:", val.DialConnectionID, " acpt id", val.AcptConnectionID, " total:", val.totalRpcDial, " total write:", val.totalWriteRpc)
		Global_TestSrvMgr.total += val.totalRpcDial
		Global_TestSrvMgr.totalReqRecv += val.totalRpcRecv
	}

	totalDiff := Global_TestSrvMgr.total - Global_TestSrvMgr.totalLast
	totalReqDiff := Global_TestSrvMgr.totalReqRecv - Global_TestSrvMgr.totalReqRecvLast
	tnow := float64(time.Now().UnixMilli())
	tdiff := (tnow - Global_TestSrvMgr.timestamp)/float64(1000)
	aver := float64(totalDiff) / tdiff
	reqAver := float64(totalReqDiff) / tdiff
	fmt.Println(" client count:", len(Global_TestSrvMgr.mapSrv), " total:", Global_TestSrvMgr.total, " last:", Global_TestSrvMgr.totalLast,
		" totalDiff:", totalDiff, " tdiff:", tdiff, " aver:", aver, " req aver:", reqAver)
	Global_TestSrvMgr.timestamp = tnow
	Global_TestSrvMgr.totalLast = Global_TestSrvMgr.total
	Global_TestSrvMgr.totalReqRecvLast = Global_TestSrvMgr.totalReqRecv

	return ""
}


func commandLineInit(){
	lin_common.AddCmd("dump", "dump id",CommandDump)
	lin_common.AddCmd("help", "help", lin_common.DumpAllCmd)
	lin_common.AddCmd("ms", "multi server", CommandMultiSrv)
	lin_common.AddCmd("tr", "test rpc", CommandTestRPC)
}
