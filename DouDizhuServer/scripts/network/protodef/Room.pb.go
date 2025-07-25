// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.2
// source: Room.proto

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

type PRoomState int32

const (
	PRoomState_PROOM_STATE_UNKNOWN PRoomState = 0
	PRoomState_PROOM_STATE_WAITING PRoomState = 1
	PRoomState_PROOM_STATE_PLAYING PRoomState = 2
)

// Enum value maps for PRoomState.
var (
	PRoomState_name = map[int32]string{
		0: "PROOM_STATE_UNKNOWN",
		1: "PROOM_STATE_WAITING",
		2: "PROOM_STATE_PLAYING",
	}
	PRoomState_value = map[string]int32{
		"PROOM_STATE_UNKNOWN": 0,
		"PROOM_STATE_WAITING": 1,
		"PROOM_STATE_PLAYING": 2,
	}
)

func (x PRoomState) Enum() *PRoomState {
	p := new(PRoomState)
	*p = x
	return p
}

func (x PRoomState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PRoomState) Descriptor() protoreflect.EnumDescriptor {
	return file_Room_proto_enumTypes[0].Descriptor()
}

func (PRoomState) Type() protoreflect.EnumType {
	return &file_Room_proto_enumTypes[0]
}

func (x PRoomState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PRoomState.Descriptor instead.
func (PRoomState) EnumDescriptor() ([]byte, []int) {
	return file_Room_proto_rawDescGZIP(), []int{0}
}

type PRoom struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	Id             uint32                 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name           string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Owner          *PPlayer               `protobuf:"bytes,3,opt,name=owner,proto3" json:"owner,omitempty"`
	Players        []*PPlayer             `protobuf:"bytes,4,rep,name=players,proto3" json:"players,omitempty"`
	MaxPlayerCount uint32                 `protobuf:"varint,5,opt,name=maxPlayerCount,proto3" json:"maxPlayerCount,omitempty"`
	State          PRoomState             `protobuf:"varint,6,opt,name=state,proto3,enum=DouDizhu.PRoomState" json:"state,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *PRoom) Reset() {
	*x = PRoom{}
	mi := &file_Room_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PRoom) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PRoom) ProtoMessage() {}

func (x *PRoom) ProtoReflect() protoreflect.Message {
	mi := &file_Room_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PRoom.ProtoReflect.Descriptor instead.
func (*PRoom) Descriptor() ([]byte, []int) {
	return file_Room_proto_rawDescGZIP(), []int{0}
}

func (x *PRoom) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *PRoom) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *PRoom) GetOwner() *PPlayer {
	if x != nil {
		return x.Owner
	}
	return nil
}

func (x *PRoom) GetPlayers() []*PPlayer {
	if x != nil {
		return x.Players
	}
	return nil
}

func (x *PRoom) GetMaxPlayerCount() uint32 {
	if x != nil {
		return x.MaxPlayerCount
	}
	return 0
}

func (x *PRoom) GetState() PRoomState {
	if x != nil {
		return x.State
	}
	return PRoomState_PROOM_STATE_UNKNOWN
}

type PCreateRoomRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	RoomName      string                 `protobuf:"bytes,1,opt,name=roomName,proto3" json:"roomName,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PCreateRoomRequest) Reset() {
	*x = PCreateRoomRequest{}
	mi := &file_Room_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PCreateRoomRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PCreateRoomRequest) ProtoMessage() {}

func (x *PCreateRoomRequest) ProtoReflect() protoreflect.Message {
	mi := &file_Room_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PCreateRoomRequest.ProtoReflect.Descriptor instead.
func (*PCreateRoomRequest) Descriptor() ([]byte, []int) {
	return file_Room_proto_rawDescGZIP(), []int{1}
}

func (x *PCreateRoomRequest) GetRoomName() string {
	if x != nil {
		return x.RoomName
	}
	return ""
}

type PCreateRoomResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Room          *PRoom                 `protobuf:"bytes,1,opt,name=room,proto3" json:"room,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PCreateRoomResponse) Reset() {
	*x = PCreateRoomResponse{}
	mi := &file_Room_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PCreateRoomResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PCreateRoomResponse) ProtoMessage() {}

func (x *PCreateRoomResponse) ProtoReflect() protoreflect.Message {
	mi := &file_Room_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PCreateRoomResponse.ProtoReflect.Descriptor instead.
func (*PCreateRoomResponse) Descriptor() ([]byte, []int) {
	return file_Room_proto_rawDescGZIP(), []int{2}
}

func (x *PCreateRoomResponse) GetRoom() *PRoom {
	if x != nil {
		return x.Room
	}
	return nil
}

type PEnterRoomRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	RoomId        uint32                 `protobuf:"varint,1,opt,name=roomId,proto3" json:"roomId,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PEnterRoomRequest) Reset() {
	*x = PEnterRoomRequest{}
	mi := &file_Room_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PEnterRoomRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PEnterRoomRequest) ProtoMessage() {}

