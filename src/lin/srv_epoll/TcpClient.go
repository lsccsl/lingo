package main

import (
	"lin/lin_common"
	"lin/msgpacket"
	"net"
	"runtime"
	"time"
)

type MAP_CLIENT_STATIC map[msgpacket.MSG_TYPE]int64

type TcpClient struct {
	fd lin_common.FD_DEF
	addr *net.Addr

	clientID int64

	timerConnClose * time.Timer
	durationClose time.Duration
	pu *eSrvMgrProcessUnit
}

func ConstructorTcpClient(pu *eSrvMgrProcessUnit, fd lin_common.FD_DEF, clientID int64) *TcpClient {
	tc := &TcpClient{
		fd : fd,
		pu : pu,
		clientID : clientID,
		durationClose : time.Second*time.Duration(pu.eSrvMgr.clientCloseTimeoutSec),
	}
	runtime.SetFinalizer(tc, (*TcpClient).Destructor)
	tc.timerConnClose = time.AfterFunc(tc.durationClose,
		func(){
			lin_common.LogDebug("timeout close clientid:", tc.clientID, " fd:", tc.fd.String())
			tc.pu.eSrvMgr.lsn.EPollListenerCloseTcp(tc.fd)
/*			tc.timerConnClose.Stop()
			tc.timerConnClose = nil*/
		})

	return tc
}

func (pthis*TcpClient)Destructor() {
	lin_common.LogDebug(" clientid:", pthis.clientID, " fd:", pthis.fd.String())
	runtime.SetFinalizer(pthis, nil)
	if pthis.timerConnClose != nil {
		pthis.timerConnClose.Stop()
		pthis.timerConnClose = nil
	}
}

func (pthis*TcpClient)Process_MSG_TCP_STATIC(msg *msgpacket.MSG_TCP_STATIC) {
	lin_common.LogDebug(" seq:", msg.Seq)

	msgRes := &msgpacket.MSG_TCP_STATIC_RES{
		ByteRecv:0,
		ByteProc:0,
		ByteSend:0,
	}
	pthis.pu.eSrvMgr.lsn.EPollListenerWrite(pthis.fd, msgpacket.ProtoPacketToBin(msgpacket.MSG_TYPE__MSG_TCP_STATIC_RES, msgRes))
}
func (pthis*TcpClient)Process_MSG_TEST(msg *msgpacket.MSG_TEST) {
	//lin_common.LogDebug("clientid:", pthis.clientID, " fd:", pthis.fd.String())
	msgRes := &msgpacket.MSG_TEST_RES{}
	msgRes.Id = msg.Id
	msgRes.Str = msg.Str
	msgRes.Seq = msg.Seq
	msgRes.Timestamp = msg.Timestamp
	msgRes.TimestampArrive = msg.TimestampArrive
	msgRes.TimestampProcess = time.Now().UnixMilli()
	pthis.pu.eSrvMgr.lsn.EPollListenerWrite(pthis.fd, msgpacket.ProtoPacketToBin(msgpacket.MSG_TYPE__MSG_TEST_RES, msgRes))

}
func (pthis*TcpClient)Process_MSG_HEARTBEAT(msg *msgpacket.MSG_HEARTBEAT) {
	lin_common.LogDebug("clientid:", pthis.clientID, " fd:", pthis.fd.String())
	msgRes := &msgpacket.MSG_HEARTBEAT_RES{}
	msgRes.Id = msg.Id
	pthis.pu.eSrvMgr.lsn.EPollListenerWrite(pthis.fd, msgpacket.ProtoPacketToBin(msgpacket.MSG_TYPE__MSG_HEARTBEAT_RES, msgRes))
}

func (pthis*TcpClient)Process_protoMsg(msg *msgProto) {
	pthis.timerConnClose.Reset(pthis.durationClose)

	switch t := msg.protoMsg.(type) {
	case *msgpacket.MSG_TEST:
		pthis.Process_MSG_TEST(t)
	case *msgpacket.MSG_HEARTBEAT:
		pthis.Process_MSG_HEARTBEAT(t)
	case *msgpacket.MSG_TCP_STATIC:
		pthis.Process_MSG_TCP_STATIC(t)
	}
}


