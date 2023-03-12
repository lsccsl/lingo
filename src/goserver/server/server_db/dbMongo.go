package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"goserver/common"
)

type DBMongo struct {
	mongoClient * mongo.Client
}

func ConstructorDBMongo(DBUser string,
	DBPwd string,
	DBIp string,
	DBPort int) *DBMongo {
	db := &DBMongo{}

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d", DBUser, DBPwd, DBIp, DBPort)

	clientOptions := options.Client().ApplyURI(uri)
	var err error
	db.mongoClient, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		common.LogErr("connect to mongo db err,", uri, " err", err)
	}
	err = db.mongoClient.Ping(context.TODO(), nil)
	if err != nil {
		common.LogErr("connect to mongo db err,", uri, " err", err)
	} else {
		common.LogInfo("connect to mongo db suc,", uri, " err:", err)
	}

	return db
}
