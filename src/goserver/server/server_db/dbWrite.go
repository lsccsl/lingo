package main

import (
	"context"
	"errors"
	"github.com/golang/protobuf/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"goserver/common"
	"goserver/msgpacket"
	"goserver/server/server_common"
)

func (pthis*DBSrv)process_Go_CallBackMsg_PB_MSG_DBSERVER_WRITE(msg * msgpacket.PB_MSG_DBSERVER_WRITE) (msgType int32, protoMsg proto.Message) {
	msgRet := &msgpacket.PB_MSG_DBSERVER_WRITE_RES{}
	msgType = int32(msgpacket.PB_MSG_TYPE__PB_MSG_DBSERVER_WRITE_RES)
	protoMsg = msgRet

	// get db type by table name
	pbTbl, pbUpdateKey, _, _ := pthis.dbDataTypeMgr.GetDBType(msg.DatabaseAppName, msg.TableName, DB_DATATYPE_QUERY_BIT)
	proto.Unmarshal(msg.Key, pbUpdateKey)
	bsonKey := server_common.PBToBsonD(pbUpdateKey, 10)
	common.LogDebug(bsonKey)

	proto.Unmarshal(msg.Record, pbTbl)
	bsonRecord := server_common.PBToBsonM(pbTbl, 10)
	common.LogDebug(bsonRecord)

	dbClient := pthis.dbMgr.GetDBConnection(msg.DatabaseAppName)
	col := dbClient.MongoClient.Database(dbClient.DB).Collection(msg.TableName)
	sRes := col.FindOne(context.TODO(), bsonKey)
	if sRes.Err() != nil {
		if errors.Is(sRes.Err(), mongo.ErrNoDocuments) {
			common.LogDebug(bsonRecord)
			col.InsertOne(context.TODO(), bsonRecord)
			return
		} else {
			common.LogErr(sRes.Err())
		}
	}

	//col.Indexes().List()

	common.LogDebug(bsonRecord)
	col.UpdateOne(context.TODO(), bsonKey, bsonRecord)

	return
}

