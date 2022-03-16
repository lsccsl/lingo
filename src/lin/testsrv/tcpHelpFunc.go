package main

import (
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"lin/msgpacket"
	"net"
)

const PACK_HEAD_SIZE int = 6

func recvProtoMsg(tcpConn net.Conn) proto.Message {
	binHead := make([]byte, PACK_HEAD_SIZE)
	readSize, err := tcpConn.Read(binHead)
	fmt.Println(readSize, err)

	packLen := binary.LittleEndian.Uint32(binHead[0:4])
	packType := binary.LittleEndian.Uint16(binHead[4:6])

	binBody:= make([]byte, packLen)
	readSize, err = tcpConn.Read(binBody)
	fmt.Println(readSize, err)

	return msgpacket.ParseProtoMsg(binBody, int32(packType))
}
