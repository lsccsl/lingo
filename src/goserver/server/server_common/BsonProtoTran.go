package server_common

import (
	"github.com/golang/protobuf/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/runtime/protoiface"
	"goserver/common"
	"strconv"
)

func PBGenFieldValue(kind protoreflect.Kind, des protoreflect.FieldDescriptor, v interface{}, recursiveCount int) (vRet protoreflect.Value) {
	switch kind {
	case protoreflect.BoolKind:
		var b bool
		switch v.(type){
		case string:
			iv, _ := strconv.Atoi(v.(string))
			b = (iv != 0)
			break
		case int:
			b = v.(int)!= 0
			break
		default:
			common.LogErr("kind:", kind, " des:", des, " val:", v)
		}
		vRet = protoreflect.ValueOfBool(b)
		break

	case protoreflect.Int32Kind:
		var val int
		switch v.(type){
		case string:
			val, _ = strconv.Atoi(v.(string))
			break
		case int32:
			val = int(v.(int32))
			break
		case int:
			val = v.(int)
			break
		case int64:
			val = int(v.(int64))
			break
		default:
			common.LogErr("kind:", kind, " des:", des, " val:", v)
		}
		vRet = protoreflect.ValueOfInt32(int32(val))
		break

	case protoreflect.StringKind:
		switch t := v.(type) {
		case string:
			vRet = protoreflect.ValueOfString(v.(string))
		case primitive.ObjectID:
			vRet = protoreflect.ValueOfString(t.Hex())
		default:
			common.LogErr("kind:", kind, " des:", des, " val:", v)
		}
		break

	case protoreflect.MessageKind:
		var msg proto.Message
		switch t := v.(type){
		case proto.Message:
			msg, _ = v.(proto.Message)
			vRet = protoreflect.ValueOfMessage(proto.MessageReflect(msg))
		case bson.M:
			mv1 := BsonToPBByName(string(des.Message().FullName()), t, recursiveCount - 1).(protoiface.MessageV1)
			vRet = protoreflect.ValueOfMessage(proto.MessageReflect(mv1))
		default:
			common.LogErr("kind:", kind, " des:", des, " val:", v)
		}
		break

	case protoreflect.EnumKind:
		switch t := v.(type) {
		case string:
			val, _ := strconv.Atoi(t)
			vRet = protoreflect.ValueOfEnum(protoreflect.EnumNumber(int32(val)))
		case int32:
			vRet = protoreflect.ValueOfEnum(protoreflect.EnumNumber(t))
		default:
			common.LogErr("kind:", kind, " des:", des, " val:", v)
		}
		break

	case protoreflect.Int64Kind:
		var val int64
		switch v.(type) {
		case string:
			val, _ = strconv.ParseInt(v.(string), 10, 64)
			break
		case int64:
			val = v.(int64)
			break
		case int:
			val = int64(v.(int))
			break
		case int32:
			val = int64(v.(int32))
			break

		default:
			common.LogErr("kind:", kind, " des:", des, " val:", v)
		}
		vRet = protoreflect.ValueOfInt64(val)
		break
	}
	return
}

func BsonMToPB(msgIns proto.Message, MapKV map[string]interface{}, recursiveCount int)interface{} {
	msgInsRef:= proto.MessageReflect(msgIns)
	if nil == msgInsRef{
		common.LogErr("no msg ref")
		return nil
	}

	MsgDesAll := msgInsRef.Descriptor()
	if nil == MsgDesAll{
		common.LogErr("no msg des all")
		return nil
	}
	MsgFields := MsgDesAll.Fields()
	if nil == MsgFields{
		common.LogErr("no msg fields")
		return nil
	}

	if nil == MapKV{
		common.LogErr("no map kv")
		return msgIns
	}

	for k,v := range MapKV{
		MsgDesField := MsgFields.ByTextName(k)
		if nil == MsgDesField{
			return nil
		}

		if MsgDesField.Cardinality() == protoreflect.Repeated{
			if MsgDesField.IsList(){
				// repeated
				muList := msgInsRef.Mutable(MsgDesField).List()
				var lst_input []interface{}
				switch t := v.(type) {
				case bson.A:
					{
						lst_input = t
						for _,lv := range lst_input{
							muList.Append(PBGenFieldValue(MsgDesField.Kind(), MsgDesField, lv, recursiveCount))
						}
					}

				default:
					common.LogErr(t)
				}
			} else  if MsgDesField.IsMap(){
				// map
				muMap := msgInsRef.Mutable(MsgDesField).Map()
				switch t := v.(type) {
				case bson.M:
					for km,vm := range t{
						tmpK := PBGenFieldValue(MsgDesField.MapKey().Kind(), MsgDesField.MapKey(), km, recursiveCount).MapKey()
						tmpV := PBGenFieldValue(MsgDesField.MapValue().Kind(), MsgDesField.MapValue(), vm, recursiveCount)
						muMap.Set(tmpK, tmpV)
					}

				default:
					common.LogErr(t)
				}
			}
		} else {
			msgInsRef.Set(MsgDesField, PBGenFieldValue(MsgDesField.Kind(), MsgDesField, v, recursiveCount))
		}
	}

	return msgIns
}
func BsonToPBByName(MsgName string, MapKV map[string]interface{}, recursiveCount int)interface{} {
	if (recursiveCount < 0) {
		common.LogErr("too many recursive")
		return nil
	}

	msgName := protoreflect.FullName(MsgName)
	msgType, err := protoregistry.GlobalTypes.FindMessageByName(msgName)
	if nil == msgType{
		common.LogErr(err)
		return nil
	}
	msgIns := proto.MessageV1(msgType.New())
	if nil == msgIns{
		common.LogErr("no msg ins")
		return nil
	}

	return BsonMToPB(msgIns, MapKV, recursiveCount)
}