func (x *PEnterRoomRequest) ProtoReflect() protoreflect.Message {
	mi := &file_Room_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PEnterRoomRequest.ProtoReflect.Descriptor instead.
func (*PEnterRoomRequest) Descriptor() ([]byte, []int) {
	return file_Room_proto_rawDescGZIP(), []int{3}
}

func (x *PEnterRoomRequest) GetRoomId() uint32 {
	if x != nil {
		return x.RoomId
	}
	return 0
}

type PEnterRoomResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Room          *PRoom                 `protobuf:"bytes,1,opt,name=room,proto3" json:"room,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PEnterRoomResponse) Reset() {
	*x = PEnterRoomResponse{}
	mi := &file_Room_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PEnterRoomResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PEnterRoomResponse) ProtoMessage() {}

func (x *PEnterRoomResponse) ProtoReflect() protoreflect.Message {
	mi := &file_Room_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PEnterRoomResponse.ProtoReflect.Descriptor instead.
func (*PEnterRoomResponse) Descriptor() ([]byte, []int) {
	return file_Room_proto_rawDescGZIP(), []int{4}
}

func (x *PEnterRoomResponse) GetRoom() *PRoom {
	if x != nil {
		return x.Room
	}
	return nil
}

type PLeaveRoomRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PLeaveRoomRequest) Reset() {
	*x = PLeaveRoomRequest{}
	mi := &file_Room_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PLeaveRoomRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PLeaveRoomRequest) ProtoMessage() {}

func (x *PLeaveRoomRequest) ProtoReflect() protoreflect.Message {
	mi := &file_Room_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PLeaveRoomRequest.ProtoReflect.Descriptor instead.
func (*PLeaveRoomRequest) Descriptor() ([]byte, []int) {
	return file_Room_proto_rawDescGZIP(), []int{5}
}

type PGetRoomListRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PGetRoomListRequest) Reset() {
	*x = PGetRoomListRequest{}
	mi := &file_Room_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PGetRoomListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PGetRoomListRequest) ProtoMessage() {}

func (x *PGetRoomListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_Room_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PGetRoomListRequest.ProtoReflect.Descriptor instead.
func (*PGetRoomListRequest) Descriptor() ([]byte, []int) {
	return file_Room_proto_rawDescGZIP(), []int{6}
}

type PGetRoomListResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Rooms         []*PRoom               `protobuf:"bytes,1,rep,name=rooms,proto3" json:"rooms,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PGetRoomListResponse) Reset() {
	*x = PGetRoomListResponse{}
	mi := &file_Room_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PGetRoomListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PGetRoomListResponse) ProtoMessage() {}

func (x *PGetRoomListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_Room_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PGetRoomListResponse.ProtoReflect.Descriptor instead.
func (*PGetRoomListResponse) Descriptor() ([]byte, []int) {
	return file_Room_proto_rawDescGZIP(), []int{7}
}

func (x *PGetRoomListResponse) GetRooms() []*PRoom {
	if x != nil {
		return x.Rooms
	}
	return nil
}

type PSyncRoomInfoRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	RoomId        uint32                 `protobuf:"varint,1,opt,name=roomId,proto3" json:"roomId,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PSyncRoomInfoRequest) Reset() {
	*x = PSyncRoomInfoRequest{}
	mi := &file_Room_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PSyncRoomInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PSyncRoomInfoRequest) ProtoMessage() {}

func (x *PSyncRoomInfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_Room_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PSyncRoomInfoRequest.ProtoReflect.Descriptor instead.
func (*PSyncRoomInfoRequest) Descriptor() ([]byte, []int) {
	return file_Room_proto_rawDescGZIP(), []int{8}
}

func (x *PSyncRoomInfoRequest) GetRoomId() uint32 {
	if x != nil {
		return x.RoomId
	}
	return 0
}

type PSyncRoomInfoResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Room          *PRoom                 `protobuf:"bytes,1,opt,name=room,proto3" json:"room,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PSyncRoomInfoResponse) Reset() {
	*x = PSyncRoomInfoResponse{}
	mi := &file_Room_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PSyncRoomInfoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PSyncRoomInfoResponse) ProtoMessage() {}

