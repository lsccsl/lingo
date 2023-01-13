package main

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"lin/lin_common"
	"lin/msgpacket"
	"lin/server/server_common"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
)

// MsgQueCenterSrv this struct is access by multi go coroutine, all member must be 'thread safe'
type MsgQueCenterSrv struct {
	lsn *lin_common.EPollListener

	mapMsgQueSrv sync.Map // server_common.MSGQUE_SRV_ID - msgQueSrvInfo

	srvIDSeed atomic.Int32
}

type msgQueSrvInfo struct {
	fd lin_common.FD_DEF
	ip string
	port int32
	queSrvID server_common.MSGQUE_SRV_ID
}

func (pthis*msgQueSrvInfo)String()(str string){
	str = pthis.queSrvID.ToString() + " fd:" + pthis.fd.String() +
		"[" + pthis.ip + ":" + strconv.FormatInt(int64(pthis.port), 10) + "]"
	return
}

type tcpAttachDataMsgQueSrv struct {
	queSrvID server_common.MSGQUE_SRV_ID
}

func (pthis*MsgQueCenterSrv)TcpAcceptConnection(fd lin_common.FD_DEF, addr net.Addr, inAttachData interface{})(outAttachData interface{}) {
	lin_common.LogDebug(addr)
	return nil
}

func (pthis*MsgQueCenterSrv)TcpDialConnection(fd lin_common.FD_DEF, addr net.Addr, inAttachData interface{})(outAttachData interface{}) {
	lin_common.LogDebug(addr)
	return nil
}

func (pthis*MsgQueCenterSrv)TcpData(fd lin_common.FD_DEF, readBuf *bytes.Buffer, inAttachData interface{})(bytesProcess int, outAttachData interface{}) {
	lin_common.LogDebug(fd)
	packType, bytesProcess, protoMsg := msgpacket.ProtoUnPacketFromBin(readBuf)
	if protoMsg == nil {
		return
	}
	lin_common.LogInfo("packType:", packType, " bytesProcess:", bytesProcess)

	outAttachData = nil
	switch msgpacket.PB_MSG_INTER_TYPE(packType) {
	case msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_REGISTER:
		{
			outAttachData = pthis.process_PB_MSG_INTER_QUESRV_REGISTER(fd, protoMsg)
			return
		}
	}

	return
}

func (pthis*MsgQueCenterSrv)TcpClose(fd lin_common.FD_DEF, closeReason lin_common.EN_TCP_CLOSE_REASON, inAttachData interface{}) {
	lin_common.LogInfo(fd, " closeReason:", closeReason, " inAttachData:", inAttachData)

	switch t := inAttachData.(type) {
	case *tcpAttachDataMsgQueSrv:
		{
			pthis.process_TcpClose_MsgQueSrv(fd, t)
		}
	}
}

func (pthis*MsgQueCenterSrv)TcpOutBandData(fd lin_common.FD_DEF, data interface{}, inAttachData interface{}) {
	lin_common.LogDebug(fd)
}




func (pthis*MsgQueCenterSrv)process_PB_MSG_INTER_QUESRV_REGISTER(fd lin_common.FD_DEF, pbMsg proto.Message) interface{}{

	regMsg, ok := pbMsg.(*msgpacket.PB_MSG_INTER_QUESRV_REGISTER)
	if !ok || regMsg == nil {
		return nil
	}

	//assign id
	qsiReg := msgQueSrvInfo{
		queSrvID : server_common.MSGQUE_SRV_ID(pthis.genQueSrvID()),
		fd :       fd,
		ip :       regMsg.Ip,
		port:      regMsg.Port,
	}
	//add msg que server list
	pthis.mapMsgQueSrv.Store(qsiReg.queSrvID, qsiReg)

	//response
	regRet := &msgpacket.PB_MSG_INTER_QUESRV_REGISTER_RES{}
	regRet.QueSrvId = int64(qsiReg.queSrvID)

	pthis.mapMsgQueSrv.Range(func(key, value any) bool{

		qsi, ok := value.(msgQueSrvInfo)
		if !ok {
			return true
		}
		queSrvInfo := &msgpacket.PB_MSG_INTER_QUESRV_INFO {
			QueSrvId: int64(qsi.queSrvID),
			Ip:qsi.ip,
			Port: qsi.port,
		}
		regRet.QueSrvInfo = append(regRet.QueSrvInfo, queSrvInfo)

		return true
	})

	pthis.SendProtoMsg(fd, msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_REGISTER_RES, regRet)

	pthis.mapMsgQueSrv.Range(func(key, value any) bool{
		qsi, ok := value.(msgQueSrvInfo)
		if !ok {
			return true
		}
		if qsi.queSrvID == qsiReg.queSrvID {
			return true
		}
		//notify other msg que srv online
		ntf := &msgpacket.PB_MSG_INTER_QUESRV_ONLINE_NTF{
			QueSrvInfo : &msgpacket.PB_MSG_INTER_QUESRV_INFO{
				QueSrvId:int64(qsiReg.queSrvID),
				Ip: qsiReg.ip,
				Port: qsiReg.port,
			},
		}
		pthis.SendProtoMsg(qsi.fd, msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_ONLINE_NTF, ntf)

		return true
	})

	return &tcpAttachDataMsgQueSrv{
		qsiReg.queSrvID,
	}
}

