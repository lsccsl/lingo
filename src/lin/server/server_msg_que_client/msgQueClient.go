package msgque_client

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"lin/lin_common"
	"lin/msgpacket"
	"lin/server/server_common"
	"net"
	"strconv"
	"sync"
)

type MgrQueClient struct {
	lsn *lin_common.EPollListener
	srvUUID server_common.SRV_ID
	sryType server_common.SRV_TYPE

	fdCenter lin_common.FD_DEF
	msgqueCenterAddr string // "x.x.x.x:n"

	fdQueSrv lin_common.FD_DEF
	queSrvAddr string
	queDialWait sync.WaitGroup

	msgMgr *ClientSrvMsgMgr
}



type tcpAttachDataMsgQueCenter struct {
}
type tcpAttachDataMsgQueDial struct {

}

func (pthis*MgrQueClient)TcpAcceptConnection(fd lin_common.FD_DEF, addr net.Addr, inAttachData interface{})(outAttachData interface{}) {
	return nil
}

func (pthis*MgrQueClient)TcpDialConnection(fd lin_common.FD_DEF, addr net.Addr, inAttachData interface{})(outAttachData interface{}) {

	switch t := inAttachData.(type) {
	case *tcpAttachDataMsgQueCenter:
		{
			pthis.process_TcpDialToMsgCenter(fd)
		}

	default:
		{
			lin_common.LogDebug(fd, addr, t)
		}
	}

	return nil
}

func (pthis*MgrQueClient)TcpData(fd lin_common.FD_DEF, readBuf *bytes.Buffer, inAttachData interface{})(bytesProcess int, outAttachData interface{}) {
	packType, bytesProcess, protoMsg := msgpacket.ProtoUnPacketFromBin(readBuf)

	switch msgpacket.PB_MSG_TYPE(packType) {
	case msgpacket.PB_MSG_TYPE__PB_MSG_INTER_CLISRV_REG_MSGQUE_CENTER_RES:
		{
			pthis.process_PB_MSG_INTER_CLISRV_REG_MSGQUE_CENTER_RES(fd, protoMsg)
		}

	case msgpacket.PB_MSG_TYPE__PB_MSG_INTER_CLISRV_REG_TO_QUE_RES:
		{
			pthis.queDialWait.Done()
		}

	case msgpacket.PB_MSG_TYPE__PB_MSG_INTER_MSG:
		{
			pthis.process_PB_MSG_INTER_MSG(fd, protoMsg, inAttachData)
		}

	case msgpacket.PB_MSG_TYPE__PB_MSG_INTER_MSG_RES:
		{
			pthis.process_PB_MSG_INTER_MSG_RES(fd, protoMsg, inAttachData)
		}

	case msgpacket.PB_MSG_TYPE__PB_MSG_INTER_CLISRV_HEARTBEAT_RES:
		{
			lin_common.LogDebug(fd, "PB_MSG_INTER_CLISRV_HEARTBEAT_RES", " attach:", inAttachData)
		}

	default:
		{
			lin_common.LogDebug(fd, "packType:", packType, " bytesProcess:", bytesProcess, " proto msg", protoMsg, " attach:", inAttachData)
		}
	}

	return
}

func (pthis*MgrQueClient)TcpClose(fd lin_common.FD_DEF, closeReason lin_common.EN_TCP_CLOSE_REASON, inAttachData interface{}) {
	// redial to que center
	switch inAttachData.(type) {
	case *tcpAttachDataMsgQueDial:
		{
			var err error
			pthis.fdCenter, err = pthis.lsn.EPollListenerDial(pthis.msgqueCenterAddr, &tcpAttachDataMsgQueCenter{}, false)
			if err != nil {
				lin_common.LogErr("dial to msg que center err:", err)
			}
		}
	}
}

func (pthis*MgrQueClient)TcpOutBandData(fd lin_common.FD_DEF, data interface{}, inAttachData interface{}) {
}

func (pthis*MgrQueClient)TcpTick(fd lin_common.FD_DEF, tNowMill int64, inAttachData interface{}) {
	//lin_common.LogDebug(fd, " tNowMill:", tNowMill, " inAttachData:", inAttachData)
	pbHB := &msgpacket.PB_MSG_INTER_CLISRV_HEARTBEAT{SrvUuid: int64(pthis.srvUUID)}
	pthis.SendProtoMsg(fd, msgpacket.PB_MSG_TYPE__PB_MSG_INTER_CLISRV_HEARTBEAT, pbHB)
}



