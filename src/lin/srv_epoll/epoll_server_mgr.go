package main

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"lin/lin_common"
	"lin/msgpacket"
	"net"
)


/* begin process unit msg define */
type msgProto struct {
	fd lin_common.FD_DEF
	packType msgpacket.MSG_TYPE
	protoMsg proto.Message
}
/* end process unit msg define */


type EpollServerMgr struct {
	lsn *lin_common.EPollListener
	processUnit []*eSrvMgrProcessUnit
}

func (pthis*EpollServerMgr)TcpAcceptConnection(fd lin_common.FD_DEF, addr net.Addr) {
	lin_common.LogDebug(" accept connection fd:", fd.FD, " magic:", fd.Magic, " addr:", addr)
	//ConstructorTcpClient
}
func (pthis*EpollServerMgr)TcpDialConnection(fd lin_common.FD_DEF, addr net.Addr) {

}
func (pthis*EpollServerMgr)TcpData(fd lin_common.FD_DEF, readBuf *bytes.Buffer)(bytesProcess int) {
	processUnitCount := len(pthis.processUnit)
	if processUnitCount <= 0 {
		return readBuf.Len()
	}

	idx := fd.FD % processUnitCount
	if idx >= processUnitCount {
		return readBuf.Len()
	}
	pu := pthis.processUnit[idx]
	if pu == nil {
		return readBuf.Len()
	}

	packType, bytesProcess, protoMsg := msgpacket.ProtoUnPacketFromBin(readBuf)
	if protoMsg == nil {
		//log.LogErr("can't parse msg:", tcpConn.ByteRecv, " proc:", tcpConn.ByteProc)
		return
	}

	pu.chMsg <- &msgProto{fd, packType,protoMsg}
	return
}
func (pthis*EpollServerMgr)TcpClose(fd lin_common.FD_DEF) {
	lin_common.LogDebug("tcp close:", fd.String())
}


func ConstructorEpollServerMgr(addr string, processUnitCount int) (*EpollServerMgr, error) {
	defer func() {
		err := recover()
		if err != nil {
			lin_common.LogErr(err)
		}
	}()
	msgpacket.InitMsgParseVirtualTable()

	eSrvMgr := &EpollServerMgr{
		processUnit : make([]*eSrvMgrProcessUnit, 0, processUnitCount),
	}
	lsn, err := lin_common.ConstructorEPollListener(eSrvMgr, addr, 10, lin_common.ParamEPollListener{})
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