func PBGetValue(field protoreflect.FieldDescriptor, val protoreflect.Value, bsonCon BsonContainer, recursiveCount int) interface{} {
	switch field.Kind() {
	case protoreflect.BoolKind:
		return val.Bool()
	case protoreflect.EnumKind:
		return int32(val.Enum())
	case protoreflect.Int32Kind:
		return val.Int()
	case protoreflect.Uint32Kind:
		return val.Uint()
	case protoreflect.Int64Kind:
		return val.Int()
	case protoreflect.Uint64Kind:
		return val.Uint()
	case protoreflect.FloatKind:
		return val.Float()
	case protoreflect.DoubleKind:
		return val.Float()
	case protoreflect.StringKind:
		return val.String()
	case protoreflect.BytesKind:
		return val.Bytes()
	case protoreflect.MessageKind:
		{
			pbMsg, ok := val.Message().Interface().(proto.Message)
			if !ok {
				common.LogErr("val:", val, " field:", field)
			} else {
				subContainer := bsonCon.GetSubContainer()
				PBToBson(pbMsg, subContainer, recursiveCount - 1)
				return subContainer.GetBsonContainer()
			}
		}

	case protoreflect.Sint64Kind:
		return val.Int()
	case protoreflect.Sint32Kind:
		return val.Int()
	case protoreflect.Sfixed32Kind:
		return val.Int()
	case protoreflect.Fixed32Kind:
		return val.Int()
	case protoreflect.Sfixed64Kind:
		return val.Int()
	case protoreflect.Fixed64Kind:
		return val.Int()
	default:
		common.LogErr("unknown protoc type:", field.Kind())
	}

	return nil
}

type BsonContainer interface {
	PutKV(key string, val interface{})
	GetBsonContainer()interface{}
	GetSubContainer()BsonContainer
}


// BsonContainerM bson.M container
type BsonContainerM struct {
	bsonMap bson.M
}

func (pthis*BsonContainerM)PutKV(key string, val interface{}) {
	pthis.bsonMap[key] = val
}

func (pthis*BsonContainerM)GetBsonContainer()interface{} {
	return pthis.bsonMap
}

func (pthis*BsonContainerM)GetSubContainer()BsonContainer {
	return &BsonContainerM{
		bsonMap: bson.M{},
	}
}

func PBToBsonM(msg proto.Message, recursiveCount int) bson.M {
	bsonM := &BsonContainerM{
		bsonMap: bson.M{},
	}

	PBToBson(msg, bsonM, recursiveCount)

	return bsonM.bsonMap
}


//
type BsonContainerD struct {
	bsonD bson.D
}
func (pthis*BsonContainerD)PutKV(key string, val interface{}) {
	pthis.bsonD = append(pthis.bsonD, bson.E{key, val})
}

func (pthis*BsonContainerD)GetBsonContainer()interface{} {
	return pthis.bsonD
}

func (pthis*BsonContainerD)GetSubContainer()BsonContainer {
	return &BsonContainerD{
		bsonD: bson.D{},
	}
}
func PBToBsonD(msg proto.Message, recursiveCount int) bson.D {
	bsonM := &BsonContainerD{
		bsonD: bson.D{},
	}

	PBToBson(msg, bsonM, recursiveCount)

	return bsonM.bsonD
}


func PBToBson(msg proto.Message, bsonCon BsonContainer, recursiveCount int) {
	if recursiveCount < 0 {
		common.LogErr("too many recursive")
		return
	}
	msgInsRef := proto.MessageReflect(msg)

	MsgDesAll := msgInsRef.Descriptor()
	if nil == MsgDesAll{
		return
	}
	MsgFields := MsgDesAll.Fields()
	if nil == MsgFields{
		return
	}

	for i := 0; i < MsgFields.Len(); i ++ {
		field := MsgFields.Get(i)
		if nil == field {
			continue
		}

		if !msgInsRef.Has(field) {
			continue
		}
		val := msgInsRef.Get(field)
		fieldName := string(field.Name())

		if field.Cardinality() == protoreflect.Repeated {
			if field.IsMap() {
				bsonSubMap := bson.M{}
				valMap := val.Map()
				valMap.Range(func(key protoreflect.MapKey, value protoreflect.Value) bool {
					bsonVal := PBGetValue(field.MapValue(), value, bsonCon, recursiveCount)
					bsonSubMap[key.Value().String()] = bsonVal
					return true
				})
				bsonCon.PutKV(fieldName, bsonSubMap)
			}
			if field.IsList() {
				// repeated
				bsonArray := primitive.A{}
				valList := val.List()
				for i := 0; i < valList.Len(); i ++ {
					value := valList.Get(i)
					bsonArray = append(bsonArray, PBGetValue(field, value, bsonCon, recursiveCount))
				}
				bsonCon.PutKV(fieldName, bsonArray)
			}
			continue
		}

		bsonCon.PutKV(fieldName, PBGetValue(field, val, bsonCon, recursiveCount))
	}

	return
}
