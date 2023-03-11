package main

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"goserver/common"

	"github.com/golang/protobuf/proto"

"google.golang.org/protobuf/reflect/protoregistry"
)

type DBDataTypeMgr struct {
	mapDBDataType map[string]*dbDataType
}

type dbDataType struct {
	fullName string
	msgRef protoreflect.MessageType
}

func ConstructDBDataTypeMgr() *DBDataTypeMgr {
	dbTypeMgr := &DBDataTypeMgr{
		mapDBDataType:make(map[string]*dbDataType),
	}

	return dbTypeMgr
}

func (pthis*DBDataTypeMgr)GetDataTypeIns(name string) proto.Message {
	dbType, _ := pthis.mapDBDataType[name]
	if nil != dbType {
		msgIns := proto.MessageV1(dbType.msgRef.New())
		return msgIns
	}

	msgFullName := protoreflect.FullName(name)
	msgRef, err := protoregistry.GlobalTypes.FindMessageByName(msgFullName)
	if nil == msgRef{
		common.LogInfo(err, " msgFullName:", msgFullName)
		return nil
	}

	pthis.mapDBDataType[name] = &dbDataType{
		fullName:string(msgFullName),
		msgRef:msgRef,
	}

	msgIns := proto.MessageV1(msgRef.New())
	return msgIns
}
