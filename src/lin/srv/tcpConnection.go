package main

import (
	"bytes"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"lin/log"
	"lin/msg"
	"net"
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
	CBConnect(tcpConn * TcpConnection, err error)
	CBConnectClose(id TCP_CONNECTION_ID)
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
	isAccept bool

	AppID int64
	AppType int64
	AppData interface{}
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
		isAccept:true,
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
			tcpConn.cbTcpConnection.CBConnect(tcpConn, nil)
		}()
	}

	return tcpConn, nil
}

func startTcpDial(connMgr InterfaceConnManage,	ip string, port int, closeExpireSec int, dialTimeoutSec int) (*TcpConnection, error) {
	defer func() {
		err := recover()
		if err != nil {
			log.LogErr(err)
		}
	}()

	tcpConn := &TcpConnection{
		connectionID:connMgr.CBGenConnectionID(),
		netConn:nil,
		cbTcpConnection:connMgr.CBGetConnectionCB(),
		closeExpireSec:closeExpireSec,
		connMgr:connMgr,
		canWrite:0,
		chMsgWrite:make(chan*interMsgTcpWrite),
		isAccept:false,
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
			tcpConn.netConn, err = net.DialTimeout("tcp", addr, time.Second * time.Duration(dialTimeoutSec))
			if err != nil {
				log.LogErr(err)
				if tcpConn.cbTcpConnection != nil {
					tcpConn.cbTcpConnection.CBConnect(nil, err)
				}
			} else {
				connMgr.CBAddTcpConn(tcpConn)
				if tcpConn.cbTcpConnection != nil {
					tcpConn.cbTcpConnection.CBConnect(tcpConn, nil)
				}
			}
		}()
	} else {
		var err error
		if dialTimeoutSec == 0 {
			tcpConn.netConn, err = net.Dial("tcp", addr)
		} else {
			tcpConn.netConn, err = net.DialTimeout("tcp", addr, time.Second * time.Duration(-dialTimeoutSec))
		}
		if err != nil {
			log.LogErr(err)
			if tcpConn.cbTcpConnection != nil {
				tcpConn.cbTcpConnection.CBConnect(nil, err)
			}
			return nil, err
		}
		if tcpConn.cbTcpConnection != nil {
			tcpConn.cbTcpConnection.CBConnect(tcpConn, nil)
		}

		connMgr.CBAddTcpConn(tcpConn)
		return tcpConn, nil
	}

	return tcpConn, nil
}

func (pthis * TcpConnection)go_tcpConnRead() {
	defer func() {
		pthis.cbTcpConnection.CBConnectClose(pthis.connectionID)
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
	var TimerConnClose * time.Timer = nil
	if pthis.closeExpireSec > 0 {
		TimerConnClose = time.AfterFunc(expireInterval, func() {
			log.LogDebug("time out close tcp connection", pthis.connectionID, pthis.AppID)
			pthis.TcpConnectClose()
		})
	}

	READ_LOOP:
	for {
		readSize, err := pthis.netConn.Read(TmpBuf)
		if err != nil {
			break READ_LOOP
		}

		if pthis.closeExpireSec > 0 {
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
	if atomic.LoadInt32(&pthis.canWrite) != 0 {
		pthis.chMsgWrite <- nil
	}
	pthis.netConn.Close()
}

func (pthis * TcpConnection)TcpConnectionID() TCP_CONNECTION_ID {
	return pthis.connectionID
}
