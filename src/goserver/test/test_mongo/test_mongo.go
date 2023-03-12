
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/encoding/protojson"
	"goserver/common"
	"goserver/msgpacket"
	"log"
	"time"
)

type StudentOut struct {
	Id              string            `json:"id" bson:"_id"`
	Name string
	Age int
	Detail StudentDetail
}

type StudentOutHex struct {
	Id              primitive.ObjectID            `json:"_id" bson:"_id"`
	Name string
	Age int
	Detail StudentDetail
}


type StudentDetail struct {
	From string
	Count int
}
type StudentIn struct {
	Name string
	Age int
	Detail StudentDetail
}

func getAll(collection *mongo.Collection) {
	doc := bson.D{}
	cur, err := collection.Find(context.TODO(), doc)
	fmt.Println(cur, err)
	curCur := cur
	for {
		stuTmp := &StudentOut{}
		curCur.Decode(stuTmp)
		fmt.Println("read all:", stuTmp)
		if !curCur.Next(context.TODO()) {
			break
		}
	}
}

func test_pb_ref(client * mongo.Client) {
	collection := client.Database("test_db").Collection("test_user_collection")

	{
		delRes, err := collection.DeleteOne(context.TODO(), bson.M{"user_id":"123"})
		fmt.Println(delRes, err)
		delRes, err = collection.DeleteOne(context.TODO(), bson.M{"user_id":123})
		fmt.Println(delRes, err)
	}

	userMap := bson.M{}
	userMap["user_id"] = int64(123)

	userDetail := bson.M{}
	userDetail["detail_data"] = "detail data"
	userDetail["detail_id"] = int32(123456)
	userMap["detail"] = userDetail
	userMap["en_test"] = msgpacket.EN_TEST_EN_TEST2

	{
		var vlist []interface{}
		{
			mapRepeated := bson.M{}
			mapRepeated["repeated_str"] = "aaa repeated"
			mapRepeated["repeated_int"] = 667

			testMap := bson.M{}
			{
				testMapElem := bson.M{}
				testMapElem["map_str"] = "map string 1"
				testMapElem["map_int"] = 878

				testMap["1"] = testMapElem
			}
			{
				testMapElem := bson.M{}
				testMapElem["map_str"] = "map string 2"
				testMapElem["map_int"] = 888

				testMap["2"] = testMapElem
			}
			mapRepeated["test_map"] = testMap

			vlist = append(vlist, mapRepeated)
		}

		{
			mapRepeated := bson.M{}
			mapRepeated["repeated_str"] = "bbb repeated"
			mapRepeated["repeated_int"] = 668

			testMapElem := bson.M{}
			testMapElem["map_str"] = "map string 3"
			testMapElem["map_int"] = 889

			testMap := bson.M{}
			testMap["1"] = testMapElem
			mapRepeated["test_map"] = testMap

			vlist = append(vlist, mapRepeated)
		}

		userMap["test_repeated"] = vlist

	}

	res, _ := collection.InsertOne(context.TODO(), userMap)
	idHexMap, _ := res.InsertedID.(primitive.ObjectID)
	fmt.Println(idHexMap)
	idString := idHexMap.Hex()
	fmt.Println("insert bson map,idString:", idString)

	// read from db
	var msg proto.Message
	{
		sres := collection.FindOne(context.TODO(), bson.M{"_id": idHexMap})
		tmpMap := bson.M{}
		sres.Decode(tmpMap)
		fmt.Println("find one, decode as map:", tmpMap)

		msgParse := PBMsgGen("msgpacket.DBUserMain", tmpMap, 10)
		msg = msgParse.(proto.Message)
		fmt.Println("proto msg:", msg)
	}

	// write proto to db
	{
		var idHexMapNew primitive.ObjectID
		{
			dbMsg := msg.(*msgpacket.DBUserMain)
			dbMsg.Detail.DetailData = "new detail data"
			dbMsg.EnTest = msgpacket.EN_TEST_EN_TEST3
			dbMsg.TestRepeated[0].TestMap[1].MapStr = "update map str 11"
			idHexMapNew, _ = primitive.ObjectIDFromHex(dbMsg.XId)
		}
		bsonWrite := PBToBson(msg, 10)
		delete(bsonWrite, "_id")
		fmt.Println(bsonWrite)
		collection.UpdateOne(context.TODO(), bson.M{"_id": idHexMapNew}, bson.D{{"$set",bsonWrite}})
	}
}

