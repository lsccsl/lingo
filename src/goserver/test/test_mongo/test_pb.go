package main

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

func PBGenFieldValue(kind protoreflect.Kind, des protoreflect.FieldDescriptor, v interface{}) (vRet protoreflect.Value) {
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
			mv1 := PBMsgGen(string(des.Message().FullName()), t).(protoiface.MessageV1)
			vRet = protoreflect.ValueOfMessage(proto.MessageReflect(mv1))
		default:
			common.LogErr("kind:", kind, " des:", des, " val:", v)
		}
		break

	case protoreflect.EnumKind:
		val, _ := strconv.Atoi(v.(string))
		vRet = protoreflect.ValueOfEnum(protoreflect.EnumNumber(int32(val)))
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

func PBMsgGen(MsgName string, MapKV map[string]interface{})interface{} {
	msgName := protoreflect.FullName(MsgName)
	msgType, err := protoregistry.GlobalTypes.FindMessageByName(msgName)
	if nil == msgType{
		common.LogErr(err)
		return nil
	}
	msgIns := proto.MessageV1(msgType.New())
	if nil == msgIns{
		return nil
	}
	msgInsRef:= proto.MessageReflect(msgIns)
	if nil == msgInsRef{
		return nil
	}

	MsgDesAll := msgInsRef.Descriptor()
	if nil == MsgDesAll{
		return nil
	}
	MsgFields := MsgDesAll.Fields()
	if nil == MsgFields{
		return nil
	}

	if nil == MapKV{
		return msgIns
	}

	for k,v := range MapKV{
		MsgDesField := MsgFields.ByTextName(k)
		if nil == MsgDesField{
			return nil
		}

		if MsgDesField.Cardinality() == protoreflect.Repeated{
			if MsgDesField.IsList(){
				muList := msgInsRef.Mutable(MsgDesField).List()
				var lst_input []interface{}
				switch t := v.(type) {
				case bson.A:
					{
						lst_input = t
						for _,lv := range lst_input{
							muList.Append(PBGenFieldValue(MsgDesField.Kind(), MsgDesField, lv))
						}
					}

				default:
					common.LogErr(t)
				}
			}
			if MsgDesField.IsMap(){
				muMap := msgInsRef.Mutable(MsgDesField).Map()
				switch t := v.(type) {
				case bson.M:
					for km,vm := range t{
						tmpK := PBGenFieldValue(MsgDesField.MapKey().Kind(), MsgDesField.MapKey(), km).MapKey()
						tmpV := PBGenFieldValue(MsgDesField.MapValue().Kind(), MsgDesField.MapValue(), vm)
						muMap.Set(tmpK, tmpV)
					}

				default:
					common.LogErr(t)
				}
			}
		} else {
			msgInsRef.Set(MsgDesField, PBGenFieldValue(MsgDesField.Kind(), MsgDesField, v))
		}
	}

	return msgIns
}

func PBToBson(msg proto.Message) bson.M {
	msgInsRef := proto.MessageReflect(msg)

	MsgDesAll := msgInsRef.Descriptor()
	if nil == MsgDesAll{
		return nil
	}
	MsgFields := MsgDesAll.Fields()
	if nil == MsgFields{
		return nil
	}

	bsonMap := bson.M{}

	for i := 0; i < MsgFields.Len(); i ++ {
		field := MsgFields.Get(i)
		if nil == field {
			continue
		}

		val := msgInsRef.Get(field)
		fieldName := string(field.Name())

		if field.Cardinality() == protoreflect.Repeated {
			if field.IsMap() {

			}
			if field.IsList() {

			}
			continue
		}

		switch field.Kind() {
		case protoreflect.BoolKind:
			bsonMap[fieldName] = val.Bool()
		case protoreflect.EnumKind:
			//field.Enum().
		case protoreflect.Int32Kind:
			bsonMap[fieldName] = val.Int()
		case protoreflect.Uint32Kind:
			bsonMap[fieldName] = val.Uint()
		case protoreflect.Int64Kind:
			bsonMap[fieldName] = val.Int()
		case protoreflect.Uint64Kind:
			bsonMap[fieldName] = val.Uint()
		case protoreflect.FloatKind:
			bsonMap[fieldName] = val.Float()
		case protoreflect.DoubleKind:
			bsonMap[fieldName] = val.Float()
		case protoreflect.StringKind:
			bsonMap[fieldName] = val.String()
		case protoreflect.BytesKind:
			bsonMap[fieldName] = val.Bytes()
		case protoreflect.MessageKind:
			{

			}
		case protoreflect.GroupKind:
			{

			}
		case protoreflect.Sint64Kind:
			bsonMap[fieldName] = val.Int()
		case protoreflect.Sint32Kind:
			bsonMap[fieldName] = val.Int()
		case protoreflect.Sfixed32Kind:
			bsonMap[fieldName] = val.Int()
		case protoreflect.Fixed32Kind:
			bsonMap[fieldName] = val.Int()
		case protoreflect.Sfixed64Kind:
			bsonMap[fieldName] = val.Int()
		case protoreflect.Fixed64Kind:
			bsonMap[fieldName] = val.Int()
		default:
			common.LogErr("unknow protocal type:", field.Kind())
		}
	}

	return bsonMap
}