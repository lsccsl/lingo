package main

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"lin/lin_common"
	"lin/msgpacket"
	"lin/server_common"
	"net"
	"sync"
	"sync/atomic"
)

type MsgQueCenterSrv struct {
	lsn *lin_common.EPollListener

	mapMsgQueSrv sync.Map // 读多写少

	queSrvIDSeed atomic.Int32
}

type msgQueSrvInfo struct {
	fd lin_common.FD_DEF
	ip string
	port int32
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

	switch msgpacket.PB_MSG_INTER_TYPE(packType) {
	case msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_REGISTER:
		{
			pthis.processMsgQueReg(fd, protoMsg)
		}
	}

	return 0,nil
}

func (pthis*MsgQueCenterSrv)TcpClose(fd lin_common.FD_DEF, closeReason lin_common.EN_TCP_CLOSE_REASON, inAttachData interface{}) {
	lin_common.LogDebug(fd)
}

func (pthis*MsgQueCenterSrv)TcpOutBandData(fd lin_common.FD_DEF, data interface{}, inAttachData interface{}) {
	lin_common.LogDebug(fd)
}




func (pthis*MsgQueCenterSrv)processMsgQueReg(fd lin_common.FD_DEF, pbMsg proto.Message) {

	regMsg, ok := pbMsg.(*msgpacket.PB_MSG_INTER_QUESRV_REGISTER)
	if !ok || regMsg == nil {
		return
	}

	//分配id
	qsi := msgQueSrvInfo{
		queSrvID : server_common.MSGQUE_SRV_ID(pthis.genQueSrvID()),
		fd : fd,
		ip : regMsg.Ip,
		port: regMsg.Port,
	}
	//加入msg que server列表
	pthis.mapMsgQueSrv.Store(qsi.queSrvID, qsi)

	//回包
	regRet := &msgpacket.PB_MSG_INTER_QUESRV_REGISTER_RES{}
	regRet.QueSrvId = int32(qsi.queSrvID)

	pthis.mapMsgQueSrv.Range(func(key, value any) bool{

		qsi := value.(msgQueSrvInfo)
		queSrvInfo := &msgpacket.PB_MSG_INTER_QUESRV_INFO {
			QueSrvId: int32(qsi.queSrvID),
			Ip:qsi.ip,
			Port: qsi.port,
		}
		regRet.QueSrvInfo = append(regRet.QueSrvInfo, queSrvInfo)

		return true
	})

	pthis.SendProtoMsg(fd, msgpacket.PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_REGISTER_RES, regRet)
}

func (pthis*MsgQueCenterSrv)genQueSrvID() int32 {
	return pthis.queSrvIDSeed.Add(1)
}

func (pthis*MsgQueCenterSrv)SendProtoMsg(fd lin_common.FD_DEF, msgType msgpacket.PB_MSG_INTER_TYPE, protoMsg proto.Message){
	pthis.lsn.EPollListenerWrite(fd, msgpacket.ProtoPacketToBin(uint16(msgType), protoMsg))
}

func (pthis*MsgQueCenterSrv)Wait() {
	pthis.lsn.EPollListenerWait()
}

// ConstructMsgQueCenterSrv <addr> example 127.0.0.1:8888
func ConstructMsgQueCenterSrv(addr string, epollCoroutineCount int) *MsgQueCenterSrv{
	mqMgr := &MsgQueCenterSrv{}

	lsn, err := lin_common.ConstructorEPollListener(mqMgr, addr, epollCoroutineCount, lin_common.ParamEPollListener{ParamET: true})
	if err != nil {
		lin_common.LogErr("constructor epoll listener err:", err)
		return nil
	}
	mqMgr.lsn = lsn

	return mqMgr
}

