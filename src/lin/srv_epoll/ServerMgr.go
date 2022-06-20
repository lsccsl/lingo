package main

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"lin/lin_common"
	cor_pool "lin/lin_cor_pool"
	"lin/msgpacket"
	"net"
	"strconv"
	"time"
)


const(
	EN_CORPOOL_JOBTYPE_Rpc_req = cor_pool.EN_CORPOOL_JOBTYPE_user + 100
	EN_CORPOOL_JOBTYPE_client_Rpc_req
)

const(
	EN_TCP_CLOSE_REASON_timeout lin_common.EN_TCP_CLOSE_REASON = lin_common.EN_TCP_CLOSE_REASON_inter_max + 1
	EN_TCP_CLOSE_REASON_new_conn lin_common.EN_TCP_CLOSE_REASON = lin_common.EN_TCP_CLOSE_REASON_inter_max + 2
	EN_TCP_CLOSE_REASON_new_dial lin_common.EN_TCP_CLOSE_REASON = lin_common.EN_TCP_CLOSE_REASON_inter_max + 3
)

type TcpSrvAttachData struct {
	srvID int64
	isDial bool
}
type TcpCliAttachData struct {
	cliID int64
}
type TcpAttachData struct {
	TcpSrvAttachData
	TcpCliAttachData
}

type ServerMgrStatic struct {
	lastTotalRecv int64
	lastSampleMS int64

	clientClose int64
	serverDialClose int64
	serverAcptClose int64

	unknownTcpClose int64
}
type ServerMgr struct {
	lsn *lin_common.EPollListener
	processUnit []*TcpClientMgrUnit

	tcpSrvMgr *TcpSrvMgr

	clientCloseTimeoutSec int
	srvCloseTimeoutSec int

	ServerMgrStatic
}

func (pthis*ServerMgr)TcpAcceptConnection(fd lin_common.FD_DEF, addr net.Addr, inAttachData interface{}) (outAttachData interface{}){
	lin_common.LogDebug(" accept connection fd:", fd.String(), " addr:", addr)
	return nil
}
func (pthis*ServerMgr)TcpDialConnection(fd lin_common.FD_DEF, addr net.Addr, inAttachData interface{}) (outAttachData interface{}) {
	attachData, ok := inAttachData.(*TcpAttachData)
	if !ok || attachData == nil {
		lin_common.LogErr(" dial connection no attach data, fd:", fd.String(), " addr:", addr, " inAttachData:", inAttachData)
		return
	}
	//lin_common.LogDebug(" dial connection fd:", fd.String(), " addr:", addr, " srv:", attachData.srvID)

	pthis.tcpSrvMgr.TcpSrvMgrPushMsgToUnit(attachData.srvID, &srvEvt_TcpDialSuc{attachData.srvID, fd})
	return nil
}
func (pthis*ServerMgr)TcpData(fd lin_common.FD_DEF, readBuf *bytes.Buffer, inAttachData interface{})(bytesProcess int, retAttachData interface{}) {
	packType, bytesProcess, protoMsg := msgpacket.ProtoUnPacketFromBin(readBuf)
	if protoMsg == nil {
		return
	}

	var pu *TcpClientMgrUnit = nil

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
			retAttachData = &TcpAttachData{
				TcpSrvAttachData:TcpSrvAttachData {
					srvID: msgR.SrvId,
					isDial: false,
				},
			}
			pthis.tcpSrvMgr.TcpSrvMgrPushMsgToUnit(msgR.SrvId,
				&srvEvt_SrvReport{
					srvID : msgR.SrvId,
					fdAcpt : fd,
				})
			return
		}
	case msgpacket.MSG_TYPE__MSG_LOGIN:
		{
			msgL := protoMsg.(*msgpacket.MSG_LOGIN)
			retAttachData =&TcpAttachData {
				TcpCliAttachData:TcpCliAttachData{
					cliID: msgL.Id,
				},
			}
			pu = pthis.GetProcessUnitByClientID(msgL.Id)
			if pu != nil {
				pu.PushTcpLoginMsg(msgL.Id, fd)
			}
			return
		}
	}

	if inAttachData == nil {
		lin_common.LogErr("fd:", fd.String(), " no attach data")
		return
	} else {
		tcpAttachData, ok := inAttachData.(*TcpAttachData)
		if !ok {
			lin_common.LogErr("fd:", fd.String(), " unknown attach data", inAttachData)
			return
		}

		if tcpAttachData.srvID != 0 {
			pthis.tcpSrvMgr.TcpSrvMgrPushMsgToUnit(tcpAttachData.srvID,
				&srvEvt_protoMsg{
					srvID:tcpAttachData.srvID,
					fd:fd,
					msg:protoMsg,
					msgType:packType,
				})
		} else {
			pu = pthis.GetProcessUnitByClientID(tcpAttachData.cliID)
			if pu == nil {
				lin_common.LogErr("fd:", fd.String(), " not process clientid:", tcpAttachData.cliID)
				return readBuf.Len(), nil
			}
			pu.PushProtoMsg(tcpAttachData.cliID, fd, protoMsg, tcpAttachData)
		}
	}

	return
}
func (pthis*ServerMgr)TcpClose(fd lin_common.FD_DEF, closeReason lin_common.EN_TCP_CLOSE_REASON, inAttachData interface{}) {
	lin_common.LogDebug("recv tcp close ", fd.String(), " attach data:", inAttachData, " closeReason:", closeReason)
	if inAttachData == nil{
		lin_common.LogErr("fd:", fd.String(), " not attach data")
		pthis.unknownTcpClose ++
		return
	}
	tcpAttachData, ok := inAttachData.(*TcpAttachData)
	if !ok {
		lin_common.LogErr("fd:", fd.String(), " unknown attach data", inAttachData)
		pthis.unknownTcpClose ++
		return
	}
	if tcpAttachData.srvID != 0 {
		pthis.tcpSrvMgr.TcpSrvMgrPushMsgToUnit(tcpAttachData.srvID,
			&srvEvt_TcpClose{
				srvID : tcpAttachData.srvID,
				fd : fd,
			})

		if tcpAttachData.isDial {
			pthis.serverDialClose ++
		} else {
			pthis.serverAcptClose ++
		}
	} else {
		pu := pthis.GetProcessUnitByClientID(tcpAttachData.cliID)
		if pu != nil {
			pu.PushTcpCloseMsg(tcpAttachData.cliID, fd)
		}

		pthis.clientClose ++
	}
}

