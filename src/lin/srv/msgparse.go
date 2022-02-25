package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	. "lin/msg"
)

type MAP_PARSE_FUNC map[int32]func(binMsg[]byte) proto.Message
var mapVirtualTable = make(MAP_PARSE_FUNC)

func InitMsgParseVirtualTable(){
	mapVirtualTable[int32(MSG_TYPE__MSG_LOGIN)] = func (binMsg []byte)proto.Message {
		msg := &MSG_LOGIN{}
		proto.Unmarshal(binMsg, msg)
		return msg
	}

	ProtoParseAddText("msg.MSG_TEST", "_MSG_TEST")
}

func ParseProtoMsg(binMsg []byte, msgType int32) proto.Message {
	if nil == mapVirtualTable{
		fmt.Println("parse table not init")
	}
	parsor, ok := mapVirtualTable[msgType]
	if ok && parsor != nil {
		return parsor(binMsg)
	}
	return ProtoParseByName(binMsg, msgType)
}



// ========================================= //
type ProtoMsgParse struct{
	msgType int32
	fullName string
	msgRef protoreflect.MessageType
}

var mapProtoMsgParse map[int32]ProtoMsgParse

func ProtoParseAdd(name string, msgTye int32){
	if nil == mapProtoMsgParse{
		mapProtoMsgParse = make(map[int32]ProtoMsgParse)
	}

	_, ok := mapProtoMsgParse[msgTye]
	if ok{
		return
	}

	msgFullName := protoreflect.FullName(name)
	msgType, err := protoregistry.GlobalTypes.FindMessageByName(msgFullName)
	if nil == msgType{
		fmt.Println(err)
		return
	}
	mapProtoMsgParse[msgTye] = ProtoMsgParse{msgTye, name, msgType}
}
func ProtoParseAddText(name string, msgType string){
	intType, ok := MSG_TYPE_value[msgType]
	if !ok{
		fmt.Println("no msgtype:", msgType)
		return
	}
	ProtoParseAdd(name, intType)
}

func ProtoParseByName(binMsg []byte, msgType int32)proto.Message {
	parse, ok := mapProtoMsgParse[msgType]
	if !ok{
		return nil
	}

	if nil == parse.msgRef{
		return nil
	}
	msgIns := proto.MessageV1(parse.msgRef.New())
	if nil == msgIns{
		return nil
	}
	proto.Unmarshal(binMsg, msgIns)
	return msgIns
}