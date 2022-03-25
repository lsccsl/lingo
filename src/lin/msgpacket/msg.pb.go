// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1-devel
// 	protoc        v3.19.1
// source: msg.proto

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

type MSG_TYPE int32

const (
	MSG_TYPE__MSG_NULL           MSG_TYPE = 0
	MSG_TYPE__MSG_RPC            MSG_TYPE = 1
	MSG_TYPE__MSG_RPC_RES        MSG_TYPE = 2
	MSG_TYPE__MSG_SRV_REPORT     MSG_TYPE = 3
	MSG_TYPE__MSG_SRV_REPORT_RES MSG_TYPE = 4
	MSG_TYPE__MSG_HEARTBEAT      MSG_TYPE = 5
	MSG_TYPE__MSG_HEARTBEAT_RES  MSG_TYPE = 6
	MSG_TYPE__MSG_TCP_STATIC     MSG_TYPE = 7
	MSG_TYPE__MSG_TCP_STATIC_RES MSG_TYPE = 8
	MSG_TYPE__MSG_MAX            MSG_TYPE = 100
	MSG_TYPE__MSG_TEST           MSG_TYPE = 101
	MSG_TYPE__MSG_TEST_RES       MSG_TYPE = 102
	MSG_TYPE__MSG_LOGIN          MSG_TYPE = 103
	MSG_TYPE__MSG_LOGIN_RES      MSG_TYPE = 104
)

// Enum value maps for MSG_TYPE.
var (
	MSG_TYPE_name = map[int32]string{
		0:   "_MSG_NULL",
		1:   "_MSG_RPC",
		2:   "_MSG_RPC_RES",
		3:   "_MSG_SRV_REPORT",
		4:   "_MSG_SRV_REPORT_RES",
		5:   "_MSG_HEARTBEAT",
		6:   "_MSG_HEARTBEAT_RES",
		7:   "_MSG_TCP_STATIC",
		8:   "_MSG_TCP_STATIC_RES",
		100: "_MSG_MAX",
		101: "_MSG_TEST",
		102: "_MSG_TEST_RES",
		103: "_MSG_LOGIN",
		104: "_MSG_LOGIN_RES",
	}
	MSG_TYPE_value = map[string]int32{
		"_MSG_NULL":           0,
		"_MSG_RPC":            1,
		"_MSG_RPC_RES":        2,
		"_MSG_SRV_REPORT":     3,
		"_MSG_SRV_REPORT_RES": 4,
		"_MSG_HEARTBEAT":      5,
		"_MSG_HEARTBEAT_RES":  6,
		"_MSG_TCP_STATIC":     7,
		"_MSG_TCP_STATIC_RES": 8,
		"_MSG_MAX":            100,
		"_MSG_TEST":           101,
		"_MSG_TEST_RES":       102,
		"_MSG_LOGIN":          103,
		"_MSG_LOGIN_RES":      104,
	}
)

func (x MSG_TYPE) Enum() *MSG_TYPE {
	p := new(MSG_TYPE)
	*p = x
	return p
}

func (x MSG_TYPE) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MSG_TYPE) Descriptor() protoreflect.EnumDescriptor {
	return file_msg_proto_enumTypes[0].Descriptor()
}

func (MSG_TYPE) Type() protoreflect.EnumType {
	return &file_msg_proto_enumTypes[0]
}

func (x MSG_TYPE) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MSG_TYPE.Descriptor instead.
func (MSG_TYPE) EnumDescriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{0}
}

type RESPONSE_CODE int32

const (
	RESPONSE_CODE_RESPONSE_CODE_NONE           RESPONSE_CODE = 0
	RESPONSE_CODE_RESPONSE_CODE_Fail           RESPONSE_CODE = 1
	RESPONSE_CODE_RESPONSE_CODE_Rpc_not_accept RESPONSE_CODE = 2
)

// Enum value maps for RESPONSE_CODE.
var (
	RESPONSE_CODE_name = map[int32]string{
		0: "RESPONSE_CODE_NONE",
		1: "RESPONSE_CODE_Fail",
		2: "RESPONSE_CODE_Rpc_not_accept",
	}
	RESPONSE_CODE_value = map[string]int32{
		"RESPONSE_CODE_NONE":           0,
		"RESPONSE_CODE_Fail":           1,
		"RESPONSE_CODE_Rpc_not_accept": 2,
	}
)

