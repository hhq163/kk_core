// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.21.0
// 	protoc        v3.11.3
// source: example.proto

package protocol

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// 定义路由 CMD 命令字
type Cmd int32

const (
	Cmd_CBeat      Cmd = 0 // 客户端发送心跳
	Cmd_SBeat      Cmd = 1 // 服务端心跳响应
	Cmd_CSecretKey Cmd = 2 // 客户端请求秘钥
	Cmd_SSecretKey Cmd = 3 // 服务端响应秘钥
	Cmd_CLogin     Cmd = 4 // 客户端登录请求
	Cmd_SLogin     Cmd = 5 // 服务端响应请求
)

// Enum value maps for Cmd.
var (
	Cmd_name = map[int32]string{
		0: "CBeat",
		1: "SBeat",
		2: "CSecretKey",
		3: "SSecretKey",
		4: "CLogin",
		5: "SLogin",
	}
	Cmd_value = map[string]int32{
		"CBeat":      0,
		"SBeat":      1,
		"CSecretKey": 2,
		"SSecretKey": 3,
		"CLogin":     4,
		"SLogin":     5,
	}
)

func (x Cmd) Enum() *Cmd {
	p := new(Cmd)
	*p = x
	return p
}

func (x Cmd) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Cmd) Descriptor() protoreflect.EnumDescriptor {
	return file_example_proto_enumTypes[0].Descriptor()
}

func (Cmd) Type() protoreflect.EnumType {
	return &file_example_proto_enumTypes[0]
}

func (x Cmd) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Cmd.Descriptor instead.
func (Cmd) EnumDescriptor() ([]byte, []int) {
	return file_example_proto_rawDescGZIP(), []int{0}
}

// 服务端响应状态码
type Result int32

const (
	Result_Result_SUCC           Result = 0
	Result_Result_FAIL           Result = 1
	Result_Result_USER_NOT_EXIST Result = 2
	Result_Result_PARAM_ERROR    Result = 3
	Result_Result_DECRYPT_ERROR  Result = 4
)

// Enum value maps for Result.
var (
	Result_name = map[int32]string{
		0: "Result_SUCC",
		1: "Result_FAIL",
		2: "Result_USER_NOT_EXIST",
		3: "Result_PARAM_ERROR",
		4: "Result_DECRYPT_ERROR",
	}
	Result_value = map[string]int32{
		"Result_SUCC":           0,
		"Result_FAIL":           1,
		"Result_USER_NOT_EXIST": 2,
		"Result_PARAM_ERROR":    3,
		"Result_DECRYPT_ERROR":  4,
	}
)

func (x Result) Enum() *Result {
	p := new(Result)
	*p = x
	return p
}

func (x Result) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Result) Descriptor() protoreflect.EnumDescriptor {
	return file_example_proto_enumTypes[1].Descriptor()
}

func (Result) Type() protoreflect.EnumType {
	return &file_example_proto_enumTypes[1]
}

func (x Result) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Result.Descriptor instead.
func (Result) EnumDescriptor() ([]byte, []int) {
	return file_example_proto_rawDescGZIP(), []int{1}
}

type ClientLogin struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OpenId         string `protobuf:"bytes,1,opt,name=OpenId,proto3" json:"OpenId,omitempty"` //用户系统生成的OpenID
	Token          string `protobuf:"bytes,2,opt,name=Token,proto3" json:"Token,omitempty"`   //用户系统生成的token
	DeviceId       string `protobuf:"bytes,3,opt,name=DeviceId,proto3" json:"DeviceId,omitempty"`
	ChannelGroupId string `protobuf:"bytes,4,opt,name=ChannelGroupId,proto3" json:"ChannelGroupId,omitempty"` //渠道商ID
	ImeiIdfa       string `protobuf:"bytes,5,opt,name=ImeiIdfa,proto3" json:"ImeiIdfa,omitempty"`             //IMEI或IDFA,安卓：IMEI / IOS:IDFA
	Platform       int32  `protobuf:"varint,6,opt,name=Platform,proto3" json:"Platform,omitempty"`            //android=1、ios=2
}

func (x *ClientLogin) Reset() {
	*x = ClientLogin{}
	if protoimpl.UnsafeEnabled {
		mi := &file_example_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientLogin) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientLogin) ProtoMessage() {}

func (x *ClientLogin) ProtoReflect() protoreflect.Message {
	mi := &file_example_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientLogin.ProtoReflect.Descriptor instead.
func (*ClientLogin) Descriptor() ([]byte, []int) {
	return file_example_proto_rawDescGZIP(), []int{0}
}

func (x *ClientLogin) GetOpenId() string {
	if x != nil {
		return x.OpenId
	}
	return ""
}

func (x *ClientLogin) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *ClientLogin) GetDeviceId() string {
	if x != nil {
		return x.DeviceId
	}
	return ""
}

func (x *ClientLogin) GetChannelGroupId() string {
	if x != nil {
		return x.ChannelGroupId
	}
	return ""
}

func (x *ClientLogin) GetImeiIdfa() string {
	if x != nil {
		return x.ImeiIdfa
	}
	return ""
}

func (x *ClientLogin) GetPlatform() int32 {
	if x != nil {
		return x.Platform
	}
	return 0
}

// 登陆响应结构
type ServerLogin struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code      int32  `protobuf:"varint,1,opt,name=Code,proto3" json:"Code,omitempty"`          //返回状态，0：成功，其他失败,需要跳转到登录界面
	Uid       string `protobuf:"bytes,2,opt,name=Uid,proto3" json:"Uid,omitempty"`             //对应的后端uid
	GameToken string `protobuf:"bytes,3,opt,name=GameToken,proto3" json:"GameToken,omitempty"` //游戏生成的token，用于断线重连使用
}