func (pthis*MsgQueCenterSrv)process_TcpClose_MsgQueSrv(fd lin_common.FD_DEF, attachData *tcpAttachDataMsgQueSrv) {
	lin_common.LogInfo(attachData.queSrvID.ToString())

	val, ok := pthis.mapMsgQueSrv.Load(attachData.queSrvID)
	if !ok {
		lin_common.LogErr(attachData.queSrvID.ToString(), " can't find")
		return
	}
	qsi, ok := val.(msgQueSrvInfo)
	if !ok {
		lin_common.LogErr(attachData.queSrvID.ToString(), " data convert err")
		return
	}

	if !qsi.fd.IsSame(&fd) {
		lin_common.LogErr(attachData.queSrvID.ToString(), " fd is not same, current:", qsi.fd, " close:", fd)
		return
	}

	pthis.mapMsgQueSrv.Delete(attachData.queSrvID)

	pthis.mapMsgQueSrv.Range(func(key, value any) bool{
		qsi, ok := value.(msgQueSrvInfo)
		if !ok {
			return true
		}
		//notify other msg que srv offline
		ntf := &msgpacket.PB_MSG_INTER_QUESRV_OFFLINE_NTF{
			QueSrvId:int64(attachData.queSrvID),
		}
		pthis.SendProtoMsg(qsi.fd, msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_OFFLINE_NTF, ntf)
		return true
	})
}

func (pthis*MsgQueCenterSrv)genQueSrvID() int32 {
	return pthis.srvIDSeed.Add(1)
}

func (pthis*MsgQueCenterSrv)SendProtoMsg(fd lin_common.FD_DEF, msgType msgpacket.PB_MSG_INTER_TYPE, protoMsg proto.Message){
	pthis.lsn.EPollListenerWrite(fd, msgpacket.ProtoPacketToBin(uint16(msgType), protoMsg))
}

func (pthis*MsgQueCenterSrv)Wait() {
	pthis.lsn.EPollListenerWait()
}

// ConstructMsgQueCenterSrv <addr> example 127.0.0.1:8888
func ConstructMsgQueCenterSrv(addr string, epollCoroutineCount int) *MsgQueCenterSrv {
	mqMgr := &MsgQueCenterSrv{}
	mqMgr.srvIDSeed.Store(1)

	lsn, err := lin_common.ConstructorEPollListener(mqMgr, addr, epollCoroutineCount,
		lin_common.ParamEPollListener{ParamET: true,
			ParamEpollWaitTimeoutMills:30*1000,
			ParamIdleCheckInterval:60 * 1000,
		})
	if err != nil {
		lin_common.LogErr("constructor epoll listener err:", err)
		return nil
	}
	mqMgr.lsn = lsn

	return mqMgr
}

func (pthis*MsgQueCenterSrv)Dump(bDetail bool) (str string) {

	str = "\r\nque srv id seed:" + strconv.FormatInt(int64(pthis.srvIDSeed.Load()), 10) + "\r\n"

	str += "msg que srv reg:\r\n"
	pthis.mapMsgQueSrv.Range(func(key, value any) bool{
		qsi, ok := value.(msgQueSrvInfo)
		if !ok {
			return true
		}

		str += qsi.String() + "\r\n"
		return true
	})

	return
}