func (x RESPONSE_CODE) Enum() *RESPONSE_CODE {
	p := new(RESPONSE_CODE)
	*p = x
	return p
}

func (x RESPONSE_CODE) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RESPONSE_CODE) Descriptor() protoreflect.EnumDescriptor {
	return file_msg_proto_enumTypes[1].Descriptor()
}

func (RESPONSE_CODE) Type() protoreflect.EnumType {
	return &file_msg_proto_enumTypes[1]
}

func (x RESPONSE_CODE) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RESPONSE_CODE.Descriptor instead.
func (RESPONSE_CODE) EnumDescriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{1}
}

type MSG_RPC struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MsgId   int64  `protobuf:"varint,1,opt,name=msg_id,json=msgId,proto3" json:"msg_id,omitempty"`
	MsgType int32  `protobuf:"varint,2,opt,name=msg_type,json=msgType,proto3" json:"msg_type,omitempty"`
	MsgBin  []byte `protobuf:"bytes,3,opt,name=msg_bin,json=msgBin,proto3" json:"msg_bin,omitempty"`
}

func (x *MSG_RPC) Reset() {
	*x = MSG_RPC{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MSG_RPC) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MSG_RPC) ProtoMessage() {}

func (x *MSG_RPC) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MSG_RPC.ProtoReflect.Descriptor instead.
func (*MSG_RPC) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{0}
}

func (x *MSG_RPC) GetMsgId() int64 {
	if x != nil {
		return x.MsgId
	}
	return 0
}

func (x *MSG_RPC) GetMsgType() int32 {
	if x != nil {
		return x.MsgType
	}
	return 0
}

func (x *MSG_RPC) GetMsgBin() []byte {
	if x != nil {
		return x.MsgBin
	}
	return nil
}

type MSG_RPC_RES struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MsgId   int64         `protobuf:"varint,1,opt,name=msg_id,json=msgId,proto3" json:"msg_id,omitempty"`
	MsgType int32         `protobuf:"varint,2,opt,name=msg_type,json=msgType,proto3" json:"msg_type,omitempty"`
	ResCode RESPONSE_CODE `protobuf:"varint,3,opt,name=res_code,json=resCode,proto3,enum=msgpacket.RESPONSE_CODE" json:"res_code,omitempty"`
	MsgBin  []byte        `protobuf:"bytes,4,opt,name=msg_bin,json=msgBin,proto3" json:"msg_bin,omitempty"`
}

func (x *MSG_RPC_RES) Reset() {
	*x = MSG_RPC_RES{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MSG_RPC_RES) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MSG_RPC_RES) ProtoMessage() {}

func (x *MSG_RPC_RES) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MSG_RPC_RES.ProtoReflect.Descriptor instead.
func (*MSG_RPC_RES) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{1}
}

func (x *MSG_RPC_RES) GetMsgId() int64 {
	if x != nil {
		return x.MsgId
	}
	return 0
}

func (x *MSG_RPC_RES) GetMsgType() int32 {
	if x != nil {
		return x.MsgType
	}
	return 0
}

func (x *MSG_RPC_RES) GetResCode() RESPONSE_CODE {
	if x != nil {
		return x.ResCode
	}
	return RESPONSE_CODE_RESPONSE_CODE_NONE
}

func (x *MSG_RPC_RES) GetMsgBin() []byte {
	if x != nil {
		return x.MsgBin
	}
	return nil
}

type MSG_SRV_REPORT struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SrvId int64 `protobuf:"varint,1,opt,name=srv_id,json=srvId,proto3" json:"srv_id,omitempty"`
}

func (x *MSG_SRV_REPORT) Reset() {
	*x = MSG_SRV_REPORT{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MSG_SRV_REPORT) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MSG_SRV_REPORT) ProtoMessage() {}

func (x *MSG_SRV_REPORT) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MSG_SRV_REPORT.ProtoReflect.Descriptor instead.
func (*MSG_SRV_REPORT) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{2}
}

func (x *MSG_SRV_REPORT) GetSrvId() int64 {
	if x != nil {
		return x.SrvId
	}
	return 0
}

