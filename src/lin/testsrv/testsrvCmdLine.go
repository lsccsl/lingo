package main

import (
	"fmt"
	"lin/lin_common"
	"lin/msgpacket"
	"math"
	"strconv"
	"time"
)

func MultiSrv(count int, idbase int, port_start int) {
	for i := 0; i < count; i ++ {
		srvid := int64(idbase + i)
		port := port_start + i

		ConstructTestSrv(Global_testCfg.local_ip + ":" + strconv.Itoa(port), port, Global_testCfg.remote_ip + ":" + strconv.Itoa(Global_testCfg.remote_port),
			srvid)
	}

	for i := 0; i < count; i ++ {
		srvid := int64(idbase + i)
		port := port_start + i
		if srvid == 599 {
			lin_common.LogDebug("srv:", srvid, " send http")
		}
		httpAddDial(&ServerFromHttp{
			SrvID: srvid,
			IP: Global_testCfg.local_ip,
			Port: port,
		})
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
	port_start := Global_testCfg.local_port_start
	if len(argStr) >= 3 {
		port_start, _ = strconv.Atoi(argStr[2])
	}

	MultiSrv(count, idbase, port_start)
	return ""
}

func CommandTestRPC(argStr []string) string {
	for _, val := range Global_TestSrvMgr.mapSrv {
		val.TestSrvBeginDial()
	}

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

	if len(argStr) >= 1{
		srvID, _ := strconv.ParseInt(argStr[0], 10, 64)
		val := Global_TestSrvMgr.mapSrv[srvID]
		if val != nil {
			fmt.Println("srv:", val.srvId, "dial id:", val.DialConnectionID, " acpt id", val.AcptConnectionID, " total:", val.totalRpcDial, " total write:", val.totalWriteRpc,
				" redial:", val.totalRedial, " reAcpt:", val.totalReAcpt)
			return ""
		}
	}

	Global_TestSrvMgr.total = 0
	Global_TestSrvMgr.totalReqRecv = 0
	var totalRedial int64 = 0
	var totalReAcpt int64 = 0
	minRTT := int64(math.MaxInt64)
	maxRTT := int64(0)
	totalRTT := int64(0)
	noAcpt := 0
	noDail := 0
	for _, val := range Global_TestSrvMgr.mapSrv {
		fmt.Println("srv:", val.srvId, "dial id:", val.DialConnectionID, " acpt id", val.AcptConnectionID, " total:", val.totalRpcDial, " total write:", val.totalWriteRpc,
			" redial:", val.totalRedial, " reAcpt:", val.totalReAcpt)
		Global_TestSrvMgr.total += val.totalRpcDial
		Global_TestSrvMgr.totalReqRecv += val.totalRpcRecv

		totalRedial += val.totalRedial
		totalReAcpt += val.totalReAcpt

		if minRTT > val.minRTTDialRpc {
			minRTT = val.minRTTDialRpc
		}
		if maxRTT < val.maxRTTDialRpc {
			maxRTT = val.maxRTTDialRpc
		}
		totalRTT += val.totalRTTRpc

		if val.DialConnectionID == 0{
			noDail ++
		}
		if val.AcptConnectionID == 0{
			noAcpt ++
		}
	}

	totalDiff := Global_TestSrvMgr.total - Global_TestSrvMgr.totalLast
	totalReqDiff := Global_TestSrvMgr.totalReqRecv - Global_TestSrvMgr.totalReqRecvLast
	tnow := float64(time.Now().UnixMilli())
	tdiff := (tnow - Global_TestSrvMgr.timestamp)/float64(1000)
	aver := float64(totalDiff) / tdiff
	reqAver := float64(totalReqDiff) / tdiff
	var averRTT int64 = 0
	if Global_TestSrvMgr.total >= 1 {
		averRTT = totalRTT / Global_TestSrvMgr.total
	}
	fmt.Println(" client count:", len(Global_TestSrvMgr.mapSrv), " total:", Global_TestSrvMgr.total, " last:", Global_TestSrvMgr.totalLast,
		" totalDiff:", totalDiff, " tdiff:", tdiff, "\n aver:", aver, " req aver:", reqAver,
		" totalRedial:", totalRedial, " totalReAcpt:", totalReAcpt,
		" \n minRTT:", minRTT, " maxRTT:", maxRTT, " averRTT:", averRTT,
		" \n noDail:", noDail, " noAcpt:", noAcpt)
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