func (x *PSyncRoomInfoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_Room_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PSyncRoomInfoResponse.ProtoReflect.Descriptor instead.
func (*PSyncRoomInfoResponse) Descriptor() ([]byte, []int) {
	return file_Room_proto_rawDescGZIP(), []int{9}
}

func (x *PSyncRoomInfoResponse) GetRoom() *PRoom {
	if x != nil {
		return x.Room
	}
	return nil
}

type PRoomDisbandedNotification struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PRoomDisbandedNotification) Reset() {
	*x = PRoomDisbandedNotification{}
	mi := &file_Room_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PRoomDisbandedNotification) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PRoomDisbandedNotification) ProtoMessage() {}

func (x *PRoomDisbandedNotification) ProtoReflect() protoreflect.Message {
	mi := &file_Room_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PRoomDisbandedNotification.ProtoReflect.Descriptor instead.
func (*PRoomDisbandedNotification) Descriptor() ([]byte, []int) {
	return file_Room_proto_rawDescGZIP(), []int{10}
}

type PRoomChangedNotification struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Room          *PRoom                 `protobuf:"bytes,1,opt,name=room,proto3" json:"room,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PRoomChangedNotification) Reset() {
	*x = PRoomChangedNotification{}
	mi := &file_Room_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PRoomChangedNotification) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PRoomChangedNotification) ProtoMessage() {}

func (x *PRoomChangedNotification) ProtoReflect() protoreflect.Message {
	mi := &file_Room_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PRoomChangedNotification.ProtoReflect.Descriptor instead.
func (*PRoomChangedNotification) Descriptor() ([]byte, []int) {
	return file_Room_proto_rawDescGZIP(), []int{11}
}

func (x *PRoomChangedNotification) GetRoom() *PRoom {
	if x != nil {
		return x.Room
	}
	return nil
}

var File_Room_proto protoreflect.FileDescriptor

const file_Room_proto_rawDesc = "" +
	"\n" +
	"\n" +
	"Room.proto\x12\bDouDizhu\x1a\fPlayer.proto\"\xd5\x01\n" +
	"\x05PRoom\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\rR\x02id\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12'\n" +
	"\x05owner\x18\x03 \x01(\v2\x11.DouDizhu.PPlayerR\x05owner\x12+\n" +
	"\aplayers\x18\x04 \x03(\v2\x11.DouDizhu.PPlayerR\aplayers\x12&\n" +
	"\x0emaxPlayerCount\x18\x05 \x01(\rR\x0emaxPlayerCount\x12*\n" +
	"\x05state\x18\x06 \x01(\x0e2\x14.DouDizhu.PRoomStateR\x05state\"0\n" +
	"\x12PCreateRoomRequest\x12\x1a\n" +
	"\broomName\x18\x01 \x01(\tR\broomName\":\n" +
	"\x13PCreateRoomResponse\x12#\n" +
	"\x04room\x18\x01 \x01(\v2\x0f.DouDizhu.PRoomR\x04room\"+\n" +
	"\x11PEnterRoomRequest\x12\x16\n" +
	"\x06roomId\x18\x01 \x01(\rR\x06roomId\"9\n" +
	"\x12PEnterRoomResponse\x12#\n" +
	"\x04room\x18\x01 \x01(\v2\x0f.DouDizhu.PRoomR\x04room\"\x13\n" +
	"\x11PLeaveRoomRequest\"\x15\n" +
	"\x13PGetRoomListRequest\"=\n" +
	"\x14PGetRoomListResponse\x12%\n" +
	"\x05rooms\x18\x01 \x03(\v2\x0f.DouDizhu.PRoomR\x05rooms\".\n" +
	"\x14PSyncRoomInfoRequest\x12\x16\n" +
	"\x06roomId\x18\x01 \x01(\rR\x06roomId\"<\n" +
	"\x15PSyncRoomInfoResponse\x12#\n" +
	"\x04room\x18\x01 \x01(\v2\x0f.DouDizhu.PRoomR\x04room\"\x1c\n" +
	"\x1aPRoomDisbandedNotification\"?\n" +
	"\x18PRoomChangedNotification\x12#\n" +
	"\x04room\x18\x01 \x01(\v2\x0f.DouDizhu.PRoomR\x04room*W\n" +
	"\n" +
	"PRoomState\x12\x17\n" +
	"\x13PROOM_STATE_UNKNOWN\x10\x00\x12\x17\n" +
	"\x13PROOM_STATE_WAITING\x10\x01\x12\x17\n" +
	"\x13PROOM_STATE_PLAYING\x10\x02B\"Z\x10network/protodef\xaa\x02\rNetwork.Protob\x06proto3"

var (
	file_Room_proto_rawDescOnce sync.Once
	file_Room_proto_rawDescData []byte
)

func file_Room_proto_rawDescGZIP() []byte {
	file_Room_proto_rawDescOnce.Do(func() {
		file_Room_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_Room_proto_rawDesc), len(file_Room_proto_rawDesc)))
	})
	return file_Room_proto_rawDescData
}

var file_Room_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_Room_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_Room_proto_goTypes = []any{
	(PRoomState)(0),                    // 0: DouDizhu.PRoomState
	(*PRoom)(nil),                      // 1: DouDizhu.PRoom
	(*PCreateRoomRequest)(nil),         // 2: DouDizhu.PCreateRoomRequest
	(*PCreateRoomResponse)(nil),        // 3: DouDizhu.PCreateRoomResponse
	(*PEnterRoomRequest)(nil),          // 4: DouDizhu.PEnterRoomRequest
	(*PEnterRoomResponse)(nil),         // 5: DouDizhu.PEnterRoomResponse
	(*PLeaveRoomRequest)(nil),          // 6: DouDizhu.PLeaveRoomRequest
	(*PGetRoomListRequest)(nil),        // 7: DouDizhu.PGetRoomListRequest
	(*PGetRoomListResponse)(nil),       // 8: DouDizhu.PGetRoomListResponse
	(*PSyncRoomInfoRequest)(nil),       // 9: DouDizhu.PSyncRoomInfoRequest
	(*PSyncRoomInfoResponse)(nil),      // 10: DouDizhu.PSyncRoomInfoResponse
	(*PRoomDisbandedNotification)(nil), // 11: DouDizhu.PRoomDisbandedNotification
	(*PRoomChangedNotification)(nil),   // 12: DouDizhu.PRoomChangedNotification
	(*PPlayer)(nil),                    // 13: DouDizhu.PPlayer
}
var file_Room_proto_depIdxs = []int32{
	13, // 0: DouDizhu.PRoom.owner:type_name -> DouDizhu.PPlayer
	13, // 1: DouDizhu.PRoom.players:type_name -> DouDizhu.PPlayer
	0,  // 2: DouDizhu.PRoom.state:type_name -> DouDizhu.PRoomState
	1,  // 3: DouDizhu.PCreateRoomResponse.room:type_name -> DouDizhu.PRoom
	1,  // 4: DouDizhu.PEnterRoomResponse.room:type_name -> DouDizhu.PRoom
	1,  // 5: DouDizhu.PGetRoomListResponse.rooms:type_name -> DouDizhu.PRoom
	1,  // 6: DouDizhu.PSyncRoomInfoResponse.room:type_name -> DouDizhu.PRoom
	1,  // 7: DouDizhu.PRoomChangedNotification.room:type_name -> DouDizhu.PRoom
	8,  // [8:8] is the sub-list for method output_type
	8,  // [8:8] is the sub-list for method input_type
	8,  // [8:8] is the sub-list for extension type_name
	8,  // [8:8] is the sub-list for extension extendee
	0,  // [0:8] is the sub-list for field type_name
}

func init() { file_Room_proto_init() }
func file_Room_proto_init() {
	if File_Room_proto != nil {
		return
	}
	file_Player_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_Room_proto_rawDesc), len(file_Room_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_Room_proto_goTypes,
		DependencyIndexes: file_Room_proto_depIdxs,
		EnumInfos:         file_Room_proto_enumTypes,
		MessageInfos:      file_Room_proto_msgTypes,
	}.Build()
	File_Room_proto = out.File
	file_Room_proto_goTypes = nil
	file_Room_proto_depIdxs = nil
}
