package main

import "github.com/golang/protobuf/proto"

type Client struct {
	tcpConn *TcpConnection
	clientID int64
}

func ConstructClient(tcpConn *TcpConnection,clientID int64) *Client {
	c := &Client{
		tcpConn:tcpConn,
		clientID:clientID,
	}
	return c
}

func (pthis*Client) ClientGetConnection()*TcpConnection{
	return pthis.tcpConn
}
func (pthis*Client) ClientGetClientID()int64{
	return pthis.clientID
}

func (pthis*Client) ClientProcess(protoMsg proto.Message) {
}