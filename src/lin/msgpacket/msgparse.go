package msgpacket

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"go/ast"
	"go/parser"
	"go/token"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"lin/lin_common"
	"path/filepath"
	"runtime"
	"strings"
)

type MAP_PARSE_FUNC map[int32]func(binMsg[]byte) proto.Message
var mapVirtualTable = make(MAP_PARSE_FUNC)

type MAP_MSGNAME_MSGTYPE map[string]int32
var mapMsgNameType = make(MAP_MSGNAME_MSGTYPE)

func isMsgType(str string) (b bool, msgName string, msgType string) {
	if strings.Index(str, "MSG_TYPE_") == 0 {
		return true, str[len("MSG_TYPE_"):len(str)], "msgpacket." + str[len("MSG_TYPE__"):len(str)]
	}
	if strings.Index(str, "PB_MSG_INTER_TYPE_") == 0 {
		return true, str[len("PB_MSG_INTER_TYPE_"):len(str)], "msgpacket." + str[len("PB_MSG_INTER_TYPE__"):len(str)]
	}
	return false, "", ""
}

func addMsgNameType(msgName string, msgType string) {
	intType, ok := MSG_TYPE_value[msgType]
	if ok{
		mapMsgNameType[msgName] = intType
		return
	}
	intType, ok = PB_MSG_INTER_TYPE_value[msgType]
	if ok{
		mapMsgNameType[msgName] = intType
		return
	}
	fmt.Println("no msgtype:", msgType, " msg name", msgName)
	return
}

func GetMsgTypeByMsgInstance(msg proto.Message) int32 {
	msgRef := proto.MessageReflect(msg)
	if msgRef == nil {
		return 0
	}
	des := msgRef.Descriptor()
	if des == nil {
		return 0
	}
	name := string(des.FullName())

	intType, _ := mapMsgNameType[name]
	return intType
}

func genAllMsgParse(msgprotoPath string) {
	if 0 == len(msgprotoPath) {
		_,filename,_,_ := runtime.Caller(0)
		pathBase := filepath.Dir(filename) + "/.."
		msgprotoPath = pathBase + "/msgpacket"
	}

	var fset = token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, msgprotoPath, nil, 0)
	if err != nil {
		lin_common.LogErr(err)
		return
	}
	for _, pkg := range pkgs {
		if pkg.Name != "msgpacket" {
			continue
		}
		for keyFileName, file := range pkg.Files {

			str := filepath.Base(keyFileName)
			if str != "msg.pb.go" && str != "msginter.pb.go" {
				continue
			}

			for _, valobj := range file.Scope.Objects {
				if valobj.Kind != ast.Con {
					continue
				}
				b, msgType, msgName := isMsgType(valobj.Name)
				lin_common.LogDebug("name:", msgName, " type:", msgType)
				if !b {
					continue
				}
				ProtoParseAddText(msgName, msgType)
				addMsgNameType(msgName, msgType)
			}
		}
	}
}

func InitMsgParseVirtualTable(msgprotoPath string){

	genAllMsgParse(msgprotoPath)

/*	msg := &MSG_TEST{}
	fmt.Println(GetMsgTypeByMsgInstance(msg))
*/
/*	mapVirtualTable[int32(MSG_TYPE__MSG_LOGIN)] = func (binMsg []byte)proto.Message {
		msg := &MSG_LOGIN{}
		proto.Unmarshal(binMsg, msg)
		return msg
	}*/

	//ProtoParseAddText("msg.MSG_TEST", "_MSG_TEST")
}

func ParseProtoMsg(binMsg []byte, msgType int32) proto.Message {
	msg := ProtoParseByName(binMsg, msgType)
	if msg != nil {
		return msg
	}

	if nil == mapVirtualTable {
		fmt.Println("parse table not init")
	}
	parsor, ok := mapVirtualTable[msgType]
	if ok && parsor != nil {
		return parsor(binMsg)
	}
	return nil
}



// ========================================= //
type ProtoMsgParse struct{
	msgType int32
	fullName string
	msgRef protoreflect.MessageType
}

var mapProtoMsgParse map[int32]ProtoMsgParse

func ProtoParseAdd(name string, msgTye int32){
	if nil == mapProtoMsgParse {
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
	if ok{
		ProtoParseAdd(name, intType)
		return
	}

	intType, ok = PB_MSG_INTER_TYPE_value[msgType]
	if ok {
		ProtoParseAdd(name, intType)
		return
	}
	fmt.Println("no msg name", name, " msg type", msgType)
	return
}

func ProtoParseByName(binMsg []byte, msgType int32)proto.Message {
	parse, ok := mapProtoMsgParse[msgType]
	if !ok{
		lin_common.LogDebug("can't find msg type:", msgType)
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

func ProtoPacketToBin(msgType MSG_TYPE, protoMsg proto.Message) []byte {
	binMsg, _ := proto.Marshal(protoMsg)
	var wb []byte
	var buf bytes.Buffer
	_ = binary.Write(&buf,binary.LittleEndian,uint32(6 + len(binMsg)))
	_ = binary.Write(&buf,binary.LittleEndian,uint16(msgType))
	wb = buf.Bytes()
	wb = append(wb, binMsg...)

	return wb
}

func ProtoUnPacketFromBin(recvBuf * bytes.Buffer) (MSG_TYPE, int, proto.Message) {
	if recvBuf.Len() < 6 {
		return MSG_TYPE__MSG_NULL, 0, nil
	}
	binHead := recvBuf.Bytes()[0:6]

	packLen := binary.LittleEndian.Uint32(binHead[0:4])
	packType := binary.LittleEndian.Uint16(binHead[4:6])

	if recvBuf.Len() < int(packLen){
		return MSG_TYPE__MSG_NULL, 0, nil
	}

	binBody := recvBuf.Bytes()[6:packLen]

	return MSG_TYPE(packType), int(packLen), ParseProtoMsg(binBody, int32(packType))
}

