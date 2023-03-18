// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1-devel
// 	protoc        v3.19.1
// source: msgDB.proto

package msgpacket

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type EN_TEST int32

const (
	EN_TEST_EN_TEST0 EN_TEST = 0
	EN_TEST_EN_TEST1 EN_TEST = 1
	EN_TEST_EN_TEST2 EN_TEST = 2
	EN_TEST_EN_TEST3 EN_TEST = 3
)

// Enum value maps for EN_TEST.
var (
	EN_TEST_name = map[int32]string{
		0: "EN_TEST0",
		1: "EN_TEST1",
		2: "EN_TEST2",
		3: "EN_TEST3",
	}
	EN_TEST_value = map[string]int32{
		"EN_TEST0": 0,
		"EN_TEST1": 1,
		"EN_TEST2": 2,
		"EN_TEST3": 3,
	}
)

func (x EN_TEST) Enum() *EN_TEST {
	p := new(EN_TEST)
	*p = x
	return p
}

func (x EN_TEST) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EN_TEST) Descriptor() protoreflect.EnumDescriptor {
	return file_msgDB_proto_enumTypes[0].Descriptor()
}

func (EN_TEST) Type() protoreflect.EnumType {
	return &file_msgDB_proto_enumTypes[0]
}

func (x EN_TEST) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use EN_TEST.Descriptor instead.
func (EN_TEST) EnumDescriptor() ([]byte, []int) {
	return file_msgDB_proto_rawDescGZIP(), []int{0}
}

type DBUserMainKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId int64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *DBUserMainKey) Reset() {
	*x = DBUserMainKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msgDB_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DBUserMainKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DBUserMainKey) ProtoMessage() {}

func (x *DBUserMainKey) ProtoReflect() protoreflect.Message {
	mi := &file_msgDB_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DBUserMainKey.ProtoReflect.Descriptor instead.
func (*DBUserMainKey) Descriptor() ([]byte, []int) {
	return file_msgDB_proto_rawDescGZIP(), []int{0}
}

func (x *DBUserMainKey) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type DBUserMain struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	XId    string `protobuf:"bytes,1,opt,name=_id,json=Id,proto3" json:"_id,omitempty"`
	UserId int64  `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *DBUserMain) Reset() {
	*x = DBUserMain{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msgDB_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DBUserMain) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DBUserMain) ProtoMessage() {}

func (x *DBUserMain) ProtoReflect() protoreflect.Message {
	mi := &file_msgDB_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DBUserMain.ProtoReflect.Descriptor instead.
func (*DBUserMain) Descriptor() ([]byte, []int) {
	return file_msgDB_proto_rawDescGZIP(), []int{1}
}

func (x *DBUserMain) GetXId() string {
	if x != nil {
		return x.XId
	}
	return ""
}

func (x *DBUserMain) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type DBUserDetailKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId int64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *DBUserDetailKey) Reset() {
	*x = DBUserDetailKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msgDB_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DBUserDetailKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DBUserDetailKey) ProtoMessage() {}

func (x *DBUserDetailKey) ProtoReflect() protoreflect.Message {
	mi := &file_msgDB_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DBUserDetailKey.ProtoReflect.Descriptor instead.
func (*DBUserDetailKey) Descriptor() ([]byte, []int) {
	return file_msgDB_proto_rawDescGZIP(), []int{2}
}

func (x *DBUserDetailKey) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type DBUserDetail struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	XId        string `protobuf:"bytes,1,opt,name=_id,json=Id,proto3" json:"_id,omitempty"`
	DetailData string `protobuf:"bytes,2,opt,name=detail_data,json=detailData,proto3" json:"detail_data,omitempty"`
	DetailId   int32  `protobuf:"varint,3,opt,name=detail_id,json=detailId,proto3" json:"detail_id,omitempty"`
}

func (x *DBUserDetail) Reset() {
	*x = DBUserDetail{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msgDB_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DBUserDetail) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DBUserDetail) ProtoMessage() {}

func (x *DBUserDetail) ProtoReflect() protoreflect.Message {
	mi := &file_msgDB_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DBUserDetail.ProtoReflect.Descriptor instead.
func (*DBUserDetail) Descriptor() ([]byte, []int) {
	return file_msgDB_proto_rawDescGZIP(), []int{3}
}

func (x *DBUserDetail) GetXId() string {
	if x != nil {
		return x.XId
	}
	return ""
}

func (x *DBUserDetail) GetDetailData() string {
	if x != nil {
		return x.DetailData
	}
	return ""
}

func (x *DBUserDetail) GetDetailId() int32 {
	if x != nil {
		return x.DetailId
	}
	return 0
}

// ============================
type DBUserMainTest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	XId          string            `protobuf:"bytes,1,opt,name=_id,json=Id,proto3" json:"_id,omitempty"`
	UserId       int64             `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Detail       *DBUserDetail     `protobuf:"bytes,3,opt,name=detail,proto3" json:"detail,omitempty"`
	TestRepeated []*DBRepeatedTest `protobuf:"bytes,4,rep,name=test_repeated,json=testRepeated,proto3" json:"test_repeated,omitempty"`
	EnTest       EN_TEST           `protobuf:"varint,5,opt,name=en_test,json=enTest,proto3,enum=msgpacket.EN_TEST" json:"en_test,omitempty"`
	Str1         string            `protobuf:"bytes,6,opt,name=str1,proto3" json:"str1,omitempty"`
	Str2         string            `protobuf:"bytes,7,opt,name=str2,proto3" json:"str2,omitempty"`
	Int1         int32             `protobuf:"varint,8,opt,name=int1,proto3" json:"int1,omitempty"`
	Int2         int32             `protobuf:"varint,9,opt,name=int2,proto3" json:"int2,omitempty"`
}

