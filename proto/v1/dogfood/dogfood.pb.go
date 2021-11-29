// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.1
// source: proto/v1/dogfood/dogfood.proto

package dogfoodpb

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CreateRecordRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// dog_food name is a name of dogfood brand.
	DogfoodName string `protobuf:"bytes,1,opt,name=dogfood_name,json=dogfoodName,proto3" json:"dogfood_name,omitempty"`
	// grap specifies how grams a dog eat dogfood.
	Gram int32 `protobuf:"varint,2,opt,name=gram,proto3" json:"gram,omitempty"`
	// dog_name specifies a name of dog.
	DogName string `protobuf:"bytes,3,opt,name=dog_name,json=dogName,proto3" json:"dog_name,omitempty"`
}

func (x *CreateRecordRequest) Reset() {
	*x = CreateRecordRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_dogfood_dogfood_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateRecordRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateRecordRequest) ProtoMessage() {}

func (x *CreateRecordRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_dogfood_dogfood_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateRecordRequest.ProtoReflect.Descriptor instead.
func (*CreateRecordRequest) Descriptor() ([]byte, []int) {
	return file_proto_v1_dogfood_dogfood_proto_rawDescGZIP(), []int{0}
}

func (x *CreateRecordRequest) GetDogfoodName() string {
	if x != nil {
		return x.DogfoodName
	}
	return ""
}

func (x *CreateRecordRequest) GetGram() int32 {
	if x != nil {
		return x.Gram
	}
	return 0
}

func (x *CreateRecordRequest) GetDogName() string {
	if x != nil {
		return x.DogName
	}
	return ""
}

type ListRecordsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// from specifies the start time of eaten_at.
	From *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=from,proto3" json:"from,omitempty"`
	// page_size specifies a requested length of records.
	PageSize int32 `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	// to specifies the end time of eaten_at.
	To *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=to,proto3" json:"to,omitempty"`
}

func (x *ListRecordsRequest) Reset() {
	*x = ListRecordsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_dogfood_dogfood_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListRecordsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRecordsRequest) ProtoMessage() {}

func (x *ListRecordsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_dogfood_dogfood_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRecordsRequest.ProtoReflect.Descriptor instead.
func (*ListRecordsRequest) Descriptor() ([]byte, []int) {
	return file_proto_v1_dogfood_dogfood_proto_rawDescGZIP(), []int{1}
}

func (x *ListRecordsRequest) GetFrom() *timestamppb.Timestamp {
	if x != nil {
		return x.From
	}
	return nil
}

func (x *ListRecordsRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *ListRecordsRequest) GetTo() *timestamppb.Timestamp {
	if x != nil {
		return x.To
	}
	return nil
}

type ListRecordsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// records specify an array of Record.
	Records []*Record `protobuf:"bytes,1,rep,name=records,proto3" json:"records,omitempty"`
	// to specifies the end time of eaten_at.
	To *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=to,proto3" json:"to,omitempty"`
}

func (x *ListRecordsResponse) Reset() {
	*x = ListRecordsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_dogfood_dogfood_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListRecordsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRecordsResponse) ProtoMessage() {}

func (x *ListRecordsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_dogfood_dogfood_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRecordsResponse.ProtoReflect.Descriptor instead.
func (*ListRecordsResponse) Descriptor() ([]byte, []int) {
	return file_proto_v1_dogfood_dogfood_proto_rawDescGZIP(), []int{2}
}

func (x *ListRecordsResponse) GetRecords() []*Record {
	if x != nil {
		return x.Records
	}
	return nil
}

func (x *ListRecordsResponse) GetTo() *timestamppb.Timestamp {
	if x != nil {
		return x.To
	}
	return nil
}

type Record struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// dog_food name is a name of dogfood brand.
	DogfoodName string `protobuf:"bytes,1,opt,name=dogfood_name,json=dogfoodName,proto3" json:"dogfood_name,omitempty"`
	// grap specifies how grams a dog eat dogfood.
	Gram int32 `protobuf:"varint,2,opt,name=gram,proto3" json:"gram,omitempty"`
	// dog_name specifies a name of dog.
	DogName string `protobuf:"bytes,3,opt,name=dog_name,json=dogName,proto3" json:"dog_name,omitempty"`
	// eaten_at specifies what time a dog ate a dogfood.
	EatenAt *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=eaten_at,json=eatenAt,proto3" json:"eaten_at,omitempty"`
}

