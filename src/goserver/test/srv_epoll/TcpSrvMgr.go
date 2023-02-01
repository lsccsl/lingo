package main

import (
	"github.com/golang/protobuf/proto"
	"goserver/common"
	"goserver/msgpacket"
	"strconv"
	"time"
)



/* begin srv event */
type srvEvt_addremote struct {
	srvID int64
	addr string
	closeExpireSec int
}
type srvEvt_TcpDialSuc struct {
	srvID  int64
	fdDial common.FD_DEF
}
type srvEvt_TcpClose struct {
	srvID int64
	fd    common.FD_DEF
}
type srvEvt_SrvReport struct {
	srvID  int64
	fdAcpt common.FD_DEF
}
type srvEvt_protoMsg struct {
	srvID int64
	fd    common.FD_DEF
	msg   proto.Message
	msgType msgpacket.MSG_TYPE
}
type srvEvt_RPC struct {
	srvID int64
	rpcUUID     int64
	chRouteBack CHAN_RPC_ROUTEBACK
	msgType     msgpacket.MSG_TYPE
	msg proto.Message
	timeoutMills int64
}
type srvEvt_RPC_Del struct {
	srvID int64
	rpcUUID int64
}
type srvEvt_static struct {
	chBack chan TcpSrvMgrUnit
}

type EN_TIMER_TYPE int
const(
	EN_TIMER_TYPE_close_dial EN_TIMER_TYPE = 1
	EN_TIMER_TYPE_close_acpt EN_TIMER_TYPE = 2
	EN_TIMER_TYPE_heartbeat  EN_TIMER_TYPE = 3
)
type srvEvt_timer struct {
	srvID     int64
	timerType EN_TIMER_TYPE
	timerData interface{}
}
/* end srv event */

type TcpSrvMgrStatic struct {
	lastTotalRPCOut int64
	lastTotalRPCIn int64
	lastSampleMS int64
}
type TcpSrvMgr struct {
	srvID int64 // self

	eSrvMgr *ServerMgr

	mgrUnit []*TcpSrvMgrUnit

	rpcPool *common.CorPool

	TcpSrvMgrStatic
}


func (pthis*TcpSrvMgr)getSrvProcessUnit(srvID int64)*TcpSrvMgrUnit {
	processUnitCount := int64(len(pthis.mgrUnit))
	idx := srvID % processUnitCount
	if idx >= processUnitCount {
		return nil
	}
	pu := pthis.mgrUnit[idx]
	if pu == nil {
		return nil
	}
	return pu
}

func (pthis*TcpSrvMgr)TcpSrvMgrPushMsgToUnit(srvID int64, msg interface{}) {
	pu := pthis.getSrvProcessUnit(srvID)
	if pu != nil {
		pu.chSrv <- msg
	} else {
		common.LogErr("srv:", srvID, " push msg to srv process unit err")
	}
}

func (pthis*TcpSrvMgr)TcpSrvMgrAddRemoteSrv(srvID int64, addr string, closeExpireSec int){
	pthis.TcpSrvMgrPushMsgToUnit(srvID, &srvEvt_addremote{
		srvID : srvID,
		addr : addr,
		closeExpireSec : closeExpireSec,
	})
}

func (pthis*TcpSrvMgr)TcpSrvMgrRPCSync(srvID int64, msgType msgpacket.MSG_TYPE, protoMsg proto.Message, timeoutMilliSec int) (proto.Message, error){
	chRouteBack := make(CHAN_RPC_ROUTEBACK)
	evt := &srvEvt_RPC{
		srvID :       srvID,
		rpcUUID :     common.GenUUID64_V4(),
		chRouteBack : chRouteBack,
		msgType :     msgType,
		msg :         protoMsg,
	}

	pthis.TcpSrvMgrPushMsgToUnit(srvID, evt)

	chTimer := time.After(time.Millisecond * time.Duration(timeoutMilliSec))
	select{
	case msgRsp := <-chRouteBack:
		return msgRsp, nil
	case <-chTimer:
		pthis.TcpSrvMgrPushMsgToUnit(srvID, &srvEvt_RPC_Del{srvID,evt.rpcUUID})
	}
	close(chRouteBack)
	return nil, common.GenErrNoERR_NUM("srv:", srvID, " rpc err, msg:", protoMsg)
}

func ConstructorTcpSrvMgr(eSrvMgr *ServerMgr, srvProcessUnitCount int) *TcpSrvMgr {
	tcpSrvMgr := &TcpSrvMgr{
		eSrvMgr : eSrvMgr,
		mgrUnit : make([]*TcpSrvMgrUnit, 0, srvProcessUnitCount),
		rpcPool : common.CorPoolInit(150, 5, 300),
	}

	for i := 0; i < srvProcessUnitCount; i ++ {
		processUnit := ConstructorTcpSrvMgrUnit(tcpSrvMgr)
		tcpSrvMgr.mgrUnit = append(tcpSrvMgr.mgrUnit, processUnit)
		go processUnit._go_srvProcess_unit()
	}

	return tcpSrvMgr
}



func (pthis*TcpSrvMgr)Dump(bDetail bool)string{
	arrayMgrUnit := make([]TcpSrvMgrUnit, 0, len(pthis.mgrUnit))
	for _, val := range pthis.mgrUnit {
		chBack := make(chan TcpSrvMgrUnit)
		val.chSrv <- &srvEvt_static{chBack}
		arrayMgrUnit = append(arrayMgrUnit, <- chBack)
	}

	var str string
	serverCount := 0
	noDial := 0
	noAcpt := 0
	var totalRPCOut int64 = 0
	var totalRPCIn int64 = 0
	for _, val := range arrayMgrUnit {
		totalRPCOut += val.totalRPCOut
		totalRPCIn += val.totalRPCIn
		serverCount += len(val.mapSrv)
		for _, valSrv := range val.mapSrv {
			if bDetail {
				str += "\r\n srv:" + strconv.FormatInt(valSrv.srvID, 10) +
					" HB timestamp:" + strconv.FormatInt(valSrv.timestampLastHeartbeat, 10) +
					" fdDial:" + valSrv.fdDial.String() + " fdAcpt:" + valSrv.fdAcpt.String() +
					" addr:" + valSrv.addr
			}

			if valSrv.fdAcpt.IsNull() {
				noAcpt++
			}
			if valSrv.fdDial.IsNull() {
				noDial++
			}
		}
	}

	tnowMS := time.Now().UnixMilli()
	tdiffSec := float64(tnowMS - pthis.lastSampleMS) / 1000.0
	diffRPCOut := totalRPCOut - pthis.lastTotalRPCOut
	diffRPCIn := totalRPCIn - pthis.lastTotalRPCIn

	averRPCOut := float64(diffRPCOut) / tdiffSec
	averRPCIn :=  float64(diffRPCIn) / tdiffSec

	str += "\r\nserver count:" + strconv.Itoa(serverCount) +
		" noAcpt:" + strconv.Itoa(noAcpt) + " noDial:" + strconv.Itoa(noDial) +
		"\r\ntotalRPCOut:" + strconv.FormatInt(totalRPCOut, 10) + " totalRPCIn:" + strconv.FormatInt(totalRPCIn, 10) +
		" averRPCOut:" + strconv.FormatFloat(averRPCOut, 'f', 2,64) +
		" averRPCIn:" + strconv.FormatFloat(averRPCIn, 'f', 2,64)

	pthis.lastTotalRPCOut = totalRPCOut
	pthis.lastTotalRPCIn = totalRPCIn
	pthis.lastSampleMS = tnowMS

	return str
}
