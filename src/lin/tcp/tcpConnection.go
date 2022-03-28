package tcp

import (
	"bytes"
	"context"
	"lin/lin_common"
	"net"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

const G_MTU int = 1536
const MAX_PACK_LEN int = 65535

type TCP_CONNECTION_CLOSE_REASON int
const(
	TCP_CONNECTION_CLOSE_REASON_none     TCP_CONNECTION_CLOSE_REASON = 0
	TCP_CONNECTION_CLOSE_REASON_timeout  TCP_CONNECTION_CLOSE_REASON = 1
	TCP_CONNECTION_CLOSE_REASON_readerr  TCP_CONNECTION_CLOSE_REASON = 2
	TCP_CONNECTION_CLOSE_REASON_dialfail TCP_CONNECTION_CLOSE_REASON = 3
	TCP_CONNECTION_CLOSE_REASON_writeerr TCP_CONNECTION_CLOSE_REASON = 4
	TCP_CONNECTION_CLOSE_REASON_relogin TCP_CONNECTION_CLOSE_REASON = 5
)

type TCP_CONNECTIOON_TYPE int
const(
	TCP_CONNECTIOON_TYPE_client TCP_CONNECTIOON_TYPE = 1
	TCP_CONNECTIOON_TYPE_server TCP_CONNECTIOON_TYPE = 2
)

type TCP_CONNECTION_ID int64
type MAP_TCPCONN map[TCP_CONNECTION_ID]*TcpConnection

type InterfaceTcpConnection interface {
	CBReadProcess(tcpConn *TcpConnection, recvBuf * bytes.Buffer)(bytesProcess int)
	CBConnectAccept(tcpConn *TcpConnection, err error) // accept connection
	CBConnectDial(tcpConn *TcpConnection, err error)   // dial connection
	CBConnectClose(tcpConn *TcpConnection, closeReason TCP_CONNECTION_CLOSE_REASON)
}

type InterfaceConnManage interface {
	CBAddTcpConn(tcpConn *TcpConnection)
	CBGenConnectionID() TCP_CONNECTION_ID
	CBGetConnectionCB() InterfaceTcpConnection
	CBDelTcpConn(id TCP_CONNECTION_ID)
}

type interMsgTcpWrite struct {
	bin []byte
}

type TcpConnection struct {
	connectionID    TCP_CONNECTION_ID
	netConn         *net.TCPConn
	cbTcpConnection InterfaceTcpConnection
	chMsgWrite      chan *interMsgTcpWrite
	closeExpireSec  int
	connMgr         InterfaceConnManage

	canWrite int32

	IsAccept bool
	SrvID int64
	ClientID int64

	ConnType TCP_CONNECTIOON_TYPE
	ConnData interface{}

	// stats
	ByteRecv int64
	ByteSend int64
	SendCount int64
	ByteProc int64
	clsRsn   TCP_CONNECTION_CLOSE_REASON
}

func startTcpConnection(connMgr InterfaceConnManage, conn net.Conn, closeExpireSec int) (*TcpConnection, error) {

	tcpC, ok := conn.(*net.TCPConn)
	if !ok {
		return nil, lin_common.GenErr(lin_common.ERR_not_tcp_connection)
	}

	tcpConn := &TcpConnection{
		connectionID:    connMgr.CBGenConnectionID(),
		netConn:         tcpC,
		cbTcpConnection: connMgr.CBGetConnectionCB(),
		closeExpireSec:  closeExpireSec,
		connMgr:         connMgr,
		canWrite:        1,
		chMsgWrite:      make(chan*interMsgTcpWrite, 100),
		IsAccept:        true,
		ByteRecv:        0,
		ByteSend:        0,
		SendCount:       0,
		ByteProc:        0,
		clsRsn:          TCP_CONNECTION_CLOSE_REASON_none,
	}
	runtime.SetFinalizer(tcpConn, (*TcpConnection).TcpConnectClose)

	realTcpConn := conn.(*net.TCPConn)
	if realTcpConn != nil {
		realTcpConn.SetNoDelay(true)
		realTcpConn.SetReadBuffer(65535)
		realTcpConn.SetWriteBuffer(65535)
	}

	connMgr.CBAddTcpConn(tcpConn)
	if tcpConn.cbTcpConnection != nil {
		func(){
			defer func() {
				err := recover()
				if err != nil {
					lin_common.LogErr(err)
				}
			}()
			tcpConn.cbTcpConnection.CBConnectAccept(tcpConn, nil)
		}()
	}
	go tcpConn.go_tcpConnRead()
	go tcpConn.go_tcpConnWrite()

	return tcpConn, nil
}

func startTcpDial(connMgr InterfaceConnManage, SrvID int64, ip string, port int,
	closeExpireSec int, dialTimeoutSec int, redialCount int,
	ctx context.Context) (*TcpConnection, error) {
	defer func() {
		err := recover()
		if err != nil {
			lin_common.LogErr(err)
		}
	}()

	tcpConn := &TcpConnection{
		SrvID:           SrvID,
		connectionID:    connMgr.CBGenConnectionID(),
		netConn:         nil,
		cbTcpConnection: connMgr.CBGetConnectionCB(),
		closeExpireSec:  closeExpireSec,
		connMgr:         connMgr,
		canWrite:        1,
		chMsgWrite:      make(chan*interMsgTcpWrite, 100),
		IsAccept:        false,
		ByteRecv:        0,
		ByteSend:        0,
		SendCount:       0,
		ByteProc:        0,
		clsRsn:          TCP_CONNECTION_CLOSE_REASON_none,
	}
	runtime.SetFinalizer(tcpConn, (*TcpConnection).TcpConnectClose)

	addr := ip + ":" + strconv.Itoa(port)

	if dialTimeoutSec > 0 {
		go func() {
			defer func() {
				err := recover()
				if err != nil {
					lin_common.LogErr(err)
				}
			}()

			var err error
			var conn net.Conn
			DIAL_LOOP:
			for i := 0; i < redialCount; i ++ {
				tBegin := time.Now()
				lin_common.LogDebug(tcpConn.connectionID, " begin dial:", addr, " connection id:", tcpConn.connectionID)

				d := net.Dialer{Timeout: time.Second * time.Duration(dialTimeoutSec)}
				conn, err = d.DialContext(ctx, "tcp", addr)
				lin_common.LogDebug(tcpConn.connectionID, " end dial err:", err, " connection id:", tcpConn.connectionID)
				//conn, err = net.DialTimeout("tcp", addr, time.Second * time.Duration(dialTimeoutSec))

				tEnd := time.Now()
				if err != nil || conn == nil {
					lin_common.LogDebug(tcpConn.connectionID, " will retry ", i, " ", redialCount, " ", tcpConn.netConn, " ", err)
					if strings.Index(err.Error(), "operation was canceled") >= 0 {
						break DIAL_LOOP
					}
					interval := int64(dialTimeoutSec) - (tEnd.Unix() - tBegin.Unix())
					runtime.Gosched()
					if interval <= 0 {
						interval = 0
					}
					time.Sleep(time.Second * time.Duration(interval + 1))
					continue
				}
				tcpConn.netConn = conn.(*net.TCPConn)
				break
			}

			if err != nil || conn == nil{
				lin_common.LogDebug("fail ", err, conn, " connection id:", tcpConn.connectionID)
				if tcpConn.cbTcpConnection != nil {
					tcpConn.cbTcpConnection.CBConnectClose(tcpConn, TCP_CONNECTION_CLOSE_REASON_dialfail)
				}
				if tcpConn.netConn != nil {
					tcpConn.TcpConnectClose()
				}
				return
			}

			connMgr.CBAddTcpConn(tcpConn)
			if tcpConn.cbTcpConnection != nil {
				tcpConn.cbTcpConnection.CBConnectDial(tcpConn, nil)
			}
			go tcpConn.go_tcpConnRead()
			go tcpConn.go_tcpConnWrite()
		}()
		return tcpConn, nil
	} else {
		var err error
		dialTimeoutSec = -dialTimeoutSec
		if dialTimeoutSec < 1 {
			dialTimeoutSec = 1
		}
		for i:=0 ; i < redialCount; i ++ {
			con, err := net.DialTimeout("tcp", addr, time.Second * time.Duration(dialTimeoutSec))
			tcpConn.netConn = con.(*net.TCPConn)
			if err != nil {
				lin_common.LogErr("will retry ", i, " ", redialCount, " ", err)
				continue
			}
		}
		if err != nil {
			lin_common.LogErr("fail ", err)
/*			if tcpConn.cbTcpConnection != nil {
				tcpConn.cbTcpConnection.CBConnectClose(tcpConn)
			}*/
			return nil, err
		}

		connMgr.CBAddTcpConn(tcpConn)
		if tcpConn.cbTcpConnection != nil {
			tcpConn.cbTcpConnection.CBConnectDial(tcpConn, nil)
		}
		go tcpConn.go_tcpConnRead()
		go tcpConn.go_tcpConnWrite()

		return tcpConn, nil
	}
}

func (pthis *TcpConnection)go_tcpConnRead() {

	if pthis.netConn == nil {
		return
	}

	/*var TimerConnClose * time.Timer = nil*/
	defer func() {

		err := recover()
		if err != nil {
			lin_common.LogErr(err)
		}
	}()

	TmpBuf := make([]byte, G_MTU)
	recvBuf := bytes.NewBuffer(make([]byte, 0, MAX_PACK_LEN))

/*	expireInterval := time.Second * time.Duration(pthis.closeExpireSec)
	if pthis.closeExpireSec > 0 {
		TimerConnClose = time.AfterFunc(expireInterval, func() {
			pthis.TcpConnectSetCloseReason(TCP_CONNECTION_CLOSE_REASON_timeout)
			lin_common.LogDebug("time out close tcp connection:", pthis.connectionID, " srvid:", pthis.SrvID, " clientid:", pthis.ClientID,
				" expire sec:", pthis.closeExpireSec)
			pthis.TcpConnectClose()
		})
	}*/

	READ_LOOP:
	for {
		pthis.netConn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(pthis.closeExpireSec)))
		readSize, err := pthis.netConn.Read(TmpBuf)
		if err != nil {
			lin_common.LogDebug(pthis.connectionID, " clientid:", pthis.ClientID, " srvid:", pthis.SrvID, " err:", err)
			switch t := err.(type) {
			case net.Error:
				{
					if t.Timeout(){
						pthis.TcpConnectSetCloseReason(TCP_CONNECTION_CLOSE_REASON_timeout)
					} else if t.Temporary() {
						lin_common.LogDebug("temporary:", t)
						continue
					} else {
						//lin_common.LogDebug("other err:", t)
					}
				}
			case *net.OpError:
				lin_common.LogDebug(t)
			default:
				//lin_common.LogDebug(t)
			}

			pthis.TcpConnectSetCloseReason(TCP_CONNECTION_CLOSE_REASON_readerr)
			break READ_LOOP
		}
		pthis.ByteRecv += int64(readSize)

/*		if TimerConnClose != nil {
			//log.LogDebug("reset close timeout:", pthis.connectionID, " srvid:", pthis.SrvID, " clientid:", pthis.ClientID, " expire:", pthis.closeExpireSec)
			TimerConnClose.Reset(expireInterval)
		}*/

		recvBuf.Write(TmpBuf[0:readSize])

		if pthis.cbTcpConnection == nil {
			recvBuf.Next(readSize)
			continue
		}

		func(){
			defer func() {
				err := recover()
				if err != nil {
					lin_common.LogErr(err)
				}
			}()
			PROCESS_LOOP:
			for recvBuf.Len() > 0 {
				bytesProcess := pthis.cbTcpConnection.CBReadProcess(pthis, recvBuf)
				if bytesProcess <= 0 {
					break PROCESS_LOOP
				}
				pthis.ByteProc += int64(bytesProcess)
				recvBuf.Next(bytesProcess)
			}
		}()
	}

	pthis.netConn.Close()

	pthis.quitTcpWrite()