func (x *Record) Reset() {
	*x = Record{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_dogfood_dogfood_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Record) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Record) ProtoMessage() {}

func (x *Record) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_dogfood_dogfood_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Record.ProtoReflect.Descriptor instead.
func (*Record) Descriptor() ([]byte, []int) {
	return file_proto_v1_dogfood_dogfood_proto_rawDescGZIP(), []int{3}
}

func (x *Record) GetDogfoodName() string {
	if x != nil {
		return x.DogfoodName
	}
	return ""
}

func (x *Record) GetGram() int32 {
	if x != nil {
		return x.Gram
	}
	return 0
}

func (x *Record) GetDogName() string {
	if x != nil {
		return x.DogName
	}
	return ""
}

func (x *Record) GetEatenAt() *timestamppb.Timestamp {
	if x != nil {
		return x.EatenAt
	}
	return nil
}

var File_proto_v1_dogfood_dogfood_proto protoreflect.FileDescriptor

var file_proto_v1_dogfood_dogfood_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x6f, 0x67, 0x66, 0x6f,
	0x6f, 0x64, 0x2f, 0x64, 0x6f, 0x67, 0x66, 0x6f, 0x6f, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0c, 0x64, 0x6f, 0x67, 0x66, 0x6f, 0x6f, 0x64, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x1a, 0x1c,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x67, 0x0a,
	0x13, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x64, 0x6f, 0x67, 0x66, 0x6f, 0x6f, 0x64, 0x5f,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x6f, 0x67, 0x66,
	0x6f, 0x6f, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x67, 0x72, 0x61, 0x6d, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x67, 0x72, 0x61, 0x6d, 0x12, 0x19, 0x0a, 0x08, 0x64,
	0x6f, 0x67, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x64,
	0x6f, 0x67, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x8d, 0x01, 0x0a, 0x12, 0x4c, 0x69, 0x73, 0x74, 0x52,
	0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2e, 0x0a,
	0x04, 0x66, 0x72, 0x6f, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x66, 0x72, 0x6f, 0x6d, 0x12, 0x1b, 0x0a,
	0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x2a, 0x0a, 0x02, 0x74, 0x6f,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x02, 0x74, 0x6f, 0x22, 0x71, 0x0a, 0x13, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a,
	0x07, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14,
	0x2e, 0x64, 0x6f, 0x67, 0x66, 0x6f, 0x6f, 0x64, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x52, 0x07, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x12, 0x2a, 0x0a,
	0x02, 0x74, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x02, 0x74, 0x6f, 0x22, 0x91, 0x01, 0x0a, 0x06, 0x52, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x64, 0x6f, 0x67, 0x66, 0x6f, 0x6f, 0x64, 0x5f,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x6f, 0x67, 0x66,
	0x6f, 0x6f, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x67, 0x72, 0x61, 0x6d, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x67, 0x72, 0x61, 0x6d, 0x12, 0x19, 0x0a, 0x08, 0x64,
	0x6f, 0x67, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x64,
	0x6f, 0x67, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x35, 0x0a, 0x08, 0x65, 0x61, 0x74, 0x65, 0x6e, 0x5f,
	0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x65, 0x61, 0x74, 0x65, 0x6e, 0x41, 0x74, 0x32, 0xec, 0x01,
	0x0a, 0x0e, 0x44, 0x6f, 0x67, 0x46, 0x6f, 0x6f, 0x64, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x66, 0x0a, 0x0c, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64,
	0x12, 0x21, 0x2e, 0x64, 0x6f, 0x67, 0x66, 0x6f, 0x6f, 0x64, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x64, 0x6f, 0x67, 0x66, 0x6f, 0x6f, 0x64, 0x70, 0x62, 0x2e,
	0x76, 0x31, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x22, 0x1d, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x17, 0x22, 0x12, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x6f, 0x67, 0x66, 0x6f, 0x6f, 0x64, 0x2f, 0x72,
	0x65, 0x63, 0x6f, 0x72, 0x64, 0x3a, 0x01, 0x2a, 0x12, 0x72, 0x0a, 0x0b, 0x4c, 0x69, 0x73, 0x74,
	0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x12, 0x20, 0x2e, 0x64, 0x6f, 0x67, 0x66, 0x6f, 0x6f,
	0x64, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x63, 0x6f, 0x72,
	0x64, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x64, 0x6f, 0x67, 0x66,
	0x6f, 0x6f, 0x64, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x63,
	0x6f, 0x72, 0x64, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1e, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x18, 0x22, 0x13, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x6f, 0x67, 0x66, 0x6f, 0x6f,
	0x64, 0x2f, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x3a, 0x01, 0x2a, 0x42, 0x14, 0x5a, 0x12,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x6f, 0x67, 0x66, 0x6f, 0x6f, 0x64,
	0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_v1_dogfood_dogfood_proto_rawDescOnce sync.Once
	file_proto_v1_dogfood_dogfood_proto_rawDescData = file_proto_v1_dogfood_dogfood_proto_rawDesc
)

