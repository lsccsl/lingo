package main

import (
	"github.com/golang/protobuf/proto"
	"lin/log"
	"lin/msg"
	"sync/atomic"
)

type interClientMsg struct {
	msgType msg.MSG_TYPE
	protoMsg proto.Message
}

type Client struct {
	tcpConn *TcpConnection
	clientID int64
	chClientMsg chan *interClientMsg
	isStopProcess int32
}

func ConstructClient(tcpConn *TcpConnection,clientID int64) *Client {
	c := &Client{
		tcpConn:tcpConn,
		clientID:clientID,
		chClientMsg:make(chan *interClientMsg),
		isStopProcess:0,
	}

	go c.go_clientProcess()

	return c
}

func (pthis*Client) go_clientProcess() {
	defer func() {
		err := recover()
		if err != nil {
			log.LogErr(err)
		}
	}()

	MSG_LOOP:
	for {
		select {
		case clientMsg := <- pthis.chClientMsg:
			if clientMsg == nil {
				break MSG_LOOP
			}
			pthis.processClientMsg(clientMsg)
		}
	}

	atomic.StoreInt32(&pthis.isStopProcess, 1)
	close(pthis.chClientMsg)
}

func (pthis*Client) ClientGetConnection()*TcpConnection{
	return pthis.tcpConn
}
func (pthis*Client) ClientGetClientID()int64{
	return pthis.clientID
}

func (pthis*Client) PushClientMsg(msgType msg.MSG_TYPE, protoMsg proto.Message) {
	if atomic.LoadInt32(&pthis.isStopProcess) == 1 {
		return
	}

	pthis.chClientMsg <- &interClientMsg{
		msgType:msgType,
		protoMsg:protoMsg,
	}
}

func (pthis*Client) processClientMsg (interMsg * interClientMsg) {
	switch t:=interMsg.protoMsg.(type){
	case *msg.MSG_TEST:
		log.LogDebug(t)
	}
}