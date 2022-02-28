package main

import "github.com/golang/protobuf/proto"

type Client struct {
	tcpConnID TCP_CONNECTION_ID
	clientID int64
}

func ConstructClient(tcpConnID TCP_CONNECTION_ID,clientID int64) *Client {
	c := &Client{
		tcpConnID:tcpConnID,
		clientID:clientID,
	}
	return c
}

func (pthis*Client) ClientGetConnectID()TCP_CONNECTION_ID{
	return pthis.tcpConnID
}
func (pthis*Client) ClientGetClientID()int64{
	return pthis.clientID
}

func (pthis*Client) ClientProcess(protoMsg proto.Message) {
}