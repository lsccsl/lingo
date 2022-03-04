package main

import (
	"bytes"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"lin/log"
	"lin/msg"
	"net"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"
)

const G_MTU int = 1536
const MAX_PACK_LEN int = 65535

type TCP_CONNECTION_ID int64
type MAP_TCPCONN map[TCP_CONNECTION_ID]*TcpConnection

type InterfaceTcpConnection interface {
	CBReadProcess(tcpConn * TcpConnection, recvBuf * bytes.Buffer)(bytesProcess int)
	CBConnectAccept(tcpConn * TcpConnection, err error) // accept connection
	CBConnectDial(tcpConn * TcpConnection, err error) // dial connection
	CBConnectClose(tcpConn * TcpConnection)
}

type InterfaceConnManage interface {
	CBAddTcpConn(tcpConn *TcpConnection)
	CBGenConnectionID()TCP_CONNECTION_ID
	CBGetConnectionCB()InterfaceTcpConnection
	CBDelTcpConn(id TCP_CONNECTION_ID)
}

type interMsgTcpWrite struct {
	bin []byte
}

type TcpConnection struct {
	connectionID TCP_CONNECTION_ID
	netConn net.Conn
	cbTcpConnection InterfaceTcpConnection
	chMsgWrite chan *interMsgTcpWrite
	closeExpireSec int
	connMgr InterfaceConnManage

	canWrite int32

	IsAccept bool
	SrvID int64
	ClientID int64

	ConnType int64
	ConnData interface{}
}

func startTcpConnection(connMgr InterfaceConnManage, conn net.Conn, closeExpireSec int) (*TcpConnection, error) {

	tcpConn := &TcpConnection{
		connectionID:connMgr.CBGenConnectionID(),
		netConn:conn,
		cbTcpConnection:connMgr.CBGetConnectionCB(),
		closeExpireSec:closeExpireSec,
		connMgr:connMgr,
		canWrite:0,
		chMsgWrite:make(chan*interMsgTcpWrite),
		IsAccept:true,
	}

	go tcpConn.go_tcpConnRead()
	go tcpConn.go_tcpConnWrite()
	connMgr.CBAddTcpConn(tcpConn)

	if tcpConn.cbTcpConnection != nil {
		func(){
			defer func() {
				err := recover()
				if err != nil {
					log.LogErr(err)
				}
			}()
			tcpConn.cbTcpConnection.CBConnectAccept(tcpConn, nil)
		}()
	}

	return tcpConn, nil
}

func startTcpDial(connMgr InterfaceConnManage, SrvID int64, ip string, port int, closeExpireSec int, dialTimeoutSec int, redialCount int) (*TcpConnection, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.LogErr(err)
		}
	}()

	tcpConn := &TcpConnection{
		SrvID:SrvID,
		connectionID:connMgr.CBGenConnectionID(),
		netConn:nil,
		cbTcpConnection:connMgr.CBGetConnectionCB(),
		closeExpireSec:closeExpireSec,
		connMgr:connMgr,
		canWrite:0,
		chMsgWrite:make(chan*interMsgTcpWrite),
		IsAccept:false,
	}
	addr := ip + ":" + strconv.Itoa(port)

	if dialTimeoutSec > 0 {
		go func() {
			defer func() {
				err := recover()
				if err != nil {
					log.LogErr(err)
				}
			}()

			var err error
			for i := 0; i < redialCount; i ++ {
				tBegin := time.Now()
				tcpConn.netConn, err = net.DialTimeout("tcp", addr, time.Second * time.Duration(dialTimeoutSec))
				tEnd := time.Now()
				if err != nil {
					log.LogErr("will retry ", i, " ", redialCount, " ", tcpConn.netConn, " ", err)
					interval := int64(dialTimeoutSec) - (tEnd.Unix() - tBegin.Unix())
					runtime.Gosched()
					if interval <= 0 {
						interval = 0
					}
					time.Sleep(time.Second * time.Duration(interval + 1))
					continue
				}
				break
			}

			if err != nil {
				log.LogErr("fail ", err)
				if tcpConn.cbTcpConnection != nil {
					tcpConn.cbTcpConnection.CBConnectClose(tcpConn)
				}
				return
			}

			go tcpConn.go_tcpConnRead()
			go tcpConn.go_tcpConnWrite()
			connMgr.CBAddTcpConn(tcpConn)

			if tcpConn.cbTcpConnection != nil {
				tcpConn.cbTcpConnection.CBConnectDial(tcpConn, nil)
			}
		}()
		return tcpConn, nil
	} else {
		var err error
		dialTimeoutSec = -dialTimeoutSec
		if dialTimeoutSec < 1 {
			dialTimeoutSec = 1
		}
		for i:=0 ; i < redialCount; i ++ {
			tcpConn.netConn, err = net.DialTimeout("tcp", addr, time.Second * time.Duration(dialTimeoutSec))
			if err != nil {
				log.LogErr("will retry ", i, " ", redialCount, " ", err)
				continue
			}
		}
		if err != nil {
			log.LogErr("fail ", err)
/*			if tcpConn.cbTcpConnection != nil {
				tcpConn.cbTcpConnection.CBConnectClose(tcpConn)
			}
*/			return nil, err
		}

		go tcpConn.go_tcpConnRead()
		go tcpConn.go_tcpConnWrite()
		connMgr.CBAddTcpConn(tcpConn)

		if tcpConn.cbTcpConnection != nil {
			tcpConn.cbTcpConnection.CBConnectDial(tcpConn, nil)
		}

		return tcpConn, nil
	}

	return tcpConn, nil
}

