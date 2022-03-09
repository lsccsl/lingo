package main

/*import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	. "lin/msgpacket"
)

type MAP_PARSE_FUNC map[int32]func(binMsg[]byte) proto.Message
var mapVirtualTable = make(MAP_PARSE_FUNC)


func InitMsgParseVirtualTable(){
	PBParseAddText("msg.MSG_LOGIN_RES", "_MSG_LOGIN_RES")
	PBParseAddText("msg.MSG_TEST", "_MSG_TEST")
	PBParseAddText("msg.MSG_TEST_RES", "_MSG_TEST_RES")
}

func ParseProtoMsg(binMsg []byte, msgType int32)proto.Message{
	if nil == mapVirtualTable{
		fmt.Println("parse table not init")
	}
	parsor, ok := mapVirtualTable[msgType]
	if ok && parsor != nil {
		return parsor(binMsg)
	}
	return PBParseByName(binMsg, msgType)
}



// ========================================= //
type PBMsgParse struct{
	msgType int32
	fullName string
	msgRef protoreflect.MessageType
}

var mapPBMsgParse map[int32]PBMsgParse

func PBParseAdd(name string, msgTye int32){
	if nil == mapPBMsgParse{
		mapPBMsgParse = make(map[int32]PBMsgParse)
	}

	_, ok := mapPBMsgParse[msgTye]
	if ok{
		return
	}

	msgFullName := protoreflect.FullName(name)
	msgType, err := protoregistry.GlobalTypes.FindMessageByName(msgFullName)
	if nil == msgType{
		fmt.Println(err)
		return
	}
	mapPBMsgParse[msgTye] = PBMsgParse{msgTye, name, msgType}
}
func PBParseAddText(name string, msgType string){
	intType, ok := MSG_TYPE_value[msgType]
	if !ok{
		fmt.Println("no msgtype:", msgType)
		return
	}
	PBParseAdd(name, intType)
}

func PBParseByName(binMsg []byte, msgType int32)proto.Message {
	parse, ok := mapPBMsgParse[msgType]
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
}*/