func test_protocal(client * mongo.Client) {
	collection := client.Database("test_db").Collection("test_user_collection")

	{
		delRes, err := collection.DeleteOne(context.TODO(), bson.M{"user_id":"123"})
		fmt.Println(delRes, err)
		delRes, err = collection.DeleteOne(context.TODO(), bson.M{"user_id":123})
		fmt.Println(delRes, err)
	}

	userMap := bson.M{}
	userMap["user_id"] = int64(123)

	userDetail := bson.M{}
	userDetail["detail_data"] = "detail data"
	userDetail["detail_id"] = int32(123456)
	userMap["detail"] = userDetail

	res, _ := collection.InsertOne(context.TODO(), userMap)
	idHexMap, _ := res.InsertedID.(primitive.ObjectID)
	fmt.Println(idHexMap)
	idString := idHexMap.Hex()
	fmt.Println("insert bson map,idString:", idString)

	msg := &msgpacket.DBUserMain{}
	{
		sres := collection.FindOne(context.TODO(), bson.M{"_id": idHexMap})
		tmpMap := bson.M{}
		sres.Decode(tmpMap)
		fmt.Println("find one, decode as map:", tmpMap)

		jsonByte, _ := json.Marshal(tmpMap)
		fmt.Println("~~~~json byte", string(jsonByte))

		protojson.Unmarshal(jsonByte, msg)
		fmt.Println("proto msg:", msg.String())
	}

	//msg.UserId = int64(6789)
	msg.Detail.DetailData = "abbbb"
	{
		ms := protojson.MarshalOptions{UseProtoNames:true}
		jsonByte, _ := ms.Marshal(msg)
		fmt.Println("~~~~proto json byte", string(jsonByte))
		bsonMap := bson.M{}
		bson.UnmarshalExtJSON(jsonByte, true, bsonMap)
		delete(bsonMap, "_id")
		fmt.Println("bsonMap:", bsonMap)
		//bson.D{}
		uRes , err := collection.UpdateOne(context.TODO(), bson.M{"_id": idHexMap}/*bson.D{{"_id", idHexMap}}*/, bson.D{{"$set",bsonMap}})
		fmt.Println(uRes, err)
		//bsonMap["user_id"] = 9090
		//uRes , err = collection.UpdateByID(context.TODO(), idHexMap, bson.D{{"$set",bsonMap}})
		//fmt.Println(uRes, err)
	}

	{
		msg1 := &msgpacket.DBUserMain{}
		sres := collection.FindOne(context.TODO(), bson.M{"_id": idHexMap})
		tmpMap := bson.M{}
		sres.Decode(tmpMap)
		fmt.Println("find one, decode as map:", tmpMap)

		jsonByte, _ := json.Marshal(tmpMap)
		fmt.Println("~~~~json byte", string(jsonByte))

		protojson.Unmarshal(jsonByte, msg1)
		fmt.Println("proto msg:", msg1.String())
	}
}

func test_read_cfg() {
	ReadDBCfg("D:\\mywork\\git\\lingo\\cfg\\dbcfg.yml")
}

func main()  {
	common.InitLog("./test.log", "./test_err.log", true, true, true)

	//test_read_cfg()
	// 设置客户端选项
	clientOptions := options.Client().ApplyURI("mongodb://admin:123456@192.168.0.103:27017")
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

	test_pb_ref(client)
	return
	test_protocal(client)

	stu:=StudentIn{Name:"李四",
		Age:int(time.Now().Unix()),
		Detail: StudentDetail{"from",1},
	}

	var idHex primitive.ObjectID
	collection := client.Database("test_db").Collection("test_collection")
	{
		res, err := collection.InsertOne(context.TODO(), stu)
		fmt.Println(res, err)
		idHex, _ = res.InsertedID.(primitive.ObjectID)
		fmt.Println(idHex)
		idString := idHex.Hex()
		fmt.Println("idString:", idString)
	}


	getAll(collection)

	{
		sres := collection.FindOne(context.TODO(), bson.M{"_id":idHex})
		stuMap := bson.M{}
		sres.Decode(stuMap)
		fmt.Println("find one, decode as map:", stuMap)

		jsonByte, _ := json.Marshal(stuMap)
		fmt.Println("~~~~json byte", string(jsonByte))

		outStu := &StudentOutHex{}
		json.Unmarshal(jsonByte, outStu)

		fmt.Println("~~~parse json", outStu)
	}

	{
		sres := collection.FindOne(context.TODO(), bson.M{"_id":idHex})
		stuTmp := &StudentOutHex{}
		sres.Decode(stuTmp)
		fmt.Println("find one:", stuTmp)
	}

	{
		stuMap := bson.M{}
		stuMap["mapname"] = "mapname"
		stuMap["age"] = 1
		stuMap["detail"] = &StudentDetail{"from", 1}
		res, _ := collection.InsertOne(context.TODO(), stuMap)
		idHexMap, _ := res.InsertedID.(primitive.ObjectID)
		fmt.Println(idHexMap)
		idString := idHexMap.Hex()
		fmt.Println("insert bson map,idString:", idString)
	}


	//collection.DeleteOne(context.TODO(), stuTmp)
	collection.DeleteOne(context.TODO(), bson.M{"_id":idHex})

	getAll(collection)
}


