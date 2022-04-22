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

type EpollServerMgrStatic struct {
	lastTotalRecv int64
	lastSampleMS int64
}
type EpollServerMgr struct {
	lsn *lin_common.EPollListener
	processUnit []*eSrvMgrProcessUnit

	clientCloseTimeoutSec int

	EpollServerMgrStatic
}

func (pthis*EpollServerMgr)TcpAcceptConnection(fd lin_common.FD_DEF, addr net.Addr) {
	lin_common.LogDebug(" accept connection fd:", fd.FD, " magic:", fd.Magic, " addr:", addr)
	//ConstructorTcpClient
}
func (pthis*EpollServerMgr)TcpDialConnection(fd lin_common.FD_DEF, addr net.Addr) {

}
func (pthis*EpollServerMgr)TcpData(fd lin_common.FD_DEF, readBuf *bytes.Buffer)(bytesProcess int) {
	pu := pthis.GetProcessUnitByFD(fd)
	if pu == nil {
		return readBuf.Len()
	}

	packType, bytesProcess, protoMsg := msgpacket.ProtoUnPacketFromBin(readBuf)
	if protoMsg == nil {
		//log.LogErr("can't parse msg:", tcpConn.ByteRecv, " proc:", tcpConn.ByteProc)
		return
	}

	pu.chMsg <- &msgProto{fd, packType,protoMsg}

	// todo:delete for test
	if packType == msgpacket.MSG_TYPE__MSG_TEST {
		msgTest := protoMsg.(*msgpacket.MSG_TEST)
		msgTest.TimestampArrive = time.Now().UnixMilli()
	}

	return
}
func (pthis*EpollServerMgr)TcpClose(fd lin_common.FD_DEF) {
	//lin_common.LogDebug("tcp close:", fd.String())
	pu := pthis.GetProcessUnitByFD(fd)
	if pu != nil {
		pu.chMsg <- &msgTcpClose{fd}
	}
}

func (pthis*EpollServerMgr)GetProcessUnitByFD(fd lin_common.FD_DEF) *eSrvMgrProcessUnit {
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

func ConstructorEpollServerMgr(addr string, processUnitCount int, clientCloseTimeoutSec int) (*EpollServerMgr, error) {
	defer func() {
		err := recover()
		if err != nil {
			lin_common.LogErr(err)
		}
	}()
	msgpacket.InitMsgParseVirtualTable()

	eSrvMgr := &EpollServerMgr{
		processUnit : make([]*eSrvMgrProcessUnit, 0, processUnitCount),
		clientCloseTimeoutSec : clientCloseTimeoutSec,
	}
	lsn, err := lin_common.ConstructorEPollListener(eSrvMgr, addr, 20, lin_common.ParamEPollListener{ParamET: true})
	if err != nil {
		lin_common.LogErr("constructor epoll listener err:", err)
		return nil, err
	}
	eSrvMgr.lsn = lsn

	for i := 0; i < processUnitCount; i ++ {
		pu := &eSrvMgrProcessUnit{
			chMsg : make(chan interface{}),
			eSrvMgr : eSrvMgr,
			mapClient : make(MAP_CLIENT),
		}
		eSrvMgr.processUnit = append(eSrvMgr.processUnit, pu)
		go pu._go_Process_unit()
	}

	return eSrvMgr, nil
}
