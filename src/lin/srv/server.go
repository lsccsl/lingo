package main

import (
	"github.com/golang/protobuf/proto"
	"lin/lin_common"
	"lin/log"
	msgpacket "lin/msgpacket"
	"sync/atomic"
	"time"
)

type Server struct {
	srvMgr *ServerMgr
	srvID int64
	connDial *TcpConnection
	connAcptID TCP_CONNECTION_ID
	chSrvProtoMsg chan *interProtoMsg
	chInterMsg chan interface{}
	heartbeatIntervalSec int
	rpcMgr *RPCManager

	isStopProcess int32
}

type interMsgSrvReport struct {
	tcpAccept * TcpConnection
}
type interMsgConnDial struct {
	tcpDial * TcpConnection
}

func ConstructServer(srvMgr *ServerMgr, srvID int64, heartbeatIntervalSec int)*Server {
	s := &Server{
		srvMgr:srvMgr,
		srvID:srvID,
		chSrvProtoMsg:make(chan *interProtoMsg, 100),
		chInterMsg:make(chan interface{}, 100),
		heartbeatIntervalSec:heartbeatIntervalSec,
		rpcMgr:ConstructRPCManager(),
		isStopProcess:0,
	}
	go s.go_serverProcess()
	return s
}

func (pthis*Server) go_serverProcess() {
	defer func() {
		atomic.StoreInt32(&pthis.isStopProcess, 1)
		err := recover()
		if err != nil {
			log.LogErr(err)
		}
	}()

	chTimer := time.After(time.Second * time.Duration(pthis.heartbeatIntervalSec))

	MSG_LOOP:
	for {
		select {
		case ProtoMsg := <- pthis.chSrvProtoMsg:
			{
				if ProtoMsg == nil {
					break MSG_LOOP
				}
				log.LogDebug(ProtoMsg)
			}

		case interMsg := <- pthis.chInterMsg:
			{
				switch t:= interMsg.(type) {
				case *interMsgSrvReport:
					pthis.processSrvReport(t.tcpAccept)
				case *interMsgConnDial:
					pthis.processDailConnect(t.tcpDial)
				}
			}

		case <-chTimer:
			{
				chTimer = time.After(time.Second * time.Duration(pthis.heartbeatIntervalSec))
				//send heartbeat
				msgTest := &msgpacket.MSG_TEST{}
				msgTest.Id = pthis.srvID
				pthis.connDial.TcpConnectSendProtoMsg(msgpacket.MSG_TYPE__MSG_TEST, msgTest)
			}
		}
	}

	atomic.StoreInt32(&pthis.isStopProcess, 1)
	close(pthis.chSrvProtoMsg)
	close(pthis.chInterMsg)
}

func (pthis*Server) ServerClose() {
	pthis.srvMgr.tcpMgr.TcpMgrCloseConn(pthis.connAcptID)
	if pthis.connDial != nil {
		pthis.srvMgr.tcpMgr.TcpMgrCloseConn(pthis.connDial.TcpConnectionID())
	}
	pthis.chSrvProtoMsg <- nil
}

func (pthis*Server) ServerCloseAndDelDialData() {
	pthis.srvMgr.tcpMgr.TcpDialDelDialData(pthis.srvID)
	pthis.ServerClose()
}

func (pthis*Server)processSrvReport(tcpAccept * TcpConnection){
	pthis.connAcptID = tcpAccept.TcpConnectionID()

	log.LogDebug(pthis.srvID, " ", pthis)
}

func (pthis*Server)processDailConnect(tcpDial * TcpConnection){
	pthis.connDial = tcpDial

	log.LogDebug(pthis.srvID, " ", pthis)
}

func (pthis*Server)PushInterMsg(msg interface{}){
	if atomic.LoadInt32(&pthis.isStopProcess) == 1 {
		return
	}
	pthis.chInterMsg <- msg
}
func (pthis*Server)PushProtoMsg(msgType msgpacket.MSG_TYPE, protoMsg proto.Message){
	if atomic.LoadInt32(&pthis.isStopProcess) == 1 {
		return
	}
	pthis.chSrvProtoMsg <- &interProtoMsg{
		msgType:msgType,
		protoMsg:protoMsg,
	}
}

func (pthis*Server)Go_ProcessRPC(tcpConn * TcpConnection, msg *msgpacket.MSG_RPC, msgBody proto.Message) {
	var msgRes proto.Message = nil
	switch t:= msgBody.(type) {
	case *msgpacket.MSG_TEST:
		{
			msgRes = pthis.processRPCTest(tcpConn, t)
		}
	}

	msgRPCRes := &msgpacket.MSG_RPC_RES{
		MsgId:msg.MsgId,
		MsgType:msg.MsgType,
		ResCode:msgpacket.RESPONSE_CODE_RESPONSE_CODE_NONE,
	}

	if msgRes != nil {
		var err error
		msgRPCRes.MsgBin, err = proto.Marshal(msgRes)
		if err != nil {
			log.LogErr(err)
		}
	}
	tcpConn.TcpConnectSendProtoMsg(msgpacket.MSG_TYPE__MSG_RPC_RES, msgRPCRes)
}
func (pthis*Server)processRPCRes(tcpConn * TcpConnection, msg *msgpacket.MSG_RPC_RES, msgBody proto.Message) {
	defer func() {
		err := recover()
		if err != nil {
			log.LogErr(err)
		}
	}()
	if pthis.rpcMgr == nil {
		return
	}
	rreq := pthis.rpcMgr.RPCManagerFindReq(msg.MsgId)
	if rreq == nil {
		return
	}
	if rreq.chNtf != nil {
		rreq.chNtf <- msgBody
	}
}

func (pthis*Server)processRPCTest(tcpDial * TcpConnection, msg *msgpacket.MSG_TEST) *msgpacket.MSG_TEST_RES {
	log.LogDebug(msg)
	return &msgpacket.MSG_TEST_RES{Id: msg.Id}
}

// SendRPC_Async @brief will block timeoutMilliSec
func (pthis*Server)SendRPC_Async(msgType msgpacket.MSG_TYPE, protoMsg proto.Message, timeoutMilliSec int) proto.Message {

	defer func() {
		err := recover()
		if err != nil {
			log.LogErr(err)
		}
	}()

	if pthis.rpcMgr == nil || pthis.connDial == nil{
		return nil
	}

	msgRPC := msgpacket.MSG_RPC{
		MsgId:lin_common.GenUUID64_V4(),
		MsgType:int32(msgType),
		MsgBin:ProtoPacketToBin(msgType, protoMsg),
	}

	rreq := pthis.rpcMgr.RPCManagerAddReq(msgRPC.MsgId)

	pthis.connDial.TcpConnectSendProtoMsg(msgpacket.MSG_TYPE__MSG_RPC, &msgRPC)

	var res proto.Message = nil
	select{
	case resCh := <-rreq.chNtf:
		log.LogDebug(resCh)
		res, _ = resCh.(proto.Message)
	case <-time.After(time.Millisecond * time.Duration(timeoutMilliSec)):
	}

	pthis.rpcMgr.RPCManagerDelReq(msgRPC.MsgId)

	return res
}
