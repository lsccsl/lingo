package main

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"goserver/common"
	"goserver/server/server_common"
	"sync"

	"github.com/golang/protobuf/proto"

"google.golang.org/protobuf/reflect/protoregistry"
)

type DBDataTypeMgr struct {
	mapDBDataTypeLock sync.RWMutex
	mapDBDataType map[string]*dbDataType
	mapDBDef map[string]*DBDef
}

type DBDef struct {
	mapTable map[string]*DBTableDef
}

type DBTableDef struct {
	TableName      string
	TableProto     string
	QueryKeyProto  string
	UpdateKeyProto string
	DeleteKeyProto string
}

type DB_DATATYPE_GET uint
const (
	DB_DATATYPE_GET_TBL        DB_DATATYPE_GET = 0x01
	DB_DATATYPE_GET_QUERY_KEY  DB_DATATYPE_GET = 0x02
	DB_DATATYPE_GET_UPDATE_KEY DB_DATATYPE_GET = 0x04
	DB_DATATYPE_GET_DELETE_KEY DB_DATATYPE_GET = 0x08

	DB_DATATYPE_QUERY_BIT = DB_DATATYPE_GET_TBL | DB_DATATYPE_GET_QUERY_KEY
	DB_DATATYPE_UPDATE_BIT = DB_DATATYPE_GET_TBL | DB_DATATYPE_GET_UPDATE_KEY
)

type dbDataType struct {
	fullName string
	msgRef protoreflect.MessageType
}

func ConstructDBDataTypeMgr() *DBDataTypeMgr {
	dbTypeMgr := &DBDataTypeMgr{
		mapDBDataType:make(map[string]*dbDataType),
		mapDBDef:make(map[string]*DBDef),
	}

	dbCfg := server_common.GetAllDataBaseCfg()
	for _, dbDefCfg := range dbCfg{
		dbDef := &DBDef{
			mapTable : make(map[string]*DBTableDef),
		}
		dbTypeMgr.mapDBDef[dbDefCfg.DataBaseAppName] = dbDef
		for _, dbTblCfg := range dbDefCfg.Tables {
			dbTblDef := &DBTableDef{
				TableName:dbTblCfg.TableName,
				TableProto:dbTblCfg.TableProto,
				QueryKeyProto:dbTblCfg.QueryKeyProto,
				UpdateKeyProto:dbTblCfg.UpdateKeyProto,
				DeleteKeyProto:dbTblCfg.DeleteKeyProto,
			}
			dbDef.mapTable[dbTblCfg.TableName] = dbTblDef

			dbTypeMgr.getDataTypeIns(dbTblDef.TableProto)
			dbTypeMgr.getDataTypeIns(dbTblDef.QueryKeyProto)
			if 0 != len(dbTblDef.UpdateKeyProto) {
				dbTypeMgr.getDataTypeIns(dbTblDef.UpdateKeyProto)
			} else {
				dbTblDef.UpdateKeyProto = dbTblDef.QueryKeyProto
			}
			if 0!= len(dbTblDef.DeleteKeyProto) {
				dbTypeMgr.getDataTypeIns(dbTblDef.DeleteKeyProto)
			} else {
				dbTblDef.DeleteKeyProto = dbTblDef.QueryKeyProto
			}
		}
	}

	return dbTypeMgr
}

func (pthis*DBDataTypeMgr)getDataTypeInsNoAdd(name string) proto.Message {
	dbType, _ := pthis.mapDBDataType[name]
	if nil != dbType {
		msgIns := proto.MessageV1(dbType.msgRef.New())
		return msgIns
	}
	return nil
}

func (pthis*DBDataTypeMgr)getDataTypeIns(name string) proto.Message {
	msgIns := pthis.getDataTypeInsNoAdd(name)
	if nil != msgIns {
		return msgIns
	}

	msgFullName := protoreflect.FullName(name)
	msgRef, err := protoregistry.GlobalTypes.FindMessageByName(msgFullName)
	if nil == msgRef{
		common.LogErr(err, " msgFullName:", msgFullName)
		return nil
	}

	pthis.mapDBDataType[name] = &dbDataType{
		fullName:string(msgFullName),
		msgRef:msgRef,
	}

	msgIns = proto.MessageV1(msgRef.New())
	common.LogInfo("get db type:", name)
	return msgIns
}

func (pthis*DBDataTypeMgr)GetDBType(dbName string, tblName string, getFlag DB_DATATYPE_GET) (tbl proto.Message,
	queryKey proto.Message,
	updateKey proto.Message,
	deleteKey proto.Message) {

	tbl = nil
	queryKey = nil
	updateKey = nil
	deleteKey = nil

	pthis.mapDBDataTypeLock.Lock()
	defer pthis.mapDBDataTypeLock.Unlock()

	dbDef, _ := pthis.mapDBDef[dbName]
	if nil == dbDef {
		return
	}

	dbTbl, _ := dbDef.mapTable[tblName]
	if nil == dbTbl {
		return
	}

	if getFlag & DB_DATATYPE_GET_TBL != 0 {
		tbl = pthis.getDataTypeInsNoAdd(dbTbl.TableProto)
	}
	if getFlag & DB_DATATYPE_GET_QUERY_KEY != 0 {
		queryKey = pthis.getDataTypeInsNoAdd(dbTbl.QueryKeyProto)
	}
	if getFlag & DB_DATATYPE_GET_UPDATE_KEY != 0 {
		queryKey = pthis.getDataTypeInsNoAdd(dbTbl.UpdateKeyProto)
	}
	if getFlag & DB_DATATYPE_GET_DELETE_KEY != 0 {
		queryKey = pthis.getDataTypeInsNoAdd(dbTbl.DeleteKeyProto)
	}

	return
}

