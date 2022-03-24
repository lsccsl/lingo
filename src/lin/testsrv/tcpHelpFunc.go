package main

import (
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"lin/msgpacket"
	"net"
)

const PACK_HEAD_SIZE int = 6

func recvProtoMsg(tcpConn net.Conn) (proto.Message, error) {
	binHead := make([]byte, PACK_HEAD_SIZE)
	_, err := tcpConn.Read(binHead)
	if err != nil {
		return nil, err
	}

	packLen := binary.LittleEndian.Uint32(binHead[0:4])
	packType := binary.LittleEndian.Uint16(binHead[4:6])

	binBody:= make([]byte, packLen)
	_, err = tcpConn.Read(binBody)
	if err != nil {
		return nil, err
	}

	msg := msgpacket.ParseProtoMsg(binBody, int32(packType))
	return msg, nil
}
