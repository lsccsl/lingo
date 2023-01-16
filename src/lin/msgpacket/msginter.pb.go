// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1-devel
// 	protoc        v3.19.1
// source: msginter.proto

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

type PB_MSG_INTER_TYPE int32

const (
	PB_MSG_INTER_TYPE__PB_MSG_INTER_NULL                    PB_MSG_INTER_TYPE = 0
	PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_REGISTER         PB_MSG_INTER_TYPE = 1
	PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_REGISTER_RES     PB_MSG_INTER_TYPE = 2
	PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_OFFLINE_NTF      PB_MSG_INTER_TYPE = 3
	PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_ONLINE_NTF       PB_MSG_INTER_TYPE = 4
	PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_CONNECT          PB_MSG_INTER_TYPE = 5
	PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_CONNECT_RES      PB_MSG_INTER_TYPE = 6
	PB_MSG_INTER_TYPE__PB_MSG_INTER_QUECENTER_HEARTBEAT     PB_MSG_INTER_TYPE = 7
	PB_MSG_INTER_TYPE__PB_MSG_INTER_QUECENTER_HEARTBEAT_RES PB_MSG_INTER_TYPE = 8
	PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_HEARTBEAT        PB_MSG_INTER_TYPE = 9
	PB_MSG_INTER_TYPE__PB_MSG_INTER_QUESRV_HEARTBEAT_RES    PB_MSG_INTER_TYPE = 10
)

// Enum value maps for PB_MSG_INTER_TYPE.
var (
	PB_MSG_INTER_TYPE_name = map[int32]string{
		0:  "_PB_MSG_INTER_NULL",
		1:  "_PB_MSG_INTER_QUESRV_REGISTER",
		2:  "_PB_MSG_INTER_QUESRV_REGISTER_RES",
		3:  "_PB_MSG_INTER_QUESRV_OFFLINE_NTF",
		4:  "_PB_MSG_INTER_QUESRV_ONLINE_NTF",
		5:  "_PB_MSG_INTER_QUESRV_CONNECT",
		6:  "_PB_MSG_INTER_QUESRV_CONNECT_RES",
		7:  "_PB_MSG_INTER_QUECENTER_HEARTBEAT",
		8:  "_PB_MSG_INTER_QUECENTER_HEARTBEAT_RES",
		9:  "_PB_MSG_INTER_QUESRV_HEARTBEAT",
		10: "_PB_MSG_INTER_QUESRV_HEARTBEAT_RES",
	}
	PB_MSG_INTER_TYPE_value = map[string]int32{
		"_PB_MSG_INTER_NULL":                    0,
		"_PB_MSG_INTER_QUESRV_REGISTER":         1,
		"_PB_MSG_INTER_QUESRV_REGISTER_RES":     2,
		"_PB_MSG_INTER_QUESRV_OFFLINE_NTF":      3,
		"_PB_MSG_INTER_QUESRV_ONLINE_NTF":       4,
		"_PB_MSG_INTER_QUESRV_CONNECT":          5,
		"_PB_MSG_INTER_QUESRV_CONNECT_RES":      6,
		"_PB_MSG_INTER_QUECENTER_HEARTBEAT":     7,
		"_PB_MSG_INTER_QUECENTER_HEARTBEAT_RES": 8,
		"_PB_MSG_INTER_QUESRV_HEARTBEAT":        9,
		"_PB_MSG_INTER_QUESRV_HEARTBEAT_RES":    10,
	}
)

func (x PB_MSG_INTER_TYPE) Enum() *PB_MSG_INTER_TYPE {
	p := new(PB_MSG_INTER_TYPE)
	*p = x
	return p
}

func (x PB_MSG_INTER_TYPE) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PB_MSG_INTER_TYPE) Descriptor() protoreflect.EnumDescriptor {
	return file_msginter_proto_enumTypes[0].Descriptor()
}

func (PB_MSG_INTER_TYPE) Type() protoreflect.EnumType {
	return &file_msginter_proto_enumTypes[0]
}

func (x PB_MSG_INTER_TYPE) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PB_MSG_INTER_TYPE.Descriptor instead.
func (PB_MSG_INTER_TYPE) EnumDescriptor() ([]byte, []int) {
	return file_msginter_proto_rawDescGZIP(), []int{0}
}

