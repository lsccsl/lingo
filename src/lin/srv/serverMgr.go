package main

import (
	"bytes"
	"fmt"
	"github.com/golang/protobuf/proto"
	"lin/lin_common"
	cor_pool "lin/lin_cor_pool"
	"lin/msgpacket"
	"lin/tcp"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const(
	EN_CORPOOL_JOBTYPE_Rpc_req = cor_pool.EN_CORPOOL_JOBTYPE_user + 1
	EN_CORPOOL_JOBTYPE_client_Rpc_req = cor_pool.EN_CORPOOL_JOBTYPE_user + 2
)

type MAP_CLIENT map[int64/*client id*/]*Client
type MAP_SERVER map[int64/*server id*/]*Server
type interProtoMsg struct {
	msgType msgpacket.MSG_TYPE
	protoMsg  proto.Message
	tcpConn *tcp.TcpConnection
}

type ClientMapMgr struct {
	mapClientMutex sync.Mutex
	mapClient MAP_CLIENT
}
type ServerMapMgr struct {
	mapServerMutex sync.Mutex
	mapServer MAP_SERVER
}
type ServerMgrStatic struct {
	totalPacket int64
	totalRecv int64
	totalSend int64
	totalProc int64
	timestamp float64
}
type ServerMgr struct {
	srvID int64
	ClientMapMgr
	ServerMapMgr
	tcpMgr *tcp.TcpMgr
	httpSrv *lin_common.HttpSrvMgr
	rpcPool *cor_pool.CorPool

	heartbeatIntervalSec int

	ServerMgrStatic
}

func (pthis*ServerMgr)CBReadProcess(tcpConn *tcp.TcpConnection, recvBuf * bytes.Buffer) (bytesProcess int) {

	packType, bytesProcess, protoMsg := msgpacket.ProtoUnPacketFromBin(recvBuf)

	if protoMsg == nil {
		//log.LogErr("can't parse msg:", tcpConn.ByteRecv, " proc:", tcpConn.ByteProc)
		return
	}

	switch packType {
	case msgpacket.MSG_TYPE__MSG_LOGIN:
		t, ok := protoMsg.(*msgpacket.MSG_LOGIN)
		if ok && t != nil {
			pthis.processClientLogin(t.Id, tcpConn)
		}

	case msgpacket.MSG_TYPE__MSG_SRV_REPORT:
		t, ok := protoMsg.(*msgpacket.MSG_SRV_REPORT)
		if ok && t != nil {
			if tcpConn.IsAccept {
				pthis.processSrvReport(tcpConn, t.SrvId)
			}
		}

	case msgpacket.MSG_TYPE__MSG_RPC:
		t, ok := protoMsg.(*msgpacket.MSG_RPC)
		if ok && t != nil {
			pthis.processRPCReq(tcpConn, t)
		}

	case msgpacket.MSG_TYPE__MSG_RPC_RES:
		t, ok := protoMsg.(*msgpacket.MSG_RPC_RES)
		if ok && t != nil {
			pthis.processRPCRes(tcpConn, t)
		}

	default:
		pthis.processMsg(tcpConn, packType, protoMsg)
	}

	return
}

func (pthis*ServerMgr)CBConnectAccept(tcpConn *tcp.TcpConnection, err error) {
	if err != nil {
		lin_common.LogErr(err)
	}
	if tcpConn == nil {
		return
	}
	//lin_common.LogDebug(tcpConn.TcpGetConn().LocalAddr(), tcpConn.TcpGetConn().RemoteAddr(), tcpConn.TcpConnectionID())
}
func (pthis*ServerMgr)CBConnectDial(tcpConn *tcp.TcpConnection, err error) {
	if err != nil {
		lin_common.LogErr(err)
	}
	if tcpConn == nil {
		return
	}
	lin_common.LogDebug(tcpConn.TcpGetConn().LocalAddr(), tcpConn.TcpGetConn().RemoteAddr(), tcpConn.TcpConnectionID())

	pthis.processDailConnect(tcpConn)
}

func (pthis*ServerMgr)CBConnectClose(tcpConn *tcp.TcpConnection, closeReason tcp.TCP_CONNECTION_CLOSE_REASON) {
	lin_common.LogDebug("id:", tcpConn.TcpConnectionID(),
		" srvid:", tcpConn.SrvID, " clientid:", tcpConn.ClientID, " is accept:", tcpConn.IsAccept,
		" closeReason:", closeReason)

	if !tcpConn.IsAccept {
		srv := pthis.getServer(tcpConn.SrvID)
		if srv != nil {
			srv.PushInterMsg(&interMsgConnClose{tcpConn})
		} else {
			lin_common.LogDebug("will check redial, can't find srv id:", tcpConn.SrvID)
			pthis.tcpMgr.TcpDialMgrCheckReDial(tcpConn.SrvID)
		}
		//pthis.delServer(tcpConn.SrvID)
	} else {
		if tcpConn.SrvID != 0 {
			srv := pthis.getServer(tcpConn.SrvID)
			if srv != nil {
				srv.PushInterMsg(&interMsgConnClose{tcpConn})
			}
			//pthis.delServer(tcpConn.SrvID)
		} else if tcpConn.ClientID != 0 {
			oldC := pthis.getClient(tcpConn.ClientID)
			if oldC != nil {
				if oldC.ClientGetConnectionID() == tcpConn.TcpConnectionID(){
					pthis.delClient(tcpConn.ClientID)
				}
			}
		}
	}
}

func ConstructServerMgr(srvID int64, heartbeatIntervalSec int, rpcPoolCount int) *ServerMgr {
	srvMgr := &ServerMgr{srvID: srvID}
	srvMgr.mapClient = make(MAP_CLIENT)
	srvMgr.mapServer = make(MAP_SERVER)
	srvMgr.heartbeatIntervalSec = heartbeatIntervalSec
	srvMgr.rpcPool = cor_pool.CorPoolInit(rpcPoolCount)

	return srvMgr
}


func (pthis*ServerMgr)getClient(clientID int64) *Client {
	pthis.ClientMapMgr.mapClientMutex.Lock()
	defer pthis.ClientMapMgr.mapClientMutex.Unlock()

	oldC, _ := pthis.ClientMapMgr.mapClient[clientID]
	return oldC
}
func (pthis*ServerMgr)addClient(c *Client) {
	pthis.ClientMapMgr.mapClientMutex.Lock()
	defer pthis.ClientMapMgr.mapClientMutex.Unlock()

	pthis.ClientMapMgr.mapClient[c.clientID] = c
}
func (pthis*ServerMgr)delClient(clientID int64) {
	pthis.ClientMapMgr.mapClientMutex.Lock()
	defer pthis.ClientMapMgr.mapClientMutex.Unlock()

	oldC, _ := pthis.ClientMapMgr.mapClient[clientID]
	if  oldC != nil {
		oldC.ClientClose()
	}
	delete(pthis.ClientMapMgr.mapClient, clientID)
}


func (pthis*ServerMgr)getServer(srvID int64) *Server {
	pthis.ServerMapMgr.mapServerMutex.Lock()
	defer pthis.ServerMapMgr.mapServerMutex.Unlock()

	oldS, _ := pthis.ServerMapMgr.mapServer[srvID]
	return oldS
}
func (pthis*ServerMgr)addServer(s *Server) {
	pthis.ServerMapMgr.mapServerMutex.Lock()
	defer pthis.ServerMapMgr.mapServerMutex.Unlock()

	pthis.ServerMapMgr.mapServer[s.srvID] = s
}
func (pthis*ServerMgr)delServer(srvID int64) {
	pthis.ServerMapMgr.mapServerMutex.Lock()
	defer pthis.ServerMapMgr.mapServerMutex.Unlock()

	oldS, _ := pthis.ServerMapMgr.mapServer[srvID]
	if oldS != nil {
		oldS.ServerClose()
	}
	delete(pthis.ServerMapMgr.mapServer, srvID)
}

func (pthis*ServerMgr)processClientLogin(clientID int64, tcpConn *tcp.TcpConnection) {
	if tcpConn == nil {
		return
	}

	tcpConn.ClientID = clientID

	oldC := pthis.getClient(clientID)
	if oldC != nil {
		if oldC.ClientGetConnectionID() != tcpConn.TcpConnectionID() {
			oldC.tcpConn.TcpConnectSetCloseReason(tcp.TCP_CONNECTION_CLOSE_REASON_relogin)
			pthis.delClient(clientID)

			c := ConstructClient(pthis, tcpConn, clientID)
			c.mapStaticMsgRecv = oldC.mapStaticMsgRecv
			pthis.addClient(c)
		}
	} else {
		c := ConstructClient(pthis, tcpConn, clientID)
		pthis.addClient(c)
	}

	msgRes := &msgpacket.MSG_LOGIN_RES{}
	msgRes.Id = clientID
	msgRes.ConnectId = int64(tcpConn.TcpConnectionID())
	tcpConn.TcpConnectSendBin(msgpacket.ProtoPacketToBin(msgpacket.MSG_TYPE__MSG_LOGIN_RES, msgRes))
}

func (pthis*ServerMgr)processMsg(tcpConn *tcp.TcpConnection, msgType msgpacket.MSG_TYPE, protoMsg proto.Message) {
	if tcpConn.SrvID != 0 {
		srv := pthis.getServer(tcpConn.SrvID)
		if srv != nil {
			srv.PushProtoMsg(msgType, protoMsg, tcpConn)
			return
		}
		pthis.tcpMgr.TcpMgrCloseConn(tcpConn.TcpConnectionID())
		return
	} else if tcpConn.ConnData != nil{
		cli := tcpConn.ConnData.(*Client)
		if cli != nil {
			if msgType == msgpacket.MSG_TYPE__MSG_TEST {
				msgTest := protoMsg.(*msgpacket.MSG_TEST)
				msgTest.TimestampArrive = time.Now().UnixMilli()
			}
			cli.PushProtoMsg(msgType, protoMsg)
			return
		}
		pthis.tcpMgr.TcpMgrCloseConn(tcpConn.TcpConnectionID())
		return
	}
}

func (pthis*ServerMgr)processSrvReport(tcpAccept *tcp.TcpConnection, srvID int64){
	tcpAccept.SrvID = srvID
	srv := pthis.getServer(srvID)
	if srv == nil {
		srv = ConstructServer(pthis, nil, tcpAccept, srvID, pthis.heartbeatIntervalSec)
		pthis.addServer(srv)
	}
	if srv != nil {
		srv.PushInterMsg(&interMsgSrvReport{tcpAccept})
	}
}

func (pthis*ServerMgr)processDailConnect(tcpDial *tcp.TcpConnection){
	srvID := tcpDial.SrvID
	srv := pthis.getServer(srvID)
	if srv == nil {
		srv = ConstructServer(pthis, tcpDial, nil, srvID, pthis.heartbeatIntervalSec)
		pthis.addServer(srv)
	}
	if srv != nil {
		srv.PushInterMsg(&interMsgConnDial{tcpDial})
	}
}

func (pthis*ServerMgr)processRPCReq(tcpConn *tcp.TcpConnection, msg *msgpacket.MSG_RPC) {
	msgRPC := msgpacket.ParseProtoMsg(msg.MsgBin, msg.MsgType)
	//lin_common.LogDebug(msgRPC)
	if tcpConn.SrvID != 0 {
		srv := pthis.getServer(tcpConn.SrvID)
		if srv != nil {
			err := pthis.rpcPool.CorPoolAddJob(&cor_pool.CorPoolJobData{
				JobType_ : EN_CORPOOL_JOBTYPE_Rpc_req,
				JobCB_   : func(jd cor_pool.CorPoolJobData){
					srv.Go_ProcessRPC(tcpConn, msg, msgRPC)
				},
			})
			if err != nil {
				lin_common.LogErr("put job err:", err, " srvid:", tcpConn.SrvID)
			}
		} else {
			lin_common.LogErr("can't find srv", tcpConn.SrvID)
			pthis.tcpMgr.TcpMgrCloseConn(tcpConn.TcpConnectionID())
			return
		}
	}
}
func (pthis*ServerMgr)processRPCRes(tcpConn *tcp.TcpConnection, msgRPC *msgpacket.MSG_RPC_RES) {
	msgBody := msgpacket.ParseProtoMsg(msgRPC.MsgBin, msgRPC.MsgType)
	if tcpConn.SrvID != 0 {
		srv := pthis.getServer(tcpConn.SrvID)
		if srv != nil {
			srv.processRPCRes(tcpConn, msgRPC, msgBody)
		}
	}
}

func (pthis*ServerMgr)SendRPC_Async(srvID int64, msgType msgpacket.MSG_TYPE, protoMsg proto.Message, timeoutMilliSec int) proto.Message {
	srv := pthis.getServer(srvID)
	if srv == nil {
		return nil
	}
	return srv.SendRPC_Async(msgType, protoMsg, timeoutMilliSec)
}

func TcpConnectSendProtoMsg(tcpConn *tcp.TcpConnection, msgType msgpacket.MSG_TYPE, protoMsg proto.Message) {
	tcpConn.TcpConnectSendBin(msgpacket.ProtoPacketToBin(msgType, protoMsg))
}

func (pthis*ServerMgr)Dump(bDtail bool) string {
	var str string
	str += "\r\nclient:\r\n"

	timestamp := float64(time.Now().UnixMilli())
	var totalPacket int64 = 0


	func(){
		pthis.ClientMapMgr.mapClientMutex.Lock()
		defer pthis.ClientMapMgr.mapClientMutex.Unlock()
		mapStatic := make(MAP_CLIENT_STATIC)
		for _, val := range pthis.ClientMapMgr.mapClient {
			if bDtail {
				str += fmt.Sprintf("\r\n client id:%v id:%v map:%v", val.clientID, val.tcpConnID, val.mapStaticMsgRecv)
			}
			for skey, sval := range val.mapStaticMsgRecv {
				mapStatic[skey] += sval
			}
		}
		for _, sval := range mapStatic {
			totalPacket += sval
		}
		str += fmt.Sprintf("static:%v", mapStatic)
		str += "\r\nclient count:" + strconv.Itoa(len(pthis.ClientMapMgr.mapClient)) + " totalPacket:" + strconv.FormatInt(totalPacket, 10)
	}()

	str += "\r\nserver:\r\n"
	func(){
		pthis.ServerMapMgr.mapServerMutex.Lock()
		defer pthis.ServerMapMgr.mapServerMutex.Unlock()
		/*if bDtail*/ {
			for _, val := range pthis.ServerMapMgr.mapServer {
				var connAcptID tcp.TCP_CONNECTION_ID
				var connDialID tcp.TCP_CONNECTION_ID
				if val.connAcpt != nil {
					connAcptID = val.connAcpt.TcpConnectionID()
				}
				if val.connDial != nil {
					connDialID = val.connDial.TcpConnectionID()
				}
				totalRPCPacket := atomic.LoadInt64(&val.totalRPCPacket)
				tnow := float64(time.Now().UnixMilli())
				tRPCdiff := (tnow - val.timestamp)/float64(1000)
				diffRPCTotal := totalRPCPacket - val.totalRPCPacketLast
				aver := float64(diffRPCTotal) / tRPCdiff
				str += fmt.Sprintf("\r\n server id:%v acpt:%v dial:%v totalPakcet:%v diffRPCTotal:%v aver:%v tdiff:%v",
					val.srvID, connAcptID, connDialID, totalRPCPacket, diffRPCTotal, aver, tRPCdiff)
				val.timestamp = tnow
				val.totalRPCPacketLast = totalRPCPacket
			}
		}
		str += "\r\nserver count:" + strconv.Itoa(len(pthis.ServerMapMgr.mapServer))
	}()

	strTcp, totalRecv, totalSend, totalProc := pthis.tcpMgr.TcpMgrDump(bDtail)
	str += strTcp

	diffTotal := totalPacket - pthis.totalPacket
	diffRecv := totalRecv - pthis.totalRecv
	diffSend := totalSend - pthis.totalSend
	diffProc := totalProc - pthis.totalProc
	tdiff := (timestamp - pthis.timestamp) / float64(1000)
	if tdiff <= 0 {
		tdiff = 1
	}
	str += fmt.Sprintf("\r\n diffTotal:%v diffRecv:%v diffSend:%v diffProc:%v", diffTotal, diffRecv, diffSend, diffProc)
	str += fmt.Sprintf("\r\n Total ps:%v Recv ps:%v Send ps:%v Proc ps:%v tdiff:%v",
		float64(diffTotal) / tdiff, float64(diffRecv) / tdiff, float64(diffSend) / tdiff, float64(diffProc) / tdiff, tdiff)

	pthis.timestamp = float64(time.Now().UnixMilli())

	pthis.totalPacket = totalPacket
	pthis.totalRecv = totalRecv
	pthis.totalSend = totalSend
	pthis.totalProc = totalProc

	return str
}
