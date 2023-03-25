package main

import (
	"context"
	"github.com/golang/protobuf/proto"
	"go.mongodb.org/mongo-driver/bson"
	"goserver/common"
	"goserver/msgpacket"
	"goserver/server/server_common"
)

func (pthis*DBSrv)process_Go_CallBackMsg_PB_MSG_DBSERVER_READ(msg * msgpacket.PB_MSG_DBSERVER_READ) (msgType int32, protoMsg proto.Message) {
	msgRet := &msgpacket.PB_MSG_DBSERVER_READ_RES{}
	msgType = int32(msgpacket.PB_MSG_TYPE__PB_MSG_DBSERVER_READ_RES)
	protoMsg = msgRet

	// get db type by table name
	pbTbl, pbQueryKey, _, _ := pthis.dbDataTypeMgr.GetDBType(msg.DatabaseAppName, msg.TableName, DB_DATATYPE_QUERY_BIT)
	proto.Unmarshal(msg.Key, pbQueryKey)
	bsonKey := server_common.PBToBsonD(pbQueryKey, 10)
	common.LogDebug(bsonKey)

	// read from db
	dbClient := pthis.dbMgr.GetDBConnection(msg.DatabaseAppName)
	cur, err := dbClient.MongoClient.Database(dbClient.DB).Collection(msg.TableName).Find(context.TODO(), bsonKey)
	if err != nil {
		common.LogErr(err)
	}

	msgRet.DatabaseAppName = msg.DatabaseAppName
	msgRet.TableName = msg.TableName

	for ; cur.Next(context.TODO()) ; {
		bsonRes := bson.M{}
		cur.Decode(bsonRes)
		server_common.BsonMToPB(pbTbl, bsonRes, 10)
		common.LogDebug(pbTbl)
		bin, _ := proto.Marshal(pbTbl)
		msgRet.Record = append(msgRet.Record, bin)
	}

	return
}
