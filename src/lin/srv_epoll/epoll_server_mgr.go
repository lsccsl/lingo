package main

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"lin/lin_common"
	"lin/msgpacket"
	"net"
	"strconv"
	"time"
)


/* begin process unit msg define */
type msgProto struct {
	fd lin_common.FD_DEF
	packType msgpacket.MSG_TYPE
	protoMsg proto.Message
}
type msgTimer struct {
	fd lin_common.FD_DEF
	timerData int
}
type msgTcpClose struct {
	fd lin_common.FD_DEF
}
/* end process unit msg define */

type TcpAttachData struct {
	srvID int64
	cliID int64
}

type EpollServerMgrStatic struct {
	lastTotalRecv int64
	lastSampleMS int64
}
type EpollServerMgr struct {
	lsn *lin_common.EPollListener
	processUnit []*EPollProcessUnit

	tcpSrvMgr *TcpSrvMgr

	clientCloseTimeoutSec int
	srvCloseTimeoutSec int

	EpollServerMgrStatic
}

func (pthis*EpollServerMgr)TcpAcceptConnection(fd lin_common.FD_DEF, addr net.Addr, tpcAttachData interface{}) {
	lin_common.LogDebug(" accept connection fd:", fd.String(), " addr:", addr)
	//ConstructorTcpClient
}
func (pthis*EpollServerMgr)TcpDialConnection(fd lin_common.FD_DEF, addr net.Addr, tpcAttachData interface{}) {
	lin_common.LogDebug(" dial connection fd:", fd.String(), " addr:", addr)
	attachData, ok := tpcAttachData.(*TcpAttachData)
	if !ok || attachData == nil {
		return
	}

	pthis.tcpSrvMgr.TcpSrvMgrPushMsgToUnit(attachData.srvID, srvEvt_TcpDialSuc{attachData.srvID, fd})
}
func (pthis*EpollServerMgr)TcpData(fd lin_common.FD_DEF, readBuf *bytes.Buffer, tpcAttachData interface{})(bytesProcess int, attachData interface{}) {
	packType, bytesProcess, protoMsg := msgpacket.ProtoUnPacketFromBin(readBuf)
	if protoMsg == nil {
		return
	}

	switch packType {
	case msgpacket.MSG_TYPE__MSG_TEST:
		{
			// todo:delete for test
			msgTest := protoMsg.(*msgpacket.MSG_TEST)
			msgTest.TimestampArrive = time.Now().UnixMilli()
		}
	case msgpacket.MSG_TYPE__MSG_SRV_REPORT:
		{
			msgR := protoMsg.(*msgpacket.MSG_SRV_REPORT)
			attachData = &TcpAttachData{srvID: msgR.SrvId}
			pthis.tcpSrvMgr.TcpSrvMgrPushMsgToUnit(msgR.SrvId,
				&srvEvt_TcpAcpt{
					srvID : msgR.SrvId,
					fdAcpt : fd,
				})

			return
		}
	}

	if attachData != nil {
		tcpAttachData, _ := attachData.(*TcpAttachData)
		if tcpAttachData != nil {
			pthis.tcpSrvMgr.TcpSrvMgrPushMsgToUnit(tcpAttachData.srvID,
				&srvEvt_protoMsg{
					srvID:tcpAttachData.srvID,
					fd:fd,
					msg:protoMsg,
				})
		}
	} else {
		pu := pthis.GetProcessUnitByFD(fd)
		if pu == nil {
			return readBuf.Len(), nil
		}
		pu.chMsg <- &msgProto{fd, packType,protoMsg}

	}
	return
}
func (pthis*EpollServerMgr)TcpClose(fd lin_common.FD_DEF, attachData interface{}) {
	if attachData != nil{
		tcpAttachData, _ := attachData.(*TcpAttachData)
		if tcpAttachData != nil {
			pthis.tcpSrvMgr.TcpSrvMgrPushMsgToUnit(tcpAttachData.srvID,
				&srvEvt_TcpClose{
					srvID : tcpAttachData.srvID,
					fd : fd,
				})
		}
		return
	} else {
		lin_common.LogDebug("tcp close:", fd.String())
		pu := pthis.GetProcessUnitByFD(fd)
		if pu != nil {
			pu.chMsg <- &msgTcpClose{fd}
		}
	}
}