/*	if TimerConnClose != nil {
		TimerConnClose.Stop()
	}*/
	pthis.cbTcpConnection.CBConnectClose(pthis, pthis.clsRsn)
	if pthis.connMgr != nil {
		pthis.connMgr.CBDelTcpConn(pthis.connectionID)
	}
}

func (pthis *TcpConnection)go_tcpConnWrite() {
	defer func() {
		err := recover()
		if err != nil {
			lin_common.LogErr(err)
		}
	}()

	WRITE_LOOP:
	for {
		select {
		case tcpW := <- pthis.chMsgWrite:
			if tcpW == nil {
				break WRITE_LOOP
			}
			//todo: option get more data and combine write to tcp channel
			writeSZ, err := pthis.netConn.Write(tcpW.bin)
			if err != nil {
				lin_common.LogDebug(" write tcp err:", err)
				pthis.TcpConnectSetCloseReason(TCP_CONNECTION_CLOSE_REASON_writeerr)
				pthis.netConn.Close()
				break WRITE_LOOP
			}
			if writeSZ < len(tcpW.bin) {
				lin_common.LogErr("write sz wrong:", writeSZ)
			}
			pthis.ByteSend += int64(writeSZ)
			pthis.SendCount ++
		}
	}

	atomic.StoreInt32(&pthis.canWrite, 0)
	close(pthis.chMsgWrite)
}


