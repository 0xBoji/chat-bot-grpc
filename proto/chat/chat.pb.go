// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: proto/chat/chat.proto

package chat

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

// Request to send a message
type SendMessageRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Content       string                 `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
	SenderId      int64                  `protobuf:"varint,2,opt,name=sender_id,json=senderId,proto3" json:"sender_id,omitempty"`
	RoomId        int64                  `protobuf:"varint,3,opt,name=room_id,json=roomId,proto3" json:"room_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SendMessageRequest) Reset() {
	*x = SendMessageRequest{}
	mi := &file_proto_chat_chat_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SendMessageRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendMessageRequest) ProtoMessage() {}

func (x *SendMessageRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chat_chat_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendMessageRequest.ProtoReflect.Descriptor instead.
func (*SendMessageRequest) Descriptor() ([]byte, []int) {
	return file_proto_chat_chat_proto_rawDescGZIP(), []int{0}
}

func (x *SendMessageRequest) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *SendMessageRequest) GetSenderId() int64 {
	if x != nil {
		return x.SenderId
	}
	return 0
}

func (x *SendMessageRequest) GetRoomId() int64 {
	if x != nil {
		return x.RoomId
	}
	return 0
}

// Response to a send message request
type SendMessageResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	MessageId     int64                  `protobuf:"varint,3,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SendMessageResponse) Reset() {
	*x = SendMessageResponse{}
	mi := &file_proto_chat_chat_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SendMessageResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SendMessageResponse) ProtoMessage() {}

func (x *SendMessageResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chat_chat_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SendMessageResponse.ProtoReflect.Descriptor instead.
func (*SendMessageResponse) Descriptor() ([]byte, []int) {
	return file_proto_chat_chat_proto_rawDescGZIP(), []int{1}
}

func (x *SendMessageResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *SendMessageResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *SendMessageResponse) GetMessageId() int64 {
	if x != nil {
		return x.MessageId
	}
	return 0
}

// Request to get messages from a room
type GetRoomMessagesRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	RoomId        int64                  `protobuf:"varint,1,opt,name=room_id,json=roomId,proto3" json:"room_id,omitempty"`
	UserId        int64                  `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Limit         int64                  `protobuf:"varint,3,opt,name=limit,proto3" json:"limit,omitempty"`
	Offset        int64                  `protobuf:"varint,4,opt,name=offset,proto3" json:"offset,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetRoomMessagesRequest) Reset() {
	*x = GetRoomMessagesRequest{}
	mi := &file_proto_chat_chat_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetRoomMessagesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRoomMessagesRequest) ProtoMessage() {}

func (x *GetRoomMessagesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chat_chat_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRoomMessagesRequest.ProtoReflect.Descriptor instead.
func (*GetRoomMessagesRequest) Descriptor() ([]byte, []int) {
	return file_proto_chat_chat_proto_rawDescGZIP(), []int{2}
}

func (x *GetRoomMessagesRequest) GetRoomId() int64 {
	if x != nil {
		return x.RoomId
	}
	return 0
}

func (x *GetRoomMessagesRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *GetRoomMessagesRequest) GetLimit() int64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *GetRoomMessagesRequest) GetOffset() int64 {
	if x != nil {
		return x.Offset
	}
	return 0
}

// Response to a get messages request
type GetRoomMessagesResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Messages      []*MessageResponse     `protobuf:"bytes,1,rep,name=messages,proto3" json:"messages,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetRoomMessagesResponse) Reset() {
	*x = GetRoomMessagesResponse{}
	mi := &file_proto_chat_chat_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetRoomMessagesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRoomMessagesResponse) ProtoMessage() {}