type MSG_SRV_REPORT_RES struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SrvId     int64 `protobuf:"varint,1,opt,name=srv_id,json=srvId,proto3" json:"srv_id,omitempty"`
	TcpConnId int64 `protobuf:"varint,2,opt,name=tcp_conn_id,json=tcpConnId,proto3" json:"tcp_conn_id,omitempty"`
}

func (x *MSG_SRV_REPORT_RES) Reset() {
	*x = MSG_SRV_REPORT_RES{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MSG_SRV_REPORT_RES) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MSG_SRV_REPORT_RES) ProtoMessage() {}

func (x *MSG_SRV_REPORT_RES) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MSG_SRV_REPORT_RES.ProtoReflect.Descriptor instead.
func (*MSG_SRV_REPORT_RES) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{3}
}

func (x *MSG_SRV_REPORT_RES) GetSrvId() int64 {
	if x != nil {
		return x.SrvId
	}
	return 0
}

func (x *MSG_SRV_REPORT_RES) GetTcpConnId() int64 {
	if x != nil {
		return x.TcpConnId
	}
	return 0
}

type MSG_HEARTBEAT struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *MSG_HEARTBEAT) Reset() {
	*x = MSG_HEARTBEAT{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MSG_HEARTBEAT) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MSG_HEARTBEAT) ProtoMessage() {}

func (x *MSG_HEARTBEAT) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MSG_HEARTBEAT.ProtoReflect.Descriptor instead.
func (*MSG_HEARTBEAT) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{4}
}

func (x *MSG_HEARTBEAT) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type MSG_HEARTBEAT_RES struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *MSG_HEARTBEAT_RES) Reset() {
	*x = MSG_HEARTBEAT_RES{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MSG_HEARTBEAT_RES) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MSG_HEARTBEAT_RES) ProtoMessage() {}

func (x *MSG_HEARTBEAT_RES) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MSG_HEARTBEAT_RES.ProtoReflect.Descriptor instead.
func (*MSG_HEARTBEAT_RES) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{5}
}

func (x *MSG_HEARTBEAT_RES) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type MSG_TEST struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id              int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Str             string `protobuf:"bytes,2,opt,name=str,proto3" json:"str,omitempty"`
	Seq             int64  `protobuf:"varint,3,opt,name=seq,proto3" json:"seq,omitempty"`
	Timestamp       int64  `protobuf:"varint,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	TimestampArrive int64  `protobuf:"varint,5,opt,name=timestamp_arrive,json=timestampArrive,proto3" json:"timestamp_arrive,omitempty"`
}

func (x *MSG_TEST) Reset() {
	*x = MSG_TEST{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MSG_TEST) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MSG_TEST) ProtoMessage() {}

func (x *MSG_TEST) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MSG_TEST.ProtoReflect.Descriptor instead.
func (*MSG_TEST) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{6}
}

func (x *MSG_TEST) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *MSG_TEST) GetStr() string {
	if x != nil {
		return x.Str
	}
	return ""
}

func (x *MSG_TEST) GetSeq() int64 {
	if x != nil {
		return x.Seq
	}
	return 0
}

func (x *MSG_TEST) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *MSG_TEST) GetTimestampArrive() int64 {
	if x != nil {
		return x.TimestampArrive
	}
	return 0
}

type MSG_TEST_RES struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id               int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Str              string `protobuf:"bytes,2,opt,name=str,proto3" json:"str,omitempty"`
	Seq              int64  `protobuf:"varint,3,opt,name=seq,proto3" json:"seq,omitempty"`
	Timestamp        int64  `protobuf:"varint,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	TimestampArrive  int64  `protobuf:"varint,5,opt,name=timestamp_arrive,json=timestampArrive,proto3" json:"timestamp_arrive,omitempty"`
	TimestampProcess int64  `protobuf:"varint,6,opt,name=timestamp_process,json=timestampProcess,proto3" json:"timestamp_process,omitempty"`
}

func (x *MSG_TEST_RES) Reset() {
	*x = MSG_TEST_RES{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MSG_TEST_RES) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MSG_TEST_RES) ProtoMessage() {}

func (x *MSG_TEST_RES) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MSG_TEST_RES.ProtoReflect.Descriptor instead.
func (*MSG_TEST_RES) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{7}
}

func (x *MSG_TEST_RES) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *MSG_TEST_RES) GetStr() string {
	if x != nil {
		return x.Str
	}
	return ""
}