func (pthis*ServerMgr)GetProcessUnitByClientID(cliID int64) *TcpClientMgrUnit {
	processUnitCount := int64(len(pthis.processUnit))
	idx := cliID % processUnitCount
	if idx >= processUnitCount {
		return nil
	}
	pu := pthis.processUnit[idx]
	if pu == nil {
		return nil
	}
	return pu
}

func (pthis*ServerMgr)AddRemoteSrv(srvID int64, addr string, closeExpireSec int) {
	pthis.tcpSrvMgr.TcpSrvMgrAddRemoteSrv(srvID, addr, closeExpireSec)
}

func (pthis*ServerMgr)SendProtoMsg(fd lin_common.FD_DEF, msgType msgpacket.MSG_TYPE, protoMsg proto.Message){
	pthis.lsn.EPollListenerWrite(fd, msgpacket.ProtoPacketToBin(msgType, protoMsg))
}

func (pthis*ServerMgr)Dump(bDetail bool)string{
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
		" byte unProc:" + strconv.FormatInt(es.ByteRecv - es.ByteProc, 10) +
		" clientClose:" + strconv.FormatInt(pthis.clientClose, 10) +
		" serverDialClose:" + strconv.FormatInt(pthis.serverDialClose, 10) +
		" serverAcptClose:" + strconv.FormatInt(pthis.serverAcptClose, 10) +
		" unknownTcpClose:" + strconv.FormatInt(pthis.unknownTcpClose, 10) +
		"\r\n\r\n"

	str += pthis.tcpSrvMgr.Dump(bDetail)

	pthis.lastSampleMS = tnowMs
	pthis.lastTotalRecv = totalRecv

	return str
}

func ConstructorEpollServerMgr(addr string,
	processUnitCount int, srvProcessUnitCount int,
	epollCoroutineCount int, clientCloseTimeoutSec int, srvCloseTimeoutSec int,
	bET bool) (*ServerMgr, error) {
	defer func() {
		err := recover()
		if err != nil {
			lin_common.LogErr(err)
		}
	}()

	lin_common.LogDebug("processUnitCount:", processUnitCount, "epollCoroutineCount:", epollCoroutineCount, "clientCloseTimeoutSec:", clientCloseTimeoutSec,
		"bET:", bET)

	eSrvMgr := &ServerMgr{
		processUnit : make([]*TcpClientMgrUnit, 0, processUnitCount),
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