func (x *ServerLogin) Reset() {
	*x = ServerLogin{}
	if protoimpl.UnsafeEnabled {
		mi := &file_example_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerLogin) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerLogin) ProtoMessage() {}

func (x *ServerLogin) ProtoReflect() protoreflect.Message {
	mi := &file_example_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerLogin.ProtoReflect.Descriptor instead.
func (*ServerLogin) Descriptor() ([]byte, []int) {
	return file_example_proto_rawDescGZIP(), []int{1}
}

func (x *ServerLogin) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *ServerLogin) GetUid() string {
	if x != nil {
		return x.Uid
	}
	return ""
}

func (x *ServerLogin) GetGameToken() string {
	if x != nil {
		return x.GameToken
	}
	return ""
}

var File_example_proto protoreflect.FileDescriptor

var file_example_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x22, 0xb7, 0x01, 0x0a, 0x0b, 0x43, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x4f, 0x70, 0x65,
	0x6e, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x4f, 0x70, 0x65, 0x6e, 0x49,
	0x64, 0x12, 0x14, 0x0a, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x44, 0x65, 0x76, 0x69, 0x63,
	0x65, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x44, 0x65, 0x76, 0x69, 0x63,
	0x65, 0x49, 0x64, 0x12, 0x26, 0x0a, 0x0e, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x47, 0x72,
	0x6f, 0x75, 0x70, 0x49, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x43, 0x68, 0x61,
	0x6e, 0x6e, 0x65, 0x6c, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x49,
	0x6d, 0x65, 0x69, 0x49, 0x64, 0x66, 0x61, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x49,
	0x6d, 0x65, 0x69, 0x49, 0x64, 0x66, 0x61, 0x12, 0x1a, 0x0a, 0x08, 0x50, 0x6c, 0x61, 0x74, 0x66,
	0x6f, 0x72, 0x6d, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x50, 0x6c, 0x61, 0x74, 0x66,
	0x6f, 0x72, 0x6d, 0x22, 0x51, 0x0a, 0x0b, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4c, 0x6f, 0x67,
	0x69, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x04, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x55, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x55, 0x69, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x47, 0x61, 0x6d, 0x65,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x47, 0x61, 0x6d,
	0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x2a, 0x53, 0x0a, 0x03, 0x43, 0x6d, 0x64, 0x12, 0x09, 0x0a,
	0x05, 0x43, 0x42, 0x65, 0x61, 0x74, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x53, 0x42, 0x65, 0x61,
	0x74, 0x10, 0x01, 0x12, 0x0e, 0x0a, 0x0a, 0x43, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x4b, 0x65,
	0x79, 0x10, 0x02, 0x12, 0x0e, 0x0a, 0x0a, 0x53, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x4b, 0x65,
	0x79, 0x10, 0x03, 0x12, 0x0a, 0x0a, 0x06, 0x43, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x10, 0x04, 0x12,
	0x0a, 0x0a, 0x06, 0x53, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x10, 0x05, 0x2a, 0x77, 0x0a, 0x06, 0x52,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x0f, 0x0a, 0x0b, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x5f,
	0x53, 0x55, 0x43, 0x43, 0x10, 0x00, 0x12, 0x0f, 0x0a, 0x0b, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x5f, 0x46, 0x41, 0x49, 0x4c, 0x10, 0x01, 0x12, 0x19, 0x0a, 0x15, 0x52, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x4e, 0x4f, 0x54, 0x5f, 0x45, 0x58, 0x49, 0x53, 0x54,
	0x10, 0x02, 0x12, 0x16, 0x0a, 0x12, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x5f, 0x50, 0x41, 0x52,
	0x41, 0x4d, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x03, 0x12, 0x18, 0x0a, 0x14, 0x52, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x5f, 0x44, 0x45, 0x43, 0x52, 0x59, 0x50, 0x54, 0x5f, 0x45, 0x52, 0x52,
	0x4f, 0x52, 0x10, 0x04, 0x42, 0x0a, 0x5a, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_example_proto_rawDescOnce sync.Once
	file_example_proto_rawDescData = file_example_proto_rawDesc
)

func file_example_proto_rawDescGZIP() []byte {
	file_example_proto_rawDescOnce.Do(func() {
		file_example_proto_rawDescData = protoimpl.X.CompressGZIP(file_example_proto_rawDescData)
	})
	return file_example_proto_rawDescData
}

var file_example_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_example_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_example_proto_goTypes = []interface{}{
	(Cmd)(0),            // 0: protocol.Cmd
	(Result)(0),         // 1: protocol.Result
	(*ClientLogin)(nil), // 2: protocol.ClientLogin
	(*ServerLogin)(nil), // 3: protocol.ServerLogin
}
var file_example_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_example_proto_init() }
func file_example_proto_init() {
	if File_example_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_example_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientLogin); i {
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
		file_example_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerLogin); i {
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
			RawDescriptor: file_example_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_example_proto_goTypes,
		DependencyIndexes: file_example_proto_depIdxs,
		EnumInfos:         file_example_proto_enumTypes,
		MessageInfos:      file_example_proto_msgTypes,
	}.Build()
	File_example_proto = out.File
	file_example_proto_rawDesc = nil
	file_example_proto_goTypes = nil
	file_example_proto_depIdxs = nil
}