func (x *MSG_TEST_RES) GetSeq() int64 {
	if x != nil {
		return x.Seq
	}
	return 0
}

func (x *MSG_TEST_RES) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *MSG_TEST_RES) GetTimestampArrive() int64 {
	if x != nil {
		return x.TimestampArrive
	}
	return 0
}

func (x *MSG_TEST_RES) GetTimestampProcess() int64 {
	if x != nil {
		return x.TimestampProcess
	}
	return 0
}

type MSG_LOGIN struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *MSG_LOGIN) Reset() {
	*x = MSG_LOGIN{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MSG_LOGIN) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MSG_LOGIN) ProtoMessage() {}

func (x *MSG_LOGIN) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MSG_LOGIN.ProtoReflect.Descriptor instead.
func (*MSG_LOGIN) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{8}
}

func (x *MSG_LOGIN) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type MSG_LOGIN_RES struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	ConnectId int64 `protobuf:"varint,2,opt,name=connect_id,json=connectId,proto3" json:"connect_id,omitempty"`
}

func (x *MSG_LOGIN_RES) Reset() {
	*x = MSG_LOGIN_RES{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MSG_LOGIN_RES) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MSG_LOGIN_RES) ProtoMessage() {}

func (x *MSG_LOGIN_RES) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MSG_LOGIN_RES.ProtoReflect.Descriptor instead.
func (*MSG_LOGIN_RES) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{9}
}

func (x *MSG_LOGIN_RES) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *MSG_LOGIN_RES) GetConnectId() int64 {
	if x != nil {
		return x.ConnectId
	}
	return 0
}

type MSG_TCP_STATIC struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Seq int64 `protobuf:"varint,1,opt,name=seq,proto3" json:"seq,omitempty"`
}

func (x *MSG_TCP_STATIC) Reset() {
	*x = MSG_TCP_STATIC{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MSG_TCP_STATIC) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MSG_TCP_STATIC) ProtoMessage() {}

func (x *MSG_TCP_STATIC) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MSG_TCP_STATIC.ProtoReflect.Descriptor instead.
func (*MSG_TCP_STATIC) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{10}
}

func (x *MSG_TCP_STATIC) GetSeq() int64 {
	if x != nil {
		return x.Seq
	}
	return 0
}