func file_proto_v1_dogfood_dogfood_proto_rawDescGZIP() []byte {
	file_proto_v1_dogfood_dogfood_proto_rawDescOnce.Do(func() {
		file_proto_v1_dogfood_dogfood_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_v1_dogfood_dogfood_proto_rawDescData)
	})
	return file_proto_v1_dogfood_dogfood_proto_rawDescData
}

var file_proto_v1_dogfood_dogfood_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_v1_dogfood_dogfood_proto_goTypes = []interface{}{
	(*CreateRecordRequest)(nil),   // 0: dogfoodpb.v1.CreateRecordRequest
	(*ListRecordsRequest)(nil),    // 1: dogfoodpb.v1.ListRecordsRequest
	(*ListRecordsResponse)(nil),   // 2: dogfoodpb.v1.ListRecordsResponse
	(*Record)(nil),                // 3: dogfoodpb.v1.Record
	(*timestamppb.Timestamp)(nil), // 4: google.protobuf.Timestamp
}
var file_proto_v1_dogfood_dogfood_proto_depIdxs = []int32{
	4, // 0: dogfoodpb.v1.ListRecordsRequest.from:type_name -> google.protobuf.Timestamp
	4, // 1: dogfoodpb.v1.ListRecordsRequest.to:type_name -> google.protobuf.Timestamp
	3, // 2: dogfoodpb.v1.ListRecordsResponse.records:type_name -> dogfoodpb.v1.Record
	4, // 3: dogfoodpb.v1.ListRecordsResponse.to:type_name -> google.protobuf.Timestamp
	4, // 4: dogfoodpb.v1.Record.eaten_at:type_name -> google.protobuf.Timestamp
	0, // 5: dogfoodpb.v1.DogFoodService.CreateRecord:input_type -> dogfoodpb.v1.CreateRecordRequest
	1, // 6: dogfoodpb.v1.DogFoodService.ListRecords:input_type -> dogfoodpb.v1.ListRecordsRequest
	3, // 7: dogfoodpb.v1.DogFoodService.CreateRecord:output_type -> dogfoodpb.v1.Record
	2, // 8: dogfoodpb.v1.DogFoodService.ListRecords:output_type -> dogfoodpb.v1.ListRecordsResponse
	7, // [7:9] is the sub-list for method output_type
	5, // [5:7] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_proto_v1_dogfood_dogfood_proto_init() }
func file_proto_v1_dogfood_dogfood_proto_init() {
	if File_proto_v1_dogfood_dogfood_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_v1_dogfood_dogfood_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateRecordRequest); i {
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
		file_proto_v1_dogfood_dogfood_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListRecordsRequest); i {
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
		file_proto_v1_dogfood_dogfood_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListRecordsResponse); i {
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
		file_proto_v1_dogfood_dogfood_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Record); i {
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
			RawDescriptor: file_proto_v1_dogfood_dogfood_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_v1_dogfood_dogfood_proto_goTypes,
		DependencyIndexes: file_proto_v1_dogfood_dogfood_proto_depIdxs,
		MessageInfos:      file_proto_v1_dogfood_dogfood_proto_msgTypes,
	}.Build()
	File_proto_v1_dogfood_dogfood_proto = out.File
	file_proto_v1_dogfood_dogfood_proto_rawDesc = nil
	file_proto_v1_dogfood_dogfood_proto_goTypes = nil
	file_proto_v1_dogfood_dogfood_proto_depIdxs = nil
}