func (x *DBUserMainTest) Reset() {
	*x = DBUserMainTest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msgDB_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DBUserMainTest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DBUserMainTest) ProtoMessage() {}

func (x *DBUserMainTest) ProtoReflect() protoreflect.Message {
	mi := &file_msgDB_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DBUserMainTest.ProtoReflect.Descriptor instead.
func (*DBUserMainTest) Descriptor() ([]byte, []int) {
	return file_msgDB_proto_rawDescGZIP(), []int{4}
}

func (x *DBUserMainTest) GetXId() string {
	if x != nil {
		return x.XId
	}
	return ""
}

func (x *DBUserMainTest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *DBUserMainTest) GetDetail() *DBUserDetail {
	if x != nil {
		return x.Detail
	}
	return nil
}

func (x *DBUserMainTest) GetTestRepeated() []*DBRepeatedTest {
	if x != nil {
		return x.TestRepeated
	}
	return nil
}

func (x *DBUserMainTest) GetEnTest() EN_TEST {
	if x != nil {
		return x.EnTest
	}
	return EN_TEST_EN_TEST0
}

func (x *DBUserMainTest) GetStr1() string {
	if x != nil {
		return x.Str1
	}
	return ""
}

func (x *DBUserMainTest) GetStr2() string {
	if x != nil {
		return x.Str2
	}
	return ""
}

func (x *DBUserMainTest) GetInt1() int32 {
	if x != nil {
		return x.Int1
	}
	return 0
}

func (x *DBUserMainTest) GetInt2() int32 {
	if x != nil {
		return x.Int2
	}
	return 0
}

type DBMapTest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MapStr string `protobuf:"bytes,1,opt,name=map_str,json=mapStr,proto3" json:"map_str,omitempty"`
	MapInt int64  `protobuf:"varint,2,opt,name=map_int,json=mapInt,proto3" json:"map_int,omitempty"`
}

func (x *DBMapTest) Reset() {
	*x = DBMapTest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msgDB_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DBMapTest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DBMapTest) ProtoMessage() {}