func (pthis * TcpConnection)go_tcpConnRead() {
	var TimerConnClose * time.Timer = nil
	defer func() {
		if TimerConnClose != nil {
			TimerConnClose.Stop()
		}
		pthis.quitTcpWrite()
		pthis.cbTcpConnection.CBConnectClose(pthis)
		if pthis.connMgr != nil {
			pthis.connMgr.CBDelTcpConn(pthis.connectionID)
		}

		err := recover()
		if err != nil {
			log.LogErr(err)
		}
	}()

	TmpBuf := make([]byte, G_MTU)
	recvBuf := bytes.NewBuffer(make([]byte, 0, MAX_PACK_LEN))

	expireInterval := time.Second * time.Duration(pthis.closeExpireSec)
	if pthis.closeExpireSec > 0 {
		TimerConnClose = time.AfterFunc(expireInterval, func() {
			log.LogDebug("time out close tcp connection:", pthis.connectionID, " srvid:", pthis.SrvID, " clientid:", pthis.ClientID)
			pthis.TcpConnectClose()
		})
	}

	READ_LOOP:
	for {
		readSize, err := pthis.netConn.Read(TmpBuf)
		if err != nil {
			break READ_LOOP
		}

		if TimerConnClose != nil {
			TimerConnClose.Reset(expireInterval)
		}

		recvBuf.Write(TmpBuf[0:readSize])

		if pthis.cbTcpConnection == nil {
			recvBuf.Next(readSize)
			continue
		}

		bytesProcess := 0
		func(){
			defer func() {
				err := recover()
				if err != nil {
					log.LogErr(err)
				}
			}()
			bytesProcess = pthis.cbTcpConnection.CBReadProcess(pthis, recvBuf)
		}()

		if bytesProcess < 0 {
			break READ_LOOP
		} else if bytesProcess > 0 {
			recvBuf.Next(bytesProcess)
		}
	}
}

func (pthis * TcpConnection)go_tcpConnWrite() {
	defer func() {
		err := recover()
		if err != nil {
			log.LogErr(err)
		}
	}()

	WRITE_LOOP:
	for {
		select {
		case tcpW := <- pthis.chMsgWrite:
			if tcpW == nil {
				break WRITE_LOOP
			}
			//todo: option wait for more data and combine write to tcp channel
			pthis.netConn.Write(tcpW.bin)
		}
	}

	atomic.StoreInt32(&pthis.canWrite, 1)
	close(pthis.chMsgWrite)
}

func (pthis * TcpConnection)TcpConnectWriteBin(bin []byte) {
	if atomic.LoadInt32(&pthis.canWrite) != 0 {
		return
	}

	tcpW := &interMsgTcpWrite{
		make([]byte,len(bin)),
	}
	copy(tcpW.bin, bin)
	//fmt.Println(&tcpW.bin[0], &bin[0], ret)
	pthis.chMsgWrite <- tcpW
	//tcpW.bin = append(tcpW.bin, bin...)
}
func (pthis*TcpConnection)TcpConnectWriteProtoMsg(msgType msg.MSG_TYPE, protoMsg proto.Message) {
	binMsg, _ := proto.Marshal(protoMsg)
	var wb []byte
	var buf bytes.Buffer
	_ = binary.Write(&buf,binary.LittleEndian,uint32(6 + len(binMsg)))
	_ = binary.Write(&buf,binary.LittleEndian,uint16(msgType))
	wb = buf.Bytes()
	wb = append(wb, binMsg...)

	pthis.TcpConnectWriteBin(wb)
}

func (pthis * TcpConnection)TcpGetConn() net.Conn {
	return pthis.netConn
}

func (pthis * TcpConnection)TcpConnectClose() {
	pthis.quitTcpWrite()
	pthis.netConn.Close()
}

func (pthis * TcpConnection)quitTcpWrite() {
	if atomic.LoadInt32(&pthis.canWrite) != 0 {
		pthis.chMsgWrite <- nil
	}
}

func (pthis * TcpConnection)TcpConnectionID() TCP_CONNECTION_ID {
	return pthis.connectionID
}
