
package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type StudentOut struct {
	Id              string            `json:"id" bson:"_id"`
	Name string
	Age int
}

type StudentOutHex struct {
	Id              primitive.ObjectID            `json:"id" bson:"_id"`
	Name string
	Age int
}


type StudentIn struct {
	Name string
	Age int
}

func getAll(collection *mongo.Collection) {
	doc := bson.D{}
	cur, err := collection.Find(context.TODO(), doc)
	fmt.Println(cur, err)
	curCur := cur
	for {
		stuTmp := &StudentOut{}
		curCur.Decode(stuTmp)
		fmt.Println("read:", stuTmp)
		if !curCur.Next(context.TODO()) {
			break
		}
	}
}

func main()  {
	// 设置客户端选项
	clientOptions := options.Client().ApplyURI("mongodb://admin:123456@192.168.0.105:27017")
	// 连接 MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	stu:=StudentIn{Name:"李四",Age:int(time.Now().Unix())}

	collection := client.Database("test_db").Collection("test_collection")
	res, err := collection.InsertOne(context.TODO(), stu)
	fmt.Println(res, err)
	idHex, ok := res.InsertedID.(primitive.ObjectID)
	fmt.Println(idHex, ok)
	idString := idHex.Hex()
	fmt.Println("idString:", idString)

	getAll(collection)

	sres := collection.FindOne(context.TODO(), bson.M{"_id":idHex})
	stuTmp := &StudentOutHex{}
	sres.Decode(stuTmp)
	fmt.Println("find one:", stuTmp)

	//collection.DeleteOne(context.TODO(), stuTmp)
	collection.DeleteOne(context.TODO(), bson.M{"_id":idHex})

	getAll(collection)
}