func (pthis*MgrQueClient)process_PB_MSG_INTER_CLISRV_REG_MSGQUE_CENTER_RES(fd lin_common.FD_DEF, pbMsg proto.Message) {
	lin_common.LogDebug(fd, pbMsg)
	// dial to msg que srv
	pbRes := pbMsg.(*msgpacket.PB_MSG_INTER_CLISRV_REG_MSGQUE_CENTER_RES)

	if pbRes.Res != msgpacket.PB_RESPONSE_CODE_PB_RESPONSE_CODE_OK {
		lin_common.LogErr("reg to msg que center err:", pbRes.Res)
	}

	pthis.srvUUID = server_common.SRV_ID(pbRes.SrvUuid)
	pthis.queSrvAddr = pbRes.QueSrvIp + ":" + strconv.FormatInt(int64(pbRes.QueSrvPort), 10)

	var err error
	pthis.fdQueSrv, err = pthis.lsn.EPollListenerDial(pthis.queSrvAddr, &tcpAttachDataMsgQueDial{}, true)
	if err != nil {
		lin_common.LogErr("dial to msg que center err:", err)
	}

	// send reg to msq que srv
	{
		pbReg := &msgpacket.PB_MSG_INTER_CLISRV_REG_TO_QUE{
			SrvUuid: int64(pthis.srvUUID),
			SrvType: int32(pthis.sryType),
		}
		pthis.SendProtoMsg(pthis.fdQueSrv, msgpacket.PB_MSG_TYPE__PB_MSG_INTER_CLISRV_REG_TO_QUE, pbReg)
	}

	pthis.lsn.EPollListenerCloseTcp(fd, server_common.EN_TCP_CLOSE_REASON_srv_reg_ok)
}

func (pthis*MgrQueClient)process_TcpDialToMsgCenter(fd lin_common.FD_DEF) {
	lin_common.LogDebug("dial to msg center success")
	pbReg := &msgpacket.PB_MSG_INTER_CLISRV_REG_MSGQUE_CENTER{
		SrvUuid: int64(pthis.srvUUID),
		SrvType: int32(pthis.sryType),
	}

	pthis.SendProtoMsg(fd, msgpacket.PB_MSG_TYPE__PB_MSG_INTER_CLISRV_REG_MSGQUE_CENTER, pbReg)
}

func (pthis*MgrQueClient)SendProtoMsg(fd lin_common.FD_DEF, msgType msgpacket.PB_MSG_TYPE, protoMsg proto.Message){
	pthis.lsn.EPollListenerWrite(fd, msgpacket.ProtoPacketToBin(uint16(msgType), protoMsg))
}


func (pthis*MgrQueClient)MgrQueClientSetCallBack(cb MsgProcessCB) {
	pthis.msgMgr.cb = cb
}

func (pthis*MgrQueClient)Wait() {
	pthis.lsn.EPollListenerWait()
}

func ConstructMgrQueClient(msgqueCenterAddr string, srvType server_common.SRV_TYPE) *MgrQueClient {
	mqCli := &MgrQueClient{
		srvUUID : server_common.SRV_ID_INVALID,
		sryType: srvType,
		msgMgr:ConstructClientSrvMsgMgr(),
	}

	lsn, err := lin_common.ConstructorEPollListener(mqCli, "", 1,
		lin_common.ParamEPollListener{
			ParamET: true,
			ParamEpollWaitTimeoutMills: 30 * 1000,
			ParamIdleClose: 600*1000,
			ParamNeedTick: true,
		})
	if err != nil {
		lin_common.LogErr("constructor epoll listener err:", err)
		return nil
	}
	mqCli.lsn = lsn

	mqCli.queDialWait.Add(1)
	// dial to msg que center server
	mqCli.msgqueCenterAddr = msgqueCenterAddr
	mqCli.fdCenter, err = lsn.EPollListenerDial(msgqueCenterAddr, &tcpAttachDataMsgQueCenter{}, false)
	if err != nil {
		lin_common.LogErr("dial to msg que center err:", err)
	}

	mqCli.queDialWait.Wait()
	lin_common.LogInfo("connect end~~~~~~~~~~~~~~~~~~~~~")

	return mqCli
}

func (pthis*MgrQueClient)Dump(bDetail bool) string {
	str := pthis.lsn.EPollListenerDump()

	str += "fdQueSrv:" + pthis.fdQueSrv.String() + " addr:" + pthis.queSrvAddr
	return str
}
