package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"go/ast"
	"go/parser"
	"go/token"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"lin/log"
	. "lin/msgpacket"
	"path/filepath"
	"runtime"
	"strings"
)

type MAP_PARSE_FUNC map[int32]func(binMsg[]byte) proto.Message
var mapVirtualTable = make(MAP_PARSE_FUNC)

func isMsgType(str string) (b bool, msgName string, msgType string) {
	if strings.Index(str, "MSG_TYPE_") == 0 {
		return true, str[len("MSG_TYPE_"):len(str)], "msgpacket." + str[len("MSG_TYPE__"):len(str)]
	}
	return false, "", ""
}

func genAllMsgParse() {
	_,filename,_,_ := runtime.Caller(0)
	pathBase := filepath.Dir(filename) + "/.."
	msgprotoPath := pathBase + "/msgpacket"

	var fset = token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, msgprotoPath, nil, 0)
	if err != nil {
		log.LogErr(err)
		return
	}
	for _, pkg := range pkgs {
		if pkg.Name != "msgpacket" {
			continue
		}
		for keyFileName, file := range pkg.Files {

			str := filepath.Base(keyFileName)
			if str != "msg.pb.go" {
				continue
			}

			for _, valobj := range file.Scope.Objects {
				if valobj.Kind != ast.Con {
					continue
				}
				b, msgType, msgName := isMsgType(valobj.Name)
				log.LogDebug("name:", msgName, " type:", msgType)
				if !b {
					continue
				}
				ProtoParseAddText(msgName, msgType)
			}

/*			ast.Inspect(file, func(n ast.Node) bool {
				if n != nil {
					fmt.Println("ast.node", n)
				}

				switch x := n.(type) {
				case *ast.TypeSpec:
					if _, ok := x.Type.(*ast.StructType); ok {
						fmt.Println(x.Name.Name)
					}
				}
				return true
			})*/
		}
	}
}

func InitMsgParseVirtualTable(){

	genAllMsgParse()

/*	mapVirtualTable[int32(MSG_TYPE__MSG_LOGIN)] = func (binMsg []byte)proto.Message {
		msg := &MSG_LOGIN{}
		proto.Unmarshal(binMsg, msg)
		return msg
	}*/

	//ProtoParseAddText("msg.MSG_TEST", "_MSG_TEST")
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