func (x *GetRoomMessagesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chat_chat_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRoomMessagesResponse.ProtoReflect.Descriptor instead.
func (*GetRoomMessagesResponse) Descriptor() ([]byte, []int) {
	return file_proto_chat_chat_proto_rawDescGZIP(), []int{3}
}

func (x *GetRoomMessagesResponse) GetMessages() []*MessageResponse {
	if x != nil {
		return x.Messages
	}
	return nil
}

// Request to stream messages from a room
type StreamRoomMessagesRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	RoomId        int64                  `protobuf:"varint,1,opt,name=room_id,json=roomId,proto3" json:"room_id,omitempty"`
	UserId        int64                  `protobuf:"varint,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StreamRoomMessagesRequest) Reset() {
	*x = StreamRoomMessagesRequest{}
	mi := &file_proto_chat_chat_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StreamRoomMessagesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamRoomMessagesRequest) ProtoMessage() {}

func (x *StreamRoomMessagesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chat_chat_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamRoomMessagesRequest.ProtoReflect.Descriptor instead.
func (*StreamRoomMessagesRequest) Descriptor() ([]byte, []int) {
	return file_proto_chat_chat_proto_rawDescGZIP(), []int{4}
}

func (x *StreamRoomMessagesRequest) GetRoomId() int64 {
	if x != nil {
		return x.RoomId
	}
	return 0
}

func (x *StreamRoomMessagesRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

// Message response
type MessageResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Content       string                 `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
	SenderId      int64                  `protobuf:"varint,3,opt,name=sender_id,json=senderId,proto3" json:"sender_id,omitempty"`
	RoomId        int64                  `protobuf:"varint,4,opt,name=room_id,json=roomId,proto3" json:"room_id,omitempty"`
	SenderName    string                 `protobuf:"bytes,5,opt,name=sender_name,json=senderName,proto3" json:"sender_name,omitempty"`
	Timestamp     string                 `protobuf:"bytes,6,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MessageResponse) Reset() {
	*x = MessageResponse{}
	mi := &file_proto_chat_chat_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MessageResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageResponse) ProtoMessage() {}

func (x *MessageResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_chat_chat_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessageResponse.ProtoReflect.Descriptor instead.
func (*MessageResponse) Descriptor() ([]byte, []int) {
	return file_proto_chat_chat_proto_rawDescGZIP(), []int{5}
}

func (x *MessageResponse) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *MessageResponse) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *MessageResponse) GetSenderId() int64 {
	if x != nil {
		return x.SenderId
	}
	return 0
}

func (x *MessageResponse) GetRoomId() int64 {
	if x != nil {
		return x.RoomId
	}
	return 0
}

func (x *MessageResponse) GetSenderName() string {
	if x != nil {
		return x.SenderName
	}
	return ""
}

func (x *MessageResponse) GetTimestamp() string {
	if x != nil {
		return x.Timestamp
	}
	return ""
}

var File_proto_chat_chat_proto protoreflect.FileDescriptor

const file_proto_chat_chat_proto_rawDesc = "" +
	"\n" +
	"\x15proto/chat/chat.proto\x12\x04chat\x1a\x1cgoogle/api/annotations.proto\"d\n" +
	"\x12SendMessageRequest\x12\x18\n" +
	"\acontent\x18\x01 \x01(\tR\acontent\x12\x1b\n" +
	"\tsender_id\x18\x02 \x01(\x03R\bsenderId\x12\x17\n" +
	"\aroom_id\x18\x03 \x01(\x03R\x06roomId\"h\n" +
	"\x13SendMessageResponse\x12\x18\n" +
	"\asuccess\x18\x01 \x01(\bR\asuccess\x12\x18\n" +
	"\amessage\x18\x02 \x01(\tR\amessage\x12\x1d\n" +
	"\n" +
	"message_id\x18\x03 \x01(\x03R\tmessageId\"x\n" +
	"\x16GetRoomMessagesRequest\x12\x17\n" +
	"\aroom_id\x18\x01 \x01(\x03R\x06roomId\x12\x17\n" +
	"\auser_id\x18\x02 \x01(\x03R\x06userId\x12\x14\n" +
	"\x05limit\x18\x03 \x01(\x03R\x05limit\x12\x16\n" +
	"\x06offset\x18\x04 \x01(\x03R\x06offset\"L\n" +
	"\x17GetRoomMessagesResponse\x121\n" +
	"\bmessages\x18\x01 \x03(\v2\x15.chat.MessageResponseR\bmessages\"M\n" +
	"\x19StreamRoomMessagesRequest\x12\x17\n" +
	"\aroom_id\x18\x01 \x01(\x03R\x06roomId\x12\x17\n" +
	"\auser_id\x18\x02 \x01(\x03R\x06userId\"\xb0\x01\n" +
	"\x0fMessageResponse\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\x03R\x02id\x12\x18\n" +
	"\acontent\x18\x02 \x01(\tR\acontent\x12\x1b\n" +
	"\tsender_id\x18\x03 \x01(\x03R\bsenderId\x12\x17\n" +
	"\aroom_id\x18\x04 \x01(\x03R\x06roomId\x12\x1f\n" +
	"\vsender_name\x18\x05 \x01(\tR\n" +
	"senderName\x12\x1c\n" +
	"\ttimestamp\x18\x06 \x01(\tR\ttimestamp2\xb6\x02\n" +
	"\vChatService\x12a\n" +
	"\vSendMessage\x12\x18.chat.SendMessageRequest\x1a\x19.chat.SendMessageResponse\"\x1d\x82\xd3\xe4\x93\x02\x17:\x01*\"\x12/chat/send-message\x12r\n" +
	"\x0fGetRoomMessages\x12\x1c.chat.GetRoomMessagesRequest\x1a\x1d.chat.GetRoomMessagesResponse\"\"\x82\xd3\xe4\x93\x02\x1c:\x01*\"\x17/chat/get-room-messages\x12P\n" +
	"\x12StreamRoomMessages\x12\x1f.chat.StreamRoomMessagesRequest\x1a\x15.chat.MessageResponse\"\x000\x01B Z\x1egrpc-messenger-core/proto/chatb\x06proto3"

var (
	file_proto_chat_chat_proto_rawDescOnce sync.Once
	file_proto_chat_chat_proto_rawDescData []byte
)

func file_proto_chat_chat_proto_rawDescGZIP() []byte {
	file_proto_chat_chat_proto_rawDescOnce.Do(func() {
		file_proto_chat_chat_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_chat_chat_proto_rawDesc), len(file_proto_chat_chat_proto_rawDesc)))
	})
	return file_proto_chat_chat_proto_rawDescData
}

var file_proto_chat_chat_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_proto_chat_chat_proto_goTypes = []any{
	(*SendMessageRequest)(nil),        // 0: chat.SendMessageRequest
	(*SendMessageResponse)(nil),       // 1: chat.SendMessageResponse
	(*GetRoomMessagesRequest)(nil),    // 2: chat.GetRoomMessagesRequest
	(*GetRoomMessagesResponse)(nil),   // 3: chat.GetRoomMessagesResponse
	(*StreamRoomMessagesRequest)(nil), // 4: chat.StreamRoomMessagesRequest
	(*MessageResponse)(nil),           // 5: chat.MessageResponse
}
var file_proto_chat_chat_proto_depIdxs = []int32{
	5, // 0: chat.GetRoomMessagesResponse.messages:type_name -> chat.MessageResponse
	0, // 1: chat.ChatService.SendMessage:input_type -> chat.SendMessageRequest
	2, // 2: chat.ChatService.GetRoomMessages:input_type -> chat.GetRoomMessagesRequest
	4, // 3: chat.ChatService.StreamRoomMessages:input_type -> chat.StreamRoomMessagesRequest
	1, // 4: chat.ChatService.SendMessage:output_type -> chat.SendMessageResponse
	3, // 5: chat.ChatService.GetRoomMessages:output_type -> chat.GetRoomMessagesResponse
	5, // 6: chat.ChatService.StreamRoomMessages:output_type -> chat.MessageResponse
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_chat_chat_proto_init() }
func file_proto_chat_chat_proto_init() {
	if File_proto_chat_chat_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_chat_chat_proto_rawDesc), len(file_proto_chat_chat_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_chat_chat_proto_goTypes,
		DependencyIndexes: file_proto_chat_chat_proto_depIdxs,
		MessageInfos:      file_proto_chat_chat_proto_msgTypes,
	}.Build()
	File_proto_chat_chat_proto = out.File
	file_proto_chat_chat_proto_goTypes = nil
	file_proto_chat_chat_proto_depIdxs = nil
}
