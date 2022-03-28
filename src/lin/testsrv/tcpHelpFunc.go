package main

import (
	"bytes"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"lin/lin_common"
	"lin/msgpacket"
	"net"
	"time"
)

const PACK_HEAD_SIZE int = 6

func recvProtoMsg(tcpConn *net.TCPConn, timeOutSend int, retryCount int) (proto.Message, error) {
	defer func() {
		tcpConn.SetReadDeadline(time.Time{})
	}()

	recvBuf := bytes.NewBuffer(make([]byte, 0, MAX_PACK_LEN))

	readSize := 0
	var errR error
	for i := 0; i < retryCount && readSize < PACK_HEAD_SIZE; i ++ {
		bin := make([]byte, PACK_HEAD_SIZE - readSize)
		tcpConn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(timeOutSend)))
		var rz int = 0
		rz, errR = tcpConn.Read(bin)
		if errR != nil {
			//lin_common.LogDebug(err)
			continue
		}
		readSize += rz
		recvBuf.Write(bin)
		if readSize >= PACK_HEAD_SIZE {
			break
		}
	}

	if readSize < PACK_HEAD_SIZE {
		//lin_common.LogDebug("read msg head err:", readSize, " err:", errR)
		return nil, lin_common.GenErr(lin_common.ERR_NONE, " err:", errR.Error())
	}

	binHead := recvBuf.Bytes()[0:PACK_HEAD_SIZE]
	packLen := int(binary.LittleEndian.Uint32(binHead[0:4]))
	packType := binary.LittleEndian.Uint16(binHead[4:6])
	recvBuf.Next(PACK_HEAD_SIZE)

	packLen = packLen - PACK_HEAD_SIZE
	if packLen > 0 {
		readSize = 0
		for j := 0; j < retryCount && readSize < packLen; j ++ {
			bin:= make([]byte, packLen - readSize)
			tcpConn.SetReadDeadline(time.Now().Add(time.Second * time.Duration(timeOutSend)))
			var rz int = 0
			rz, errR = tcpConn.Read(bin)
			if errR != nil {
				//lin_common.LogDebug(errR)
				continue
			}
			readSize += rz
			recvBuf.Write(bin)
			if readSize >= packLen {
				break
			}
		}
		if readSize < packLen {
			lin_common.LogDebug("read msg body err:", readSize, " err:", errR)
			return nil, lin_common.GenErr(lin_common.ERR_NONE, " err:", errR.Error())
		}
	} else {
		packLen = 0
	}

	binBody := recvBuf.Bytes()[0:packLen]
	msg := msgpacket.ParseProtoMsg(binBody, int32(packType))
	return msg, nil
}