func (pthis*EpollServerMgr)GetProcessUnitByFD(fd lin_common.FD_DEF) *EPollProcessUnit {
	processUnitCount := len(pthis.processUnit)
	idx := fd.FD % processUnitCount
	if idx >= processUnitCount {
		return nil
	}
	pu := pthis.processUnit[idx]
	if pu == nil {
		return nil
	}
	return pu
}

func (pthis*EpollServerMgr)AddRemoteSrv(srvID int64, addr string, closeExpireSec int) {

}

func (pthis*EpollServerMgr)SendProtoMsg(fd lin_common.FD_DEF, msgType msgpacket.MSG_TYPE, protoMsg proto.Message){
	pthis.lsn.EPollListenerWrite(fd, msgpacket.ProtoPacketToBin(msgType, protoMsg))
}

func (pthis*EpollServerMgr)Dump(bDetail bool)string{
	tnowMs := time.Now().UnixMilli()

	var es lin_common.EPollListenerStatic
	pthis.lsn.EPollListenerGetStatic(&es)

	var totalRecv int64 = 0
	var totalClient int = 0
	for _, val := range pthis.processUnit {
		totalRecv += val.totalRecv
		totalClient += val.clientCount
	}

	var diff = float64(tnowMs - pthis.lastSampleMS) / 1000
	averRecv := float64(totalRecv - pthis.lastTotalRecv) / diff

	str := "averRecv:" + strconv.FormatFloat(averRecv, 'f', 2,64) +
		" tcp conn count:" + strconv.FormatInt(int64(es.TcpConnCount), 10) +
		" totalClient:" + strconv.FormatInt(int64(totalClient), 10) +
		" totalRecv:" + strconv.FormatInt(totalRecv, 10) +
		"\r\n" +
		" TcpCloseCount:" + strconv.FormatInt(es.TcpCloseCount, 10) +
		" byteRecv:" + strconv.FormatInt(es.ByteRecv, 10) +
		" byteSend:" + strconv.FormatInt(es.ByteSend, 10) +
		" byteProc:" + strconv.FormatInt(es.ByteProc, 10) +
		" byte unProc:" + strconv.FormatInt(es.ByteRecv - es.ByteProc, 10) + "\r\n\r\n"

	pthis.lastSampleMS = tnowMs
	pthis.lastTotalRecv = totalRecv

	return str
}

func ConstructorEpollServerMgr(addr string,
	processUnitCount int, srvProcessUnitCount int,
	epollCoroutineCount int, clientCloseTimeoutSec int, srvCloseTimeoutSec int,
	bET bool) (*EpollServerMgr, error) {
	defer func() {
		err := recover()
		if err != nil {
			lin_common.LogErr(err)
		}
	}()

	lin_common.LogDebug("processUnitCount:", processUnitCount, "epollCoroutineCount:", epollCoroutineCount, "clientCloseTimeoutSec:", epollCoroutineCount,
		"bET:", bET)

	msgpacket.InitMsgParseVirtualTable()

	eSrvMgr := &EpollServerMgr{
		processUnit : make([]*EPollProcessUnit, 0, processUnitCount),
		clientCloseTimeoutSec : clientCloseTimeoutSec,
		srvCloseTimeoutSec : srvCloseTimeoutSec,
	}
	eSrvMgr.tcpSrvMgr = ConstructorTcpSrvMgr(eSrvMgr, srvProcessUnitCount)
	lsn, err := lin_common.ConstructorEPollListener(eSrvMgr, addr, epollCoroutineCount, lin_common.ParamEPollListener{ParamET: bET})
	if err != nil {
		lin_common.LogErr("constructor epoll listener err:", err)
		return nil, err
	}
	eSrvMgr.lsn = lsn

	for i := 0; i < processUnitCount; i ++ {
		pu := ConstructorEPollProcessUnit(eSrvMgr)
		eSrvMgr.processUnit = append(eSrvMgr.processUnit, pu)
		go pu._go_Process_unit()
	}

	return eSrvMgr, nil
}