func (x *DBMapTest) ProtoReflect() protoreflect.Message {
	mi := &file_msgDB_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DBMapTest.ProtoReflect.Descriptor instead.
func (*DBMapTest) Descriptor() ([]byte, []int) {
	return file_msgDB_proto_rawDescGZIP(), []int{5}
}

func (x *DBMapTest) GetMapStr() string {
	if x != nil {
		return x.MapStr
	}
	return ""
}

func (x *DBMapTest) GetMapInt() int64 {
	if x != nil {
		return x.MapInt
	}
	return 0
}

type DBRepeatedTest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RepeatedStr string               `protobuf:"bytes,1,opt,name=repeated_str,json=repeatedStr,proto3" json:"repeated_str,omitempty"`
	RepeatedInt int64                `protobuf:"varint,2,opt,name=repeated_int,json=repeatedInt,proto3" json:"repeated_int,omitempty"`
	TestMap     map[int64]*DBMapTest `protobuf:"bytes,3,rep,name=test_map,json=testMap,proto3" json:"test_map,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *DBRepeatedTest) Reset() {
	*x = DBRepeatedTest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msgDB_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DBRepeatedTest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DBRepeatedTest) ProtoMessage() {}

func (x *DBRepeatedTest) ProtoReflect() protoreflect.Message {
	mi := &file_msgDB_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DBRepeatedTest.ProtoReflect.Descriptor instead.
func (*DBRepeatedTest) Descriptor() ([]byte, []int) {
	return file_msgDB_proto_rawDescGZIP(), []int{6}
}

func (x *DBRepeatedTest) GetRepeatedStr() string {
	if x != nil {
		return x.RepeatedStr
	}
	return ""
}

func (x *DBRepeatedTest) GetRepeatedInt() int64 {
	if x != nil {
		return x.RepeatedInt
	}
	return 0
}

func (x *DBRepeatedTest) GetTestMap() map[int64]*DBMapTest {
	if x != nil {
		return x.TestMap
	}
	return nil
}

var File_msgDB_proto protoreflect.FileDescriptor

var file_msgDB_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x6d, 0x73, 0x67, 0x44, 0x42, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x6d,
	0x73, 0x67, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x22, 0x28, 0x0a, 0x0d, 0x44, 0x42, 0x55, 0x73,
	0x65, 0x72, 0x4d, 0x61, 0x69, 0x6e, 0x4b, 0x65, 0x79, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x22, 0x36, 0x0a, 0x0a, 0x44, 0x42, 0x55, 0x73, 0x65, 0x72, 0x4d, 0x61, 0x69, 0x6e,
	0x12, 0x0f, 0x0a, 0x03, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49,
	0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x2a, 0x0a, 0x0f, 0x44, 0x42,
	0x55, 0x73, 0x65, 0x72, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x4b, 0x65, 0x79, 0x12, 0x17, 0x0a,
	0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x5d, 0x0a, 0x0c, 0x44, 0x42, 0x55, 0x73, 0x65, 0x72,
	0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x12, 0x0f, 0x0a, 0x03, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x64, 0x65, 0x74, 0x61, 0x69,
	0x6c, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x64, 0x65,
	0x74, 0x61, 0x69, 0x6c, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1b, 0x0a, 0x09, 0x64, 0x65, 0x74, 0x61,
	0x69, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x64, 0x65, 0x74,
	0x61, 0x69, 0x6c, 0x49, 0x64, 0x22, 0xa8, 0x02, 0x0a, 0x0e, 0x44, 0x42, 0x55, 0x73, 0x65, 0x72,
	0x4d, 0x61, 0x69, 0x6e, 0x54, 0x65, 0x73, 0x74, 0x12, 0x0f, 0x0a, 0x03, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x12, 0x2f, 0x0a, 0x06, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x17, 0x2e, 0x6d, 0x73, 0x67, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x2e, 0x44,
	0x42, 0x55, 0x73, 0x65, 0x72, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x52, 0x06, 0x64, 0x65, 0x74,
	0x61, 0x69, 0x6c, 0x12, 0x3e, 0x0a, 0x0d, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x72, 0x65, 0x70, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x6d, 0x73, 0x67,
	0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x2e, 0x44, 0x42, 0x52, 0x65, 0x70, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x54, 0x65, 0x73, 0x74, 0x52, 0x0c, 0x74, 0x65, 0x73, 0x74, 0x52, 0x65, 0x70, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x12, 0x2b, 0x0a, 0x07, 0x65, 0x6e, 0x5f, 0x74, 0x65, 0x73, 0x74, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x6d, 0x73, 0x67, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74,
	0x2e, 0x45, 0x4e, 0x5f, 0x54, 0x45, 0x53, 0x54, 0x52, 0x06, 0x65, 0x6e, 0x54, 0x65, 0x73, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x73, 0x74, 0x72, 0x31, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x73, 0x74, 0x72, 0x31, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x74, 0x72, 0x32, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x73, 0x74, 0x72, 0x32, 0x12, 0x12, 0x0a, 0x04, 0x69, 0x6e, 0x74, 0x31,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x69, 0x6e, 0x74, 0x31, 0x12, 0x12, 0x0a, 0x04,
	0x69, 0x6e, 0x74, 0x32, 0x18, 0x09, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x69, 0x6e, 0x74, 0x32,
	0x22, 0x3d, 0x0a, 0x09, 0x44, 0x42, 0x4d, 0x61, 0x70, 0x54, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a,
	0x07, 0x6d, 0x61, 0x70, 0x5f, 0x73, 0x74, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x6d, 0x61, 0x70, 0x53, 0x74, 0x72, 0x12, 0x17, 0x0a, 0x07, 0x6d, 0x61, 0x70, 0x5f, 0x69, 0x6e,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x6d, 0x61, 0x70, 0x49, 0x6e, 0x74, 0x22,
	0xeb, 0x01, 0x0a, 0x0e, 0x44, 0x42, 0x52, 0x65, 0x70, 0x65, 0x61, 0x74, 0x65, 0x64, 0x54, 0x65,
	0x73, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x72, 0x65, 0x70, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x73,
	0x74, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x72, 0x65, 0x70, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x53, 0x74, 0x72, 0x12, 0x21, 0x0a, 0x0c, 0x72, 0x65, 0x70, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x5f, 0x69, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x72, 0x65, 0x70,
	0x65, 0x61, 0x74, 0x65, 0x64, 0x49, 0x6e, 0x74, 0x12, 0x41, 0x0a, 0x08, 0x74, 0x65, 0x73, 0x74,
	0x5f, 0x6d, 0x61, 0x70, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x6d, 0x73, 0x67,
	0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x2e, 0x44, 0x42, 0x52, 0x65, 0x70, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x54, 0x65, 0x73, 0x74, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x4d, 0x61, 0x70, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x52, 0x07, 0x74, 0x65, 0x73, 0x74, 0x4d, 0x61, 0x70, 0x1a, 0x50, 0x0a, 0x0c, 0x54,
	0x65, 0x73, 0x74, 0x4d, 0x61, 0x70, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x2a, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x6d,
	0x73, 0x67, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x2e, 0x44, 0x42, 0x4d, 0x61, 0x70, 0x54, 0x65,
	0x73, 0x74, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x2a, 0x41, 0x0a,
	0x07, 0x45, 0x4e, 0x5f, 0x54, 0x45, 0x53, 0x54, 0x12, 0x0c, 0x0a, 0x08, 0x45, 0x4e, 0x5f, 0x54,
	0x45, 0x53, 0x54, 0x30, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x45, 0x4e, 0x5f, 0x54, 0x45, 0x53,
	0x54, 0x31, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x45, 0x4e, 0x5f, 0x54, 0x45, 0x53, 0x54, 0x32,
	0x10, 0x02, 0x12, 0x0c, 0x0a, 0x08, 0x45, 0x4e, 0x5f, 0x54, 0x45, 0x53, 0x54, 0x33, 0x10, 0x03,
	0x42, 0x0e, 0x5a, 0x0c, 0x2e, 0x2f, 0x3b, 0x6d, 0x73, 0x67, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_msgDB_proto_rawDescOnce sync.Once
	file_msgDB_proto_rawDescData = file_msgDB_proto_rawDesc
)

func file_msgDB_proto_rawDescGZIP() []byte {
	file_msgDB_proto_rawDescOnce.Do(func() {
		file_msgDB_proto_rawDescData = protoimpl.X.CompressGZIP(file_msgDB_proto_rawDescData)
	})
	return file_msgDB_proto_rawDescData
}

var file_msgDB_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_msgDB_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_msgDB_proto_goTypes = []interface{}{
	(EN_TEST)(0),            // 0: msgpacket.EN_TEST
	(*DBUserMainKey)(nil),   // 1: msgpacket.DBUserMainKey
	(*DBUserMain)(nil),      // 2: msgpacket.DBUserMain
	(*DBUserDetailKey)(nil), // 3: msgpacket.DBUserDetailKey
	(*DBUserDetail)(nil),    // 4: msgpacket.DBUserDetail
	(*DBUserMainTest)(nil),  // 5: msgpacket.DBUserMainTest
	(*DBMapTest)(nil),       // 6: msgpacket.DBMapTest
	(*DBRepeatedTest)(nil),  // 7: msgpacket.DBRepeatedTest
	nil,                     // 8: msgpacket.DBRepeatedTest.TestMapEntry
}
var file_msgDB_proto_depIdxs = []int32{
	4, // 0: msgpacket.DBUserMainTest.detail:type_name -> msgpacket.DBUserDetail
	7, // 1: msgpacket.DBUserMainTest.test_repeated:type_name -> msgpacket.DBRepeatedTest
	0, // 2: msgpacket.DBUserMainTest.en_test:type_name -> msgpacket.EN_TEST
	8, // 3: msgpacket.DBRepeatedTest.test_map:type_name -> msgpacket.DBRepeatedTest.TestMapEntry
	6, // 4: msgpacket.DBRepeatedTest.TestMapEntry.value:type_name -> msgpacket.DBMapTest
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_msgDB_proto_init() }
func file_msgDB_proto_init() {
	if File_msgDB_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_msgDB_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DBUserMainKey); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_msgDB_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DBUserMain); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_msgDB_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DBUserDetailKey); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_msgDB_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DBUserDetail); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_msgDB_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DBUserMainTest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_msgDB_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DBMapTest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_msgDB_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DBRepeatedTest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_msgDB_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_msgDB_proto_goTypes,
		DependencyIndexes: file_msgDB_proto_depIdxs,
		EnumInfos:         file_msgDB_proto_enumTypes,
		MessageInfos:      file_msgDB_proto_msgTypes,
	}.Build()
	File_msgDB_proto = out.File
	file_msgDB_proto_rawDesc = nil
	file_msgDB_proto_goTypes = nil
	file_msgDB_proto_depIdxs = nil
}