type PB_MSG_INTER_QUESRV_REGISTER struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ip   string `protobuf:"bytes,1,opt,name=ip,proto3" json:"ip,omitempty"`
	Port int32  `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
}

func (x *PB_MSG_INTER_QUESRV_REGISTER) Reset() {
	*x = PB_MSG_INTER_QUESRV_REGISTER{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msginter_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PB_MSG_INTER_QUESRV_REGISTER) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PB_MSG_INTER_QUESRV_REGISTER) ProtoMessage() {}

func (x *PB_MSG_INTER_QUESRV_REGISTER) ProtoReflect() protoreflect.Message {
	mi := &file_msginter_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PB_MSG_INTER_QUESRV_REGISTER.ProtoReflect.Descriptor instead.
func (*PB_MSG_INTER_QUESRV_REGISTER) Descriptor() ([]byte, []int) {
	return file_msginter_proto_rawDescGZIP(), []int{0}
}

func (x *PB_MSG_INTER_QUESRV_REGISTER) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *PB_MSG_INTER_QUESRV_REGISTER) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

type PB_MSG_INTER_QUESRV_INFO struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	QueSrvId int64  `protobuf:"varint,1,opt,name=que_srv_id,json=queSrvId,proto3" json:"que_srv_id,omitempty"`
	Ip       string `protobuf:"bytes,2,opt,name=ip,proto3" json:"ip,omitempty"`
	Port     int32  `protobuf:"varint,3,opt,name=port,proto3" json:"port,omitempty"`
}

func (x *PB_MSG_INTER_QUESRV_INFO) Reset() {
	*x = PB_MSG_INTER_QUESRV_INFO{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msginter_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PB_MSG_INTER_QUESRV_INFO) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PB_MSG_INTER_QUESRV_INFO) ProtoMessage() {}

func (x *PB_MSG_INTER_QUESRV_INFO) ProtoReflect() protoreflect.Message {
	mi := &file_msginter_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PB_MSG_INTER_QUESRV_INFO.ProtoReflect.Descriptor instead.
func (*PB_MSG_INTER_QUESRV_INFO) Descriptor() ([]byte, []int) {
	return file_msginter_proto_rawDescGZIP(), []int{1}
}

func (x *PB_MSG_INTER_QUESRV_INFO) GetQueSrvId() int64 {
	if x != nil {
		return x.QueSrvId
	}
	return 0
}

func (x *PB_MSG_INTER_QUESRV_INFO) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *PB_MSG_INTER_QUESRV_INFO) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

type PB_MSG_INTER_QUESRV_REGISTER_RES struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	QueSrvInfo []*PB_MSG_INTER_QUESRV_INFO `protobuf:"bytes,1,rep,name=que_srv_info,json=queSrvInfo,proto3" json:"que_srv_info,omitempty"`
	QueSrvId   int64                       `protobuf:"varint,2,opt,name=que_srv_id,json=queSrvId,proto3" json:"que_srv_id,omitempty"`
}

func (x *PB_MSG_INTER_QUESRV_REGISTER_RES) Reset() {
	*x = PB_MSG_INTER_QUESRV_REGISTER_RES{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msginter_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PB_MSG_INTER_QUESRV_REGISTER_RES) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PB_MSG_INTER_QUESRV_REGISTER_RES) ProtoMessage() {}

func (x *PB_MSG_INTER_QUESRV_REGISTER_RES) ProtoReflect() protoreflect.Message {
	mi := &file_msginter_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PB_MSG_INTER_QUESRV_REGISTER_RES.ProtoReflect.Descriptor instead.
func (*PB_MSG_INTER_QUESRV_REGISTER_RES) Descriptor() ([]byte, []int) {
	return file_msginter_proto_rawDescGZIP(), []int{2}
}

func (x *PB_MSG_INTER_QUESRV_REGISTER_RES) GetQueSrvInfo() []*PB_MSG_INTER_QUESRV_INFO {
	if x != nil {
		return x.QueSrvInfo
	}
	return nil
}

func (x *PB_MSG_INTER_QUESRV_REGISTER_RES) GetQueSrvId() int64 {
	if x != nil {
		return x.QueSrvId
	}
	return 0
}

type PB_MSG_INTER_QUESRV_OFFLINE_NTF struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	QueSrvId int64 `protobuf:"varint,1,opt,name=que_srv_id,json=queSrvId,proto3" json:"que_srv_id,omitempty"`
}

func (x *PB_MSG_INTER_QUESRV_OFFLINE_NTF) Reset() {
	*x = PB_MSG_INTER_QUESRV_OFFLINE_NTF{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msginter_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PB_MSG_INTER_QUESRV_OFFLINE_NTF) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PB_MSG_INTER_QUESRV_OFFLINE_NTF) ProtoMessage() {}

func (x *PB_MSG_INTER_QUESRV_OFFLINE_NTF) ProtoReflect() protoreflect.Message {
	mi := &file_msginter_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PB_MSG_INTER_QUESRV_OFFLINE_NTF.ProtoReflect.Descriptor instead.
func (*PB_MSG_INTER_QUESRV_OFFLINE_NTF) Descriptor() ([]byte, []int) {
	return file_msginter_proto_rawDescGZIP(), []int{3}
}

func (x *PB_MSG_INTER_QUESRV_OFFLINE_NTF) GetQueSrvId() int64 {
	if x != nil {
		return x.QueSrvId
	}
	return 0
}

type PB_MSG_INTER_QUESRV_ONLINE_NTF struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	QueSrvInfo *PB_MSG_INTER_QUESRV_INFO `protobuf:"bytes,1,opt,name=que_srv_info,json=queSrvInfo,proto3" json:"que_srv_info,omitempty"`
}

func (x *PB_MSG_INTER_QUESRV_ONLINE_NTF) Reset() {
	*x = PB_MSG_INTER_QUESRV_ONLINE_NTF{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msginter_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PB_MSG_INTER_QUESRV_ONLINE_NTF) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PB_MSG_INTER_QUESRV_ONLINE_NTF) ProtoMessage() {}

func (x *PB_MSG_INTER_QUESRV_ONLINE_NTF) ProtoReflect() protoreflect.Message {
	mi := &file_msginter_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PB_MSG_INTER_QUESRV_ONLINE_NTF.ProtoReflect.Descriptor instead.
func (*PB_MSG_INTER_QUESRV_ONLINE_NTF) Descriptor() ([]byte, []int) {
	return file_msginter_proto_rawDescGZIP(), []int{4}
}

func (x *PB_MSG_INTER_QUESRV_ONLINE_NTF) GetQueSrvInfo() *PB_MSG_INTER_QUESRV_INFO {
	if x != nil {
		return x.QueSrvInfo
	}
	return nil
}

type PB_MSG_INTER_QUESRV_CONNECT struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	QueSrvId int64 `protobuf:"varint,1,opt,name=que_srv_id,json=queSrvId,proto3" json:"que_srv_id,omitempty"`
}

func (x *PB_MSG_INTER_QUESRV_CONNECT) Reset() {
	*x = PB_MSG_INTER_QUESRV_CONNECT{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msginter_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PB_MSG_INTER_QUESRV_CONNECT) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PB_MSG_INTER_QUESRV_CONNECT) ProtoMessage() {}

func (x *PB_MSG_INTER_QUESRV_CONNECT) ProtoReflect() protoreflect.Message {
	mi := &file_msginter_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PB_MSG_INTER_QUESRV_CONNECT.ProtoReflect.Descriptor instead.
func (*PB_MSG_INTER_QUESRV_CONNECT) Descriptor() ([]byte, []int) {
	return file_msginter_proto_rawDescGZIP(), []int{5}
}

func (x *PB_MSG_INTER_QUESRV_CONNECT) GetQueSrvId() int64 {
	if x != nil {
		return x.QueSrvId
	}
	return 0
}

type PB_MSG_INTER_QUESRV_CONNECT_RES struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	QueSrvId int64 `protobuf:"varint,1,opt,name=que_srv_id,json=queSrvId,proto3" json:"que_srv_id,omitempty"`
}

func (x *PB_MSG_INTER_QUESRV_CONNECT_RES) Reset() {
	*x = PB_MSG_INTER_QUESRV_CONNECT_RES{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msginter_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PB_MSG_INTER_QUESRV_CONNECT_RES) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PB_MSG_INTER_QUESRV_CONNECT_RES) ProtoMessage() {}

func (x *PB_MSG_INTER_QUESRV_CONNECT_RES) ProtoReflect() protoreflect.Message {
	mi := &file_msginter_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PB_MSG_INTER_QUESRV_CONNECT_RES.ProtoReflect.Descriptor instead.
func (*PB_MSG_INTER_QUESRV_CONNECT_RES) Descriptor() ([]byte, []int) {
	return file_msginter_proto_rawDescGZIP(), []int{6}
}

func (x *PB_MSG_INTER_QUESRV_CONNECT_RES) GetQueSrvId() int64 {
	if x != nil {
		return x.QueSrvId
	}
	return 0
}

type PB_MSG_INTER_QUECENTER_HEARTBEAT struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	QueSrvId int64 `protobuf:"varint,1,opt,name=que_srv_id,json=queSrvId,proto3" json:"que_srv_id,omitempty"`
}

func (x *PB_MSG_INTER_QUECENTER_HEARTBEAT) Reset() {
	*x = PB_MSG_INTER_QUECENTER_HEARTBEAT{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msginter_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PB_MSG_INTER_QUECENTER_HEARTBEAT) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PB_MSG_INTER_QUECENTER_HEARTBEAT) ProtoMessage() {}

func (x *PB_MSG_INTER_QUECENTER_HEARTBEAT) ProtoReflect() protoreflect.Message {
	mi := &file_msginter_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PB_MSG_INTER_QUECENTER_HEARTBEAT.ProtoReflect.Descriptor instead.
func (*PB_MSG_INTER_QUECENTER_HEARTBEAT) Descriptor() ([]byte, []int) {
	return file_msginter_proto_rawDescGZIP(), []int{7}
}

func (x *PB_MSG_INTER_QUECENTER_HEARTBEAT) GetQueSrvId() int64 {
	if x != nil {
		return x.QueSrvId
	}
	return 0
}

type PB_MSG_INTER_QUECENTER_HEARTBEAT_RES struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	QueSrvId int64 `protobuf:"varint,1,opt,name=que_srv_id,json=queSrvId,proto3" json:"que_srv_id,omitempty"`
}

func (x *PB_MSG_INTER_QUECENTER_HEARTBEAT_RES) Reset() {
	*x = PB_MSG_INTER_QUECENTER_HEARTBEAT_RES{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msginter_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PB_MSG_INTER_QUECENTER_HEARTBEAT_RES) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PB_MSG_INTER_QUECENTER_HEARTBEAT_RES) ProtoMessage() {}

func (x *PB_MSG_INTER_QUECENTER_HEARTBEAT_RES) ProtoReflect() protoreflect.Message {
	mi := &file_msginter_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PB_MSG_INTER_QUECENTER_HEARTBEAT_RES.ProtoReflect.Descriptor instead.
func (*PB_MSG_INTER_QUECENTER_HEARTBEAT_RES) Descriptor() ([]byte, []int) {
	return file_msginter_proto_rawDescGZIP(), []int{8}
}

func (x *PB_MSG_INTER_QUECENTER_HEARTBEAT_RES) GetQueSrvId() int64 {
	if x != nil {
		return x.QueSrvId
	}
	return 0
}

type PB_MSG_INTER_QUESRV_HEARTBEAT struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	QueSrvId int64 `protobuf:"varint,1,opt,name=que_srv_id,json=queSrvId,proto3" json:"que_srv_id,omitempty"`
}

func (x *PB_MSG_INTER_QUESRV_HEARTBEAT) Reset() {
	*x = PB_MSG_INTER_QUESRV_HEARTBEAT{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msginter_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PB_MSG_INTER_QUESRV_HEARTBEAT) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PB_MSG_INTER_QUESRV_HEARTBEAT) ProtoMessage() {}

func (x *PB_MSG_INTER_QUESRV_HEARTBEAT) ProtoReflect() protoreflect.Message {
	mi := &file_msginter_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PB_MSG_INTER_QUESRV_HEARTBEAT.ProtoReflect.Descriptor instead.
func (*PB_MSG_INTER_QUESRV_HEARTBEAT) Descriptor() ([]byte, []int) {
	return file_msginter_proto_rawDescGZIP(), []int{9}
}

func (x *PB_MSG_INTER_QUESRV_HEARTBEAT) GetQueSrvId() int64 {
	if x != nil {
		return x.QueSrvId
	}
	return 0
}

type PB_MSG_INTER_QUESRV_HEARTBEAT_RES struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	QueSrvId int64 `protobuf:"varint,1,opt,name=que_srv_id,json=queSrvId,proto3" json:"que_srv_id,omitempty"`
}

func (x *PB_MSG_INTER_QUESRV_HEARTBEAT_RES) Reset() {
	*x = PB_MSG_INTER_QUESRV_HEARTBEAT_RES{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msginter_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PB_MSG_INTER_QUESRV_HEARTBEAT_RES) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PB_MSG_INTER_QUESRV_HEARTBEAT_RES) ProtoMessage() {}

func (x *PB_MSG_INTER_QUESRV_HEARTBEAT_RES) ProtoReflect() protoreflect.Message {
	mi := &file_msginter_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PB_MSG_INTER_QUESRV_HEARTBEAT_RES.ProtoReflect.Descriptor instead.
func (*PB_MSG_INTER_QUESRV_HEARTBEAT_RES) Descriptor() ([]byte, []int) {
	return file_msginter_proto_rawDescGZIP(), []int{10}
}

func (x *PB_MSG_INTER_QUESRV_HEARTBEAT_RES) GetQueSrvId() int64 {
	if x != nil {
		return x.QueSrvId
	}
	return 0
}

var File_msginter_proto protoreflect.FileDescriptor

var file_msginter_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x6d, 0x73, 0x67, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x09, 0x6d, 0x73, 0x67, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x22, 0x42, 0x0a, 0x1c, 0x50,
	0x42, 0x5f, 0x4d, 0x53, 0x47, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x5f, 0x51, 0x55, 0x45, 0x53,
	0x52, 0x56, 0x5f, 0x52, 0x45, 0x47, 0x49, 0x53, 0x54, 0x45, 0x52, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x70, 0x12, 0x12, 0x0a, 0x04, 0x70,
	0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x22,
	0x5c, 0x0a, 0x18, 0x50, 0x42, 0x5f, 0x4d, 0x53, 0x47, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x5f,
	0x51, 0x55, 0x45, 0x53, 0x52, 0x56, 0x5f, 0x49, 0x4e, 0x46, 0x4f, 0x12, 0x1c, 0x0a, 0x0a, 0x71,
	0x75, 0x65, 0x5f, 0x73, 0x72, 0x76, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x08, 0x71, 0x75, 0x65, 0x53, 0x72, 0x76, 0x49, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x70, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x70, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x22, 0x87, 0x01,
	0x0a, 0x20, 0x50, 0x42, 0x5f, 0x4d, 0x53, 0x47, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x5f, 0x51,
	0x55, 0x45, 0x53, 0x52, 0x56, 0x5f, 0x52, 0x45, 0x47, 0x49, 0x53, 0x54, 0x45, 0x52, 0x5f, 0x52,
	0x45, 0x53, 0x12, 0x45, 0x0a, 0x0c, 0x71, 0x75, 0x65, 0x5f, 0x73, 0x72, 0x76, 0x5f, 0x69, 0x6e,
	0x66, 0x6f, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x6d, 0x73, 0x67, 0x70, 0x61,
	0x63, 0x6b, 0x65, 0x74, 0x2e, 0x50, 0x42, 0x5f, 0x4d, 0x53, 0x47, 0x5f, 0x49, 0x4e, 0x54, 0x45,
	0x52, 0x5f, 0x51, 0x55, 0x45, 0x53, 0x52, 0x56, 0x5f, 0x49, 0x4e, 0x46, 0x4f, 0x52, 0x0a, 0x71,
	0x75, 0x65, 0x53, 0x72, 0x76, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1c, 0x0a, 0x0a, 0x71, 0x75, 0x65,
	0x5f, 0x73, 0x72, 0x76, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x71,
	0x75, 0x65, 0x53, 0x72, 0x76, 0x49, 0x64, 0x22, 0x3f, 0x0a, 0x1f, 0x50, 0x42, 0x5f, 0x4d, 0x53,
	0x47, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x5f, 0x51, 0x55, 0x45, 0x53, 0x52, 0x56, 0x5f, 0x4f,
	0x46, 0x46, 0x4c, 0x49, 0x4e, 0x45, 0x5f, 0x4e, 0x54, 0x46, 0x12, 0x1c, 0x0a, 0x0a, 0x71, 0x75,
	0x65, 0x5f, 0x73, 0x72, 0x76, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08,
	0x71, 0x75, 0x65, 0x53, 0x72, 0x76, 0x49, 0x64, 0x22, 0x67, 0x0a, 0x1e, 0x50, 0x42, 0x5f, 0x4d,
	0x53, 0x47, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x5f, 0x51, 0x55, 0x45, 0x53, 0x52, 0x56, 0x5f,
	0x4f, 0x4e, 0x4c, 0x49, 0x4e, 0x45, 0x5f, 0x4e, 0x54, 0x46, 0x12, 0x45, 0x0a, 0x0c, 0x71, 0x75,
	0x65, 0x5f, 0x73, 0x72, 0x76, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x23, 0x2e, 0x6d, 0x73, 0x67, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x2e, 0x50, 0x42, 0x5f,
	0x4d, 0x53, 0x47, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x5f, 0x51, 0x55, 0x45, 0x53, 0x52, 0x56,
	0x5f, 0x49, 0x4e, 0x46, 0x4f, 0x52, 0x0a, 0x71, 0x75, 0x65, 0x53, 0x72, 0x76, 0x49, 0x6e, 0x66,
	0x6f, 0x22, 0x3b, 0x0a, 0x1b, 0x50, 0x42, 0x5f, 0x4d, 0x53, 0x47, 0x5f, 0x49, 0x4e, 0x54, 0x45,
	0x52, 0x5f, 0x51, 0x55, 0x45, 0x53, 0x52, 0x56, 0x5f, 0x43, 0x4f, 0x4e, 0x4e, 0x45, 0x43, 0x54,
	0x12, 0x1c, 0x0a, 0x0a, 0x71, 0x75, 0x65, 0x5f, 0x73, 0x72, 0x76, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x71, 0x75, 0x65, 0x53, 0x72, 0x76, 0x49, 0x64, 0x22, 0x3f,
	0x0a, 0x1f, 0x50, 0x42, 0x5f, 0x4d, 0x53, 0x47, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x5f, 0x51,
	0x55, 0x45, 0x53, 0x52, 0x56, 0x5f, 0x43, 0x4f, 0x4e, 0x4e, 0x45, 0x43, 0x54, 0x5f, 0x52, 0x45,
	0x53, 0x12, 0x1c, 0x0a, 0x0a, 0x71, 0x75, 0x65, 0x5f, 0x73, 0x72, 0x76, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x71, 0x75, 0x65, 0x53, 0x72, 0x76, 0x49, 0x64, 0x22,
	0x40, 0x0a, 0x20, 0x50, 0x42, 0x5f, 0x4d, 0x53, 0x47, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x5f,
	0x51, 0x55, 0x45, 0x43, 0x45, 0x4e, 0x54, 0x45, 0x52, 0x5f, 0x48, 0x45, 0x41, 0x52, 0x54, 0x42,
	0x45, 0x41, 0x54, 0x12, 0x1c, 0x0a, 0x0a, 0x71, 0x75, 0x65, 0x5f, 0x73, 0x72, 0x76, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x71, 0x75, 0x65, 0x53, 0x72, 0x76, 0x49,
	0x64, 0x22, 0x44, 0x0a, 0x24, 0x50, 0x42, 0x5f, 0x4d, 0x53, 0x47, 0x5f, 0x49, 0x4e, 0x54, 0x45,
	0x52, 0x5f, 0x51, 0x55, 0x45, 0x43, 0x45, 0x4e, 0x54, 0x45, 0x52, 0x5f, 0x48, 0x45, 0x41, 0x52,
	0x54, 0x42, 0x45, 0x41, 0x54, 0x5f, 0x52, 0x45, 0x53, 0x12, 0x1c, 0x0a, 0x0a, 0x71, 0x75, 0x65,
	0x5f, 0x73, 0x72, 0x76, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x71,
	0x75, 0x65, 0x53, 0x72, 0x76, 0x49, 0x64, 0x22, 0x3d, 0x0a, 0x1d, 0x50, 0x42, 0x5f, 0x4d, 0x53,
	0x47, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x5f, 0x51, 0x55, 0x45, 0x53, 0x52, 0x56, 0x5f, 0x48,
	0x45, 0x41, 0x52, 0x54, 0x42, 0x45, 0x41, 0x54, 0x12, 0x1c, 0x0a, 0x0a, 0x71, 0x75, 0x65, 0x5f,
	0x73, 0x72, 0x76, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x71, 0x75,
	0x65, 0x53, 0x72, 0x76, 0x49, 0x64, 0x22, 0x41, 0x0a, 0x21, 0x50, 0x42, 0x5f, 0x4d, 0x53, 0x47,
	0x5f, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x5f, 0x51, 0x55, 0x45, 0x53, 0x52, 0x56, 0x5f, 0x48, 0x45,
	0x41, 0x52, 0x54, 0x42, 0x45, 0x41, 0x54, 0x5f, 0x52, 0x45, 0x53, 0x12, 0x1c, 0x0a, 0x0a, 0x71,
	0x75, 0x65, 0x5f, 0x73, 0x72, 0x76, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x08, 0x71, 0x75, 0x65, 0x53, 0x72, 0x76, 0x49, 0x64, 0x2a, 0xa6, 0x03, 0x0a, 0x11, 0x50, 0x42,
	0x5f, 0x4d, 0x53, 0x47, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x12,
	0x16, 0x0a, 0x12, 0x5f, 0x50, 0x42, 0x5f, 0x4d, 0x53, 0x47, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x52,
	0x5f, 0x4e, 0x55, 0x4c, 0x4c, 0x10, 0x00, 0x12, 0x21, 0x0a, 0x1d, 0x5f, 0x50, 0x42, 0x5f, 0x4d,
	0x53, 0x47, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x5f, 0x51, 0x55, 0x45, 0x53, 0x52, 0x56, 0x5f,
	0x52, 0x45, 0x47, 0x49, 0x53, 0x54, 0x45, 0x52, 0x10, 0x01, 0x12, 0x25, 0x0a, 0x21, 0x5f, 0x50,
	0x42, 0x5f, 0x4d, 0x53, 0x47, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x5f, 0x51, 0x55, 0x45, 0x53,
	0x52, 0x56, 0x5f, 0x52, 0x45, 0x47, 0x49, 0x53, 0x54, 0x45, 0x52, 0x5f, 0x52, 0x45, 0x53, 0x10,
	0x02, 0x12, 0x24, 0x0a, 0x20, 0x5f, 0x50, 0x42, 0x5f, 0x4d, 0x53, 0x47, 0x5f, 0x49, 0x4e, 0x54,
	0x45, 0x52, 0x5f, 0x51, 0x55, 0x45, 0x53, 0x52, 0x56, 0x5f, 0x4f, 0x46, 0x46, 0x4c, 0x49, 0x4e,
	0x45, 0x5f, 0x4e, 0x54, 0x46, 0x10, 0x03, 0x12, 0x23, 0x0a, 0x1f, 0x5f, 0x50, 0x42, 0x5f, 0x4d,
	0x53, 0x47, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x5f, 0x51, 0x55, 0x45, 0x53, 0x52, 0x56, 0x5f,
	0x4f, 0x4e, 0x4c, 0x49, 0x4e, 0x45, 0x5f, 0x4e, 0x54, 0x46, 0x10, 0x04, 0x12, 0x20, 0x0a, 0x1c,
	0x5f, 0x50, 0x42, 0x5f, 0x4d, 0x53, 0x47, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x5f, 0x51, 0x55,
	0x45, 0x53, 0x52, 0x56, 0x5f, 0x43, 0x4f, 0x4e, 0x4e, 0x45, 0x43, 0x54, 0x10, 0x05, 0x12, 0x24,
	0x0a, 0x20, 0x5f, 0x50, 0x42, 0x5f, 0x4d, 0x53, 0x47, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x5f,
	0x51, 0x55, 0x45, 0x53, 0x52, 0x56, 0x5f, 0x43, 0x4f, 0x4e, 0x4e, 0x45, 0x43, 0x54, 0x5f, 0x52,
	0x45, 0x53, 0x10, 0x06, 0x12, 0x25, 0x0a, 0x21, 0x5f, 0x50, 0x42, 0x5f, 0x4d, 0x53, 0x47, 0x5f,
	0x49, 0x4e, 0x54, 0x45, 0x52, 0x5f, 0x51, 0x55, 0x45, 0x43, 0x45, 0x4e, 0x54, 0x45, 0x52, 0x5f,
	0x48, 0x45, 0x41, 0x52, 0x54, 0x42, 0x45, 0x41, 0x54, 0x10, 0x07, 0x12, 0x29, 0x0a, 0x25, 0x5f,
	0x50, 0x42, 0x5f, 0x4d, 0x53, 0x47, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x5f, 0x51, 0x55, 0x45,
	0x43, 0x45, 0x4e, 0x54, 0x45, 0x52, 0x5f, 0x48, 0x45, 0x41, 0x52, 0x54, 0x42, 0x45, 0x41, 0x54,
	0x5f, 0x52, 0x45, 0x53, 0x10, 0x08, 0x12, 0x22, 0x0a, 0x1e, 0x5f, 0x50, 0x42, 0x5f, 0x4d, 0x53,
	0x47, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x5f, 0x51, 0x55, 0x45, 0x53, 0x52, 0x56, 0x5f, 0x48,
	0x45, 0x41, 0x52, 0x54, 0x42, 0x45, 0x41, 0x54, 0x10, 0x09, 0x12, 0x26, 0x0a, 0x22, 0x5f, 0x50,
	0x42, 0x5f, 0x4d, 0x53, 0x47, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x5f, 0x51, 0x55, 0x45, 0x53,
	0x52, 0x56, 0x5f, 0x48, 0x45, 0x41, 0x52, 0x54, 0x42, 0x45, 0x41, 0x54, 0x5f, 0x52, 0x45, 0x53,
	0x10, 0x0a, 0x42, 0x0e, 0x5a, 0x0c, 0x2e, 0x2f, 0x3b, 0x6d, 0x73, 0x67, 0x70, 0x61, 0x63, 0x6b,
	0x65, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_msginter_proto_rawDescOnce sync.Once
	file_msginter_proto_rawDescData = file_msginter_proto_rawDesc
)

func file_msginter_proto_rawDescGZIP() []byte {
	file_msginter_proto_rawDescOnce.Do(func() {
		file_msginter_proto_rawDescData = protoimpl.X.CompressGZIP(file_msginter_proto_rawDescData)
	})
	return file_msginter_proto_rawDescData
}

var file_msginter_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_msginter_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_msginter_proto_goTypes = []interface{}{
	(PB_MSG_INTER_TYPE)(0),                       // 0: msgpacket.PB_MSG_INTER_TYPE
	(*PB_MSG_INTER_QUESRV_REGISTER)(nil),         // 1: msgpacket.PB_MSG_INTER_QUESRV_REGISTER
	(*PB_MSG_INTER_QUESRV_INFO)(nil),             // 2: msgpacket.PB_MSG_INTER_QUESRV_INFO
	(*PB_MSG_INTER_QUESRV_REGISTER_RES)(nil),     // 3: msgpacket.PB_MSG_INTER_QUESRV_REGISTER_RES
	(*PB_MSG_INTER_QUESRV_OFFLINE_NTF)(nil),      // 4: msgpacket.PB_MSG_INTER_QUESRV_OFFLINE_NTF
	(*PB_MSG_INTER_QUESRV_ONLINE_NTF)(nil),       // 5: msgpacket.PB_MSG_INTER_QUESRV_ONLINE_NTF
	(*PB_MSG_INTER_QUESRV_CONNECT)(nil),          // 6: msgpacket.PB_MSG_INTER_QUESRV_CONNECT
	(*PB_MSG_INTER_QUESRV_CONNECT_RES)(nil),      // 7: msgpacket.PB_MSG_INTER_QUESRV_CONNECT_RES
	(*PB_MSG_INTER_QUECENTER_HEARTBEAT)(nil),     // 8: msgpacket.PB_MSG_INTER_QUECENTER_HEARTBEAT
	(*PB_MSG_INTER_QUECENTER_HEARTBEAT_RES)(nil), // 9: msgpacket.PB_MSG_INTER_QUECENTER_HEARTBEAT_RES
	(*PB_MSG_INTER_QUESRV_HEARTBEAT)(nil),        // 10: msgpacket.PB_MSG_INTER_QUESRV_HEARTBEAT
	(*PB_MSG_INTER_QUESRV_HEARTBEAT_RES)(nil),    // 11: msgpacket.PB_MSG_INTER_QUESRV_HEARTBEAT_RES
}
var file_msginter_proto_depIdxs = []int32{
	2, // 0: msgpacket.PB_MSG_INTER_QUESRV_REGISTER_RES.que_srv_info:type_name -> msgpacket.PB_MSG_INTER_QUESRV_INFO
	2, // 1: msgpacket.PB_MSG_INTER_QUESRV_ONLINE_NTF.que_srv_info:type_name -> msgpacket.PB_MSG_INTER_QUESRV_INFO
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_msginter_proto_init() }
func file_msginter_proto_init() {
	if File_msginter_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_msginter_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PB_MSG_INTER_QUESRV_REGISTER); i {
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
		file_msginter_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PB_MSG_INTER_QUESRV_INFO); i {
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
		file_msginter_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PB_MSG_INTER_QUESRV_REGISTER_RES); i {
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
		file_msginter_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PB_MSG_INTER_QUESRV_OFFLINE_NTF); i {
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
		file_msginter_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PB_MSG_INTER_QUESRV_ONLINE_NTF); i {
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
		file_msginter_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PB_MSG_INTER_QUESRV_CONNECT); i {
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
		file_msginter_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PB_MSG_INTER_QUESRV_CONNECT_RES); i {
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
		file_msginter_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PB_MSG_INTER_QUECENTER_HEARTBEAT); i {
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
		file_msginter_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PB_MSG_INTER_QUECENTER_HEARTBEAT_RES); i {
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
		file_msginter_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PB_MSG_INTER_QUESRV_HEARTBEAT); i {
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
		file_msginter_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PB_MSG_INTER_QUESRV_HEARTBEAT_RES); i {
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
			RawDescriptor: file_msginter_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_msginter_proto_goTypes,
		DependencyIndexes: file_msginter_proto_depIdxs,
		EnumInfos:         file_msginter_proto_enumTypes,
		MessageInfos:      file_msginter_proto_msgTypes,
	}.Build()
	File_msginter_proto = out.File
	file_msginter_proto_rawDesc = nil
	file_msginter_proto_goTypes = nil
	file_msginter_proto_depIdxs = nil
}
