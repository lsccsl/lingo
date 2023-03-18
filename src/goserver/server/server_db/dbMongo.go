package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"goserver/common"
)

type DBMongo struct {
	MongoClient * mongo.Client
	DB string
}


func ConstructorDBMongo(DBUser string,
	DBPwd string,
	DBIp string,
	DBPort int,
	DB string) *DBMongo {
	db := &DBMongo{
		DB:DB,
	}

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d", DBUser, DBPwd, DBIp, DBPort)

	clientOptions := options.Client().ApplyURI(uri)
	var err error
	db.MongoClient, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		common.LogErr("connect to mongo db err,", uri, " err", err)
	}
	err = db.MongoClient.Ping(context.TODO(), nil)
	if err != nil {
		common.LogErr("connect to mongo db err,", uri, " err", err)
	} else {
		common.LogInfo("connect to mongo db suc,", uri, " err:", err)
	}

	return db
}