func (pthis *TcpConnection)TcpConnectSendBin(bin []byte) {
	if atomic.LoadInt32(&pthis.canWrite) == 0 {
		return
	}

	tcpW := &interMsgTcpWrite{
		make([]byte,len(bin)),
	}
	copy(tcpW.bin, bin)
	pthis.chMsgWrite <- tcpW
	//tcpW.bin = append(tcpW.bin, bin...)
}


func (pthis *TcpConnection)TcpGetConn() net.Conn {
	return pthis.netConn
}

func (pthis *TcpConnection)TcpConnectClose() {
	runtime.SetFinalizer(pthis, nil)
	//lin_common.LogDebug(" close:", pthis.TcpConnectionID(), " client id:", pthis.ClientID, " srv id:", pthis.SrvID)
	if pthis.netConn != nil {
		pthis.netConn.Close()
	}
	pthis.quitTcpWrite()
}

func (pthis *TcpConnection)quitTcpWrite() {
	defer func() {
		err := recover()
		if err != nil {
			lin_common.LogErr(err)
		}
	}()

	if atomic.LoadInt32(&pthis.canWrite) != 0 {
		pthis.chMsgWrite <- nil
	}
}

func (pthis *TcpConnection)TcpConnectionID() TCP_CONNECTION_ID {
	return pthis.connectionID
}

func (pthis *TcpConnection)TcpConnectSetCloseReason(closeReason TCP_CONNECTION_CLOSE_REASON) {
	if TCP_CONNECTION_CLOSE_REASON_none == pthis.clsRsn {
		pthis.clsRsn = closeReason
	}
}