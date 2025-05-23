// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.2
// source: Error.proto

package protodef

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PError_Type int32

const (
	PError_TYPE_UNKNOWN      PError_Type = 0
	PError_TYPE_SERVER_ERROR PError_Type = 1
	PError_TYPE_BUSINESS     PError_Type = 2
)

// Enum value maps for PError_Type.
var (
	PError_Type_name = map[int32]string{
		0: "TYPE_UNKNOWN",
		1: "TYPE_SERVER_ERROR",
		2: "TYPE_BUSINESS",
	}
	PError_Type_value = map[string]int32{
		"TYPE_UNKNOWN":      0,
		"TYPE_SERVER_ERROR": 1,
		"TYPE_BUSINESS":     2,
	}
)

func (x PError_Type) Enum() *PError_Type {
	p := new(PError_Type)
	*p = x
	return p
}

func (x PError_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PError_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_Error_proto_enumTypes[0].Descriptor()
}

func (PError_Type) Type() protoreflect.EnumType {
	return &file_Error_proto_enumTypes[0]
}

func (x PError_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PError_Type.Descriptor instead.
func (PError_Type) EnumDescriptor() ([]byte, []int) {
	return file_Error_proto_rawDescGZIP(), []int{0, 0}
}

type PError struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Type          PError_Type            `protobuf:"varint,1,opt,name=type,proto3,enum=DouDizhu.PError_Type" json:"type,omitempty"`
	ErrorCode     string                 `protobuf:"bytes,2,opt,name=errorCode,proto3" json:"errorCode,omitempty"`
	Message       string                 `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PError) Reset() {
	*x = PError{}
	mi := &file_Error_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PError) ProtoMessage() {}

func (x *PError) ProtoReflect() protoreflect.Message {
	mi := &file_Error_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PError.ProtoReflect.Descriptor instead.
func (*PError) Descriptor() ([]byte, []int) {
	return file_Error_proto_rawDescGZIP(), []int{0}
}

func (x *PError) GetType() PError_Type {
	if x != nil {
		return x.Type
	}
	return PError_TYPE_UNKNOWN
}

func (x *PError) GetErrorCode() string {
	if x != nil {
		return x.ErrorCode
	}
	return ""
}

func (x *PError) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_Error_proto protoreflect.FileDescriptor

const file_Error_proto_rawDesc = "" +
	"\n" +
	"\vError.proto\x12\bDouDizhu\"\xaf\x01\n" +
	"\x06PError\x12)\n" +
	"\x04type\x18\x01 \x01(\x0e2\x15.DouDizhu.PError.TypeR\x04type\x12\x1c\n" +
	"\terrorCode\x18\x02 \x01(\tR\terrorCode\x12\x18\n" +
	"\amessage\x18\x03 \x01(\tR\amessage\"B\n" +
	"\x04Type\x12\x10\n" +
	"\fTYPE_UNKNOWN\x10\x00\x12\x15\n" +
	"\x11TYPE_SERVER_ERROR\x10\x01\x12\x11\n" +
	"\rTYPE_BUSINESS\x10\x02B\"Z\x10network/protodef\xaa\x02\rNetwork.Protob\x06proto3"

var (
	file_Error_proto_rawDescOnce sync.Once
	file_Error_proto_rawDescData []byte
)

func file_Error_proto_rawDescGZIP() []byte {
	file_Error_proto_rawDescOnce.Do(func() {
		file_Error_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_Error_proto_rawDesc), len(file_Error_proto_rawDesc)))
	})
	return file_Error_proto_rawDescData
}

var file_Error_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_Error_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_Error_proto_goTypes = []any{
	(PError_Type)(0), // 0: DouDizhu.PError.Type
	(*PError)(nil),   // 1: DouDizhu.PError
}
var file_Error_proto_depIdxs = []int32{
	0, // 0: DouDizhu.PError.type:type_name -> DouDizhu.PError.Type
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_Error_proto_init() }
func file_Error_proto_init() {
	if File_Error_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_Error_proto_rawDesc), len(file_Error_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_Error_proto_goTypes,
		DependencyIndexes: file_Error_proto_depIdxs,
		EnumInfos:         file_Error_proto_enumTypes,
		MessageInfos:      file_Error_proto_msgTypes,
	}.Build()
	File_Error_proto = out.File
	file_Error_proto_goTypes = nil
	file_Error_proto_depIdxs = nil
}