type MSG_TCP_STATIC_RES struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PacketCount      int64           `protobuf:"varint,1,opt,name=packet_count,json=packetCount,proto3" json:"packet_count,omitempty"`
	ByteRecv         int64           `protobuf:"varint,2,opt,name=byte_recv,json=byteRecv,proto3" json:"byte_recv,omitempty"`
	ByteProc         int64           `protobuf:"varint,3,opt,name=byte_proc,json=byteProc,proto3" json:"byte_proc,omitempty"`
	ByteSend         int64           `protobuf:"varint,4,opt,name=byte_send,json=byteSend,proto3" json:"byte_send,omitempty"`
	MapStaticMsgRecv map[int32]int64 `protobuf:"bytes,5,rep,name=map_static_msg_recv,json=mapStaticMsgRecv,proto3" json:"map_static_msg_recv,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (x *MSG_TCP_STATIC_RES) Reset() {
	*x = MSG_TCP_STATIC_RES{}
	if protoimpl.UnsafeEnabled {
		mi := &file_msg_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MSG_TCP_STATIC_RES) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MSG_TCP_STATIC_RES) ProtoMessage() {}

func (x *MSG_TCP_STATIC_RES) ProtoReflect() protoreflect.Message {
	mi := &file_msg_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MSG_TCP_STATIC_RES.ProtoReflect.Descriptor instead.
func (*MSG_TCP_STATIC_RES) Descriptor() ([]byte, []int) {
	return file_msg_proto_rawDescGZIP(), []int{11}
}

func (x *MSG_TCP_STATIC_RES) GetPacketCount() int64 {
	if x != nil {
		return x.PacketCount
	}
	return 0
}

func (x *MSG_TCP_STATIC_RES) GetByteRecv() int64 {
	if x != nil {
		return x.ByteRecv
	}
	return 0
}

func (x *MSG_TCP_STATIC_RES) GetByteProc() int64 {
	if x != nil {
		return x.ByteProc
	}
	return 0
}

func (x *MSG_TCP_STATIC_RES) GetByteSend() int64 {
	if x != nil {
		return x.ByteSend
	}
	return 0
}

func (x *MSG_TCP_STATIC_RES) GetMapStaticMsgRecv() map[int32]int64 {
	if x != nil {
		return x.MapStaticMsgRecv
	}
	return nil
}

var File_msg_proto protoreflect.FileDescriptor

var file_msg_proto_rawDesc = []byte{
	0x0a, 0x09, 0x6d, 0x73, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x6d, 0x73, 0x67,
	0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x22, 0x54, 0x0a, 0x07, 0x4d, 0x53, 0x47, 0x5f, 0x52, 0x50,
	0x43, 0x12, 0x15, 0x0a, 0x06, 0x6d, 0x73, 0x67, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x05, 0x6d, 0x73, 0x67, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x6d, 0x73, 0x67, 0x5f,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x6d, 0x73, 0x67, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x6d, 0x73, 0x67, 0x5f, 0x62, 0x69, 0x6e, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x6d, 0x73, 0x67, 0x42, 0x69, 0x6e, 0x22, 0x8d, 0x01, 0x0a,
	0x0b, 0x4d, 0x53, 0x47, 0x5f, 0x52, 0x50, 0x43, 0x5f, 0x52, 0x45, 0x53, 0x12, 0x15, 0x0a, 0x06,
	0x6d, 0x73, 0x67, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6d, 0x73,
	0x67, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x6d, 0x73, 0x67, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x6d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x12, 0x33,
	0x0a, 0x08, 0x72, 0x65, 0x73, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x18, 0x2e, 0x6d, 0x73, 0x67, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x2e, 0x52, 0x45, 0x53,
	0x50, 0x4f, 0x4e, 0x53, 0x45, 0x5f, 0x43, 0x4f, 0x44, 0x45, 0x52, 0x07, 0x72, 0x65, 0x73, 0x43,
	0x6f, 0x64, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x6d, 0x73, 0x67, 0x5f, 0x62, 0x69, 0x6e, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x6d, 0x73, 0x67, 0x42, 0x69, 0x6e, 0x22, 0x27, 0x0a, 0x0e,
	0x4d, 0x53, 0x47, 0x5f, 0x53, 0x52, 0x56, 0x5f, 0x52, 0x45, 0x50, 0x4f, 0x52, 0x54, 0x12, 0x15,
	0x0a, 0x06, 0x73, 0x72, 0x76, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05,
	0x73, 0x72, 0x76, 0x49, 0x64, 0x22, 0x4b, 0x0a, 0x12, 0x4d, 0x53, 0x47, 0x5f, 0x53, 0x52, 0x56,
	0x5f, 0x52, 0x45, 0x50, 0x4f, 0x52, 0x54, 0x5f, 0x52, 0x45, 0x53, 0x12, 0x15, 0x0a, 0x06, 0x73,
	0x72, 0x76, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x73, 0x72, 0x76,
	0x49, 0x64, 0x12, 0x1e, 0x0a, 0x0b, 0x74, 0x63, 0x70, 0x5f, 0x63, 0x6f, 0x6e, 0x6e, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x63, 0x70, 0x43, 0x6f, 0x6e, 0x6e,
	0x49, 0x64, 0x22, 0x1f, 0x0a, 0x0d, 0x4d, 0x53, 0x47, 0x5f, 0x48, 0x45, 0x41, 0x52, 0x54, 0x42,
	0x45, 0x41, 0x54, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x02, 0x69, 0x64, 0x22, 0x23, 0x0a, 0x11, 0x4d, 0x53, 0x47, 0x5f, 0x48, 0x45, 0x41, 0x52, 0x54,
	0x42, 0x45, 0x41, 0x54, 0x5f, 0x52, 0x45, 0x53, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x22, 0x87, 0x01, 0x0a, 0x08, 0x4d, 0x53, 0x47,
	0x5f, 0x54, 0x45, 0x53, 0x54, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x74, 0x72, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x73, 0x74, 0x72, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x65, 0x71, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x73, 0x65, 0x71, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x29, 0x0a, 0x10, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x5f, 0x61, 0x72, 0x72, 0x69, 0x76, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x0f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x41, 0x72, 0x72, 0x69,
	0x76, 0x65, 0x22, 0xb8, 0x01, 0x0a, 0x0c, 0x4d, 0x53, 0x47, 0x5f, 0x54, 0x45, 0x53, 0x54, 0x5f,
	0x52, 0x45, 0x53, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x74, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x73, 0x74, 0x72, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x65, 0x71, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x03, 0x73, 0x65, 0x71, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x29, 0x0a, 0x10, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x5f, 0x61, 0x72, 0x72, 0x69, 0x76, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x41, 0x72, 0x72, 0x69, 0x76, 0x65,
	0x12, 0x2b, 0x0a, 0x11, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x5f, 0x70, 0x72,
	0x6f, 0x63, 0x65, 0x73, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x10, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x22, 0x1b, 0x0a,
	0x09, 0x4d, 0x53, 0x47, 0x5f, 0x4c, 0x4f, 0x47, 0x49, 0x4e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x22, 0x3e, 0x0a, 0x0d, 0x4d, 0x53,
	0x47, 0x5f, 0x4c, 0x4f, 0x47, 0x49, 0x4e, 0x5f, 0x52, 0x45, 0x53, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x63,
	0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x09, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x49, 0x64, 0x22, 0x22, 0x0a, 0x0e, 0x4d, 0x53,
	0x47, 0x5f, 0x54, 0x43, 0x50, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x49, 0x43, 0x12, 0x10, 0x0a, 0x03,
	0x73, 0x65, 0x71, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x73, 0x65, 0x71, 0x22, 0xb7,
	0x02, 0x0a, 0x12, 0x4d, 0x53, 0x47, 0x5f, 0x54, 0x43, 0x50, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x49,
	0x43, 0x5f, 0x52, 0x45, 0x53, 0x12, 0x21, 0x0a, 0x0c, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x5f,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x70, 0x61, 0x63,
	0x6b, 0x65, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x62, 0x79, 0x74, 0x65,
	0x5f, 0x72, 0x65, 0x63, 0x76, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x62, 0x79, 0x74,
	0x65, 0x52, 0x65, 0x63, 0x76, 0x12, 0x1b, 0x0a, 0x09, 0x62, 0x79, 0x74, 0x65, 0x5f, 0x70, 0x72,
	0x6f, 0x63, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x62, 0x79, 0x74, 0x65, 0x50, 0x72,
	0x6f, 0x63, 0x12, 0x1b, 0x0a, 0x09, 0x62, 0x79, 0x74, 0x65, 0x5f, 0x73, 0x65, 0x6e, 0x64, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x62, 0x79, 0x74, 0x65, 0x53, 0x65, 0x6e, 0x64, 0x12,
	0x62, 0x0a, 0x13, 0x6d, 0x61, 0x70, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x69, 0x63, 0x5f, 0x6d, 0x73,
	0x67, 0x5f, 0x72, 0x65, 0x63, 0x76, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x33, 0x2e, 0x6d,
	0x73, 0x67, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x2e, 0x4d, 0x53, 0x47, 0x5f, 0x54, 0x43, 0x50,
	0x5f, 0x53, 0x54, 0x41, 0x54, 0x49, 0x43, 0x5f, 0x52, 0x45, 0x53, 0x2e, 0x4d, 0x61, 0x70, 0x53,
	0x74, 0x61, 0x74, 0x69, 0x63, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x63, 0x76, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x10, 0x6d, 0x61, 0x70, 0x53, 0x74, 0x61, 0x74, 0x69, 0x63, 0x4d, 0x73, 0x67, 0x52,
	0x65, 0x63, 0x76, 0x1a, 0x43, 0x0a, 0x15, 0x4d, 0x61, 0x70, 0x53, 0x74, 0x61, 0x74, 0x69, 0x63,
	0x4d, 0x73, 0x67, 0x52, 0x65, 0x63, 0x76, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x2a, 0x95, 0x02, 0x0a, 0x08, 0x4d, 0x53, 0x47,
	0x5f, 0x54, 0x59, 0x50, 0x45, 0x12, 0x0d, 0x0a, 0x09, 0x5f, 0x4d, 0x53, 0x47, 0x5f, 0x4e, 0x55,
	0x4c, 0x4c, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x5f, 0x4d, 0x53, 0x47, 0x5f, 0x52, 0x50, 0x43,
	0x10, 0x01, 0x12, 0x10, 0x0a, 0x0c, 0x5f, 0x4d, 0x53, 0x47, 0x5f, 0x52, 0x50, 0x43, 0x5f, 0x52,
	0x45, 0x53, 0x10, 0x02, 0x12, 0x13, 0x0a, 0x0f, 0x5f, 0x4d, 0x53, 0x47, 0x5f, 0x53, 0x52, 0x56,
	0x5f, 0x52, 0x45, 0x50, 0x4f, 0x52, 0x54, 0x10, 0x03, 0x12, 0x17, 0x0a, 0x13, 0x5f, 0x4d, 0x53,
	0x47, 0x5f, 0x53, 0x52, 0x56, 0x5f, 0x52, 0x45, 0x50, 0x4f, 0x52, 0x54, 0x5f, 0x52, 0x45, 0x53,
	0x10, 0x04, 0x12, 0x12, 0x0a, 0x0e, 0x5f, 0x4d, 0x53, 0x47, 0x5f, 0x48, 0x45, 0x41, 0x52, 0x54,
	0x42, 0x45, 0x41, 0x54, 0x10, 0x05, 0x12, 0x16, 0x0a, 0x12, 0x5f, 0x4d, 0x53, 0x47, 0x5f, 0x48,
	0x45, 0x41, 0x52, 0x54, 0x42, 0x45, 0x41, 0x54, 0x5f, 0x52, 0x45, 0x53, 0x10, 0x06, 0x12, 0x13,
	0x0a, 0x0f, 0x5f, 0x4d, 0x53, 0x47, 0x5f, 0x54, 0x43, 0x50, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x49,
	0x43, 0x10, 0x07, 0x12, 0x17, 0x0a, 0x13, 0x5f, 0x4d, 0x53, 0x47, 0x5f, 0x54, 0x43, 0x50, 0x5f,
	0x53, 0x54, 0x41, 0x54, 0x49, 0x43, 0x5f, 0x52, 0x45, 0x53, 0x10, 0x08, 0x12, 0x0c, 0x0a, 0x08,
	0x5f, 0x4d, 0x53, 0x47, 0x5f, 0x4d, 0x41, 0x58, 0x10, 0x64, 0x12, 0x0d, 0x0a, 0x09, 0x5f, 0x4d,
	0x53, 0x47, 0x5f, 0x54, 0x45, 0x53, 0x54, 0x10, 0x65, 0x12, 0x11, 0x0a, 0x0d, 0x5f, 0x4d, 0x53,
	0x47, 0x5f, 0x54, 0x45, 0x53, 0x54, 0x5f, 0x52, 0x45, 0x53, 0x10, 0x66, 0x12, 0x0e, 0x0a, 0x0a,
	0x5f, 0x4d, 0x53, 0x47, 0x5f, 0x4c, 0x4f, 0x47, 0x49, 0x4e, 0x10, 0x67, 0x12, 0x12, 0x0a, 0x0e,
	0x5f, 0x4d, 0x53, 0x47, 0x5f, 0x4c, 0x4f, 0x47, 0x49, 0x4e, 0x5f, 0x52, 0x45, 0x53, 0x10, 0x68,
	0x2a, 0x61, 0x0a, 0x0d, 0x52, 0x45, 0x53, 0x50, 0x4f, 0x4e, 0x53, 0x45, 0x5f, 0x43, 0x4f, 0x44,
	0x45, 0x12, 0x16, 0x0a, 0x12, 0x52, 0x45, 0x53, 0x50, 0x4f, 0x4e, 0x53, 0x45, 0x5f, 0x43, 0x4f,
	0x44, 0x45, 0x5f, 0x4e, 0x4f, 0x4e, 0x45, 0x10, 0x00, 0x12, 0x16, 0x0a, 0x12, 0x52, 0x45, 0x53,
	0x50, 0x4f, 0x4e, 0x53, 0x45, 0x5f, 0x43, 0x4f, 0x44, 0x45, 0x5f, 0x46, 0x61, 0x69, 0x6c, 0x10,
	0x01, 0x12, 0x20, 0x0a, 0x1c, 0x52, 0x45, 0x53, 0x50, 0x4f, 0x4e, 0x53, 0x45, 0x5f, 0x43, 0x4f,
	0x44, 0x45, 0x5f, 0x52, 0x70, 0x63, 0x5f, 0x6e, 0x6f, 0x74, 0x5f, 0x61, 0x63, 0x63, 0x65, 0x70,
	0x74, 0x10, 0x02, 0x42, 0x0e, 0x5a, 0x0c, 0x2e, 0x2f, 0x3b, 0x6d, 0x73, 0x67, 0x70, 0x61, 0x63,
	0x6b, 0x65, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_msg_proto_rawDescOnce sync.Once
	file_msg_proto_rawDescData = file_msg_proto_rawDesc
)

func file_msg_proto_rawDescGZIP() []byte {
	file_msg_proto_rawDescOnce.Do(func() {
		file_msg_proto_rawDescData = protoimpl.X.CompressGZIP(file_msg_proto_rawDescData)
	})
	return file_msg_proto_rawDescData
}

var file_msg_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_msg_proto_msgTypes = make([]protoimpl.MessageInfo, 13)
var file_msg_proto_goTypes = []interface{}{
	(MSG_TYPE)(0),              // 0: msgpacket.MSG_TYPE
	(RESPONSE_CODE)(0),         // 1: msgpacket.RESPONSE_CODE
	(*MSG_RPC)(nil),            // 2: msgpacket.MSG_RPC
	(*MSG_RPC_RES)(nil),        // 3: msgpacket.MSG_RPC_RES
	(*MSG_SRV_REPORT)(nil),     // 4: msgpacket.MSG_SRV_REPORT
	(*MSG_SRV_REPORT_RES)(nil), // 5: msgpacket.MSG_SRV_REPORT_RES
	(*MSG_HEARTBEAT)(nil),      // 6: msgpacket.MSG_HEARTBEAT
	(*MSG_HEARTBEAT_RES)(nil),  // 7: msgpacket.MSG_HEARTBEAT_RES
	(*MSG_TEST)(nil),           // 8: msgpacket.MSG_TEST
	(*MSG_TEST_RES)(nil),       // 9: msgpacket.MSG_TEST_RES
	(*MSG_LOGIN)(nil),          // 10: msgpacket.MSG_LOGIN
	(*MSG_LOGIN_RES)(nil),      // 11: msgpacket.MSG_LOGIN_RES
	(*MSG_TCP_STATIC)(nil),     // 12: msgpacket.MSG_TCP_STATIC
	(*MSG_TCP_STATIC_RES)(nil), // 13: msgpacket.MSG_TCP_STATIC_RES
	nil,                        // 14: msgpacket.MSG_TCP_STATIC_RES.MapStaticMsgRecvEntry
}
var file_msg_proto_depIdxs = []int32{
	1,  // 0: msgpacket.MSG_RPC_RES.res_code:type_name -> msgpacket.RESPONSE_CODE
	14, // 1: msgpacket.MSG_TCP_STATIC_RES.map_static_msg_recv:type_name -> msgpacket.MSG_TCP_STATIC_RES.MapStaticMsgRecvEntry
	2,  // [2:2] is the sub-list for method output_type
	2,  // [2:2] is the sub-list for method input_type
	2,  // [2:2] is the sub-list for extension type_name
	2,  // [2:2] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}

func init() { file_msg_proto_init() }
func file_msg_proto_init() {
	if File_msg_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_msg_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MSG_RPC); i {
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
		file_msg_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MSG_RPC_RES); i {
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
		file_msg_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MSG_SRV_REPORT); i {
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
		file_msg_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MSG_SRV_REPORT_RES); i {
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
		file_msg_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MSG_HEARTBEAT); i {
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
		file_msg_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MSG_HEARTBEAT_RES); i {
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
		file_msg_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MSG_TEST); i {
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
		file_msg_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MSG_TEST_RES); i {
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
		file_msg_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MSG_LOGIN); i {
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
		file_msg_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MSG_LOGIN_RES); i {
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
		file_msg_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MSG_TCP_STATIC); i {
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
		file_msg_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MSG_TCP_STATIC_RES); i {
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
			RawDescriptor: file_msg_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   13,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_msg_proto_goTypes,
		DependencyIndexes: file_msg_proto_depIdxs,
		EnumInfos:         file_msg_proto_enumTypes,
		MessageInfos:      file_msg_proto_msgTypes,
	}.Build()
	File_msg_proto = out.File
	file_msg_proto_rawDesc = nil
	file_msg_proto_goTypes = nil
	file_msg_proto_depIdxs = nil
}
