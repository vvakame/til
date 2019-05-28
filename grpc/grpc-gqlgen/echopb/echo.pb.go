// Code generated by protoc-gen-go. DO NOT EDIT.
// source: echo.proto

package echopb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	_ "github.com/vvakame/til/grpc/grpc-gqlgen/proto-extentions"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type SayRequest struct {
	MessageId            string   `protobuf:"bytes,1,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
	MessageBody          string   `protobuf:"bytes,2,opt,name=message_body,json=messageBody,proto3" json:"message_body,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SayRequest) Reset()         { *m = SayRequest{} }
func (m *SayRequest) String() string { return proto.CompactTextString(m) }
func (*SayRequest) ProtoMessage()    {}
func (*SayRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_08134aea513e0001, []int{0}
}

func (m *SayRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SayRequest.Unmarshal(m, b)
}
func (m *SayRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SayRequest.Marshal(b, m, deterministic)
}
func (m *SayRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SayRequest.Merge(m, src)
}
func (m *SayRequest) XXX_Size() int {
	return xxx_messageInfo_SayRequest.Size(m)
}
func (m *SayRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SayRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SayRequest proto.InternalMessageInfo

func (m *SayRequest) GetMessageId() string {
	if m != nil {
		return m.MessageId
	}
	return ""
}

func (m *SayRequest) GetMessageBody() string {
	if m != nil {
		return m.MessageBody
	}
	return ""
}

type SayResponse struct {
	MessageId            string               `protobuf:"bytes,1,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
	MessageBody          string               `protobuf:"bytes,2,opt,name=message_body,json=messageBody,proto3" json:"message_body,omitempty"`
	Received             *timestamp.Timestamp `protobuf:"bytes,3,opt,name=received,proto3" json:"received,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *SayResponse) Reset()         { *m = SayResponse{} }
func (m *SayResponse) String() string { return proto.CompactTextString(m) }
func (*SayResponse) ProtoMessage()    {}
func (*SayResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_08134aea513e0001, []int{1}
}

func (m *SayResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SayResponse.Unmarshal(m, b)
}
func (m *SayResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SayResponse.Marshal(b, m, deterministic)
}
func (m *SayResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SayResponse.Merge(m, src)
}
func (m *SayResponse) XXX_Size() int {
	return xxx_messageInfo_SayResponse.Size(m)
}
func (m *SayResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SayResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SayResponse proto.InternalMessageInfo

func (m *SayResponse) GetMessageId() string {
	if m != nil {
		return m.MessageId
	}
	return ""
}

func (m *SayResponse) GetMessageBody() string {
	if m != nil {
		return m.MessageBody
	}
	return ""
}

func (m *SayResponse) GetReceived() *timestamp.Timestamp {
	if m != nil {
		return m.Received
	}
	return nil
}

type Example1 struct {
	Foo                  *Example1_InMessage `protobuf:"bytes,1,opt,name=foo,proto3" json:"foo,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *Example1) Reset()         { *m = Example1{} }
func (m *Example1) String() string { return proto.CompactTextString(m) }
func (*Example1) ProtoMessage()    {}
func (*Example1) Descriptor() ([]byte, []int) {
	return fileDescriptor_08134aea513e0001, []int{2}
}

func (m *Example1) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Example1.Unmarshal(m, b)
}
func (m *Example1) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Example1.Marshal(b, m, deterministic)
}
func (m *Example1) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Example1.Merge(m, src)
}
func (m *Example1) XXX_Size() int {
	return xxx_messageInfo_Example1.Size(m)
}
func (m *Example1) XXX_DiscardUnknown() {
	xxx_messageInfo_Example1.DiscardUnknown(m)
}

var xxx_messageInfo_Example1 proto.InternalMessageInfo

func (m *Example1) GetFoo() *Example1_InMessage {
	if m != nil {
		return m.Foo
	}
	return nil
}

type Example1_InMessage struct {
	Bar                  string   `protobuf:"bytes,1,opt,name=bar,proto3" json:"bar,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Example1_InMessage) Reset()         { *m = Example1_InMessage{} }
func (m *Example1_InMessage) String() string { return proto.CompactTextString(m) }
func (*Example1_InMessage) ProtoMessage()    {}
func (*Example1_InMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_08134aea513e0001, []int{2, 0}
}

func (m *Example1_InMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Example1_InMessage.Unmarshal(m, b)
}
func (m *Example1_InMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Example1_InMessage.Marshal(b, m, deterministic)
}
func (m *Example1_InMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Example1_InMessage.Merge(m, src)
}
func (m *Example1_InMessage) XXX_Size() int {
	return xxx_messageInfo_Example1_InMessage.Size(m)
}
func (m *Example1_InMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_Example1_InMessage.DiscardUnknown(m)
}

var xxx_messageInfo_Example1_InMessage proto.InternalMessageInfo

func (m *Example1_InMessage) GetBar() string {
	if m != nil {
		return m.Bar
	}
	return ""
}

type Example2 struct {
	Hoge                 *Example2_InMessage `protobuf:"bytes,1,opt,name=hoge,proto3" json:"hoge,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *Example2) Reset()         { *m = Example2{} }
func (m *Example2) String() string { return proto.CompactTextString(m) }
func (*Example2) ProtoMessage()    {}
func (*Example2) Descriptor() ([]byte, []int) {
	return fileDescriptor_08134aea513e0001, []int{3}
}

func (m *Example2) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Example2.Unmarshal(m, b)
}
func (m *Example2) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Example2.Marshal(b, m, deterministic)
}
func (m *Example2) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Example2.Merge(m, src)
}
func (m *Example2) XXX_Size() int {
	return xxx_messageInfo_Example2.Size(m)
}
func (m *Example2) XXX_DiscardUnknown() {
	xxx_messageInfo_Example2.DiscardUnknown(m)
}

var xxx_messageInfo_Example2 proto.InternalMessageInfo

func (m *Example2) GetHoge() *Example2_InMessage {
	if m != nil {
		return m.Hoge
	}
	return nil
}

type Example2_InMessage struct {
	Fuga                 string   `protobuf:"bytes,1,opt,name=fuga,proto3" json:"fuga,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Example2_InMessage) Reset()         { *m = Example2_InMessage{} }
func (m *Example2_InMessage) String() string { return proto.CompactTextString(m) }
func (*Example2_InMessage) ProtoMessage()    {}
func (*Example2_InMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_08134aea513e0001, []int{3, 0}
}

func (m *Example2_InMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Example2_InMessage.Unmarshal(m, b)
}
func (m *Example2_InMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Example2_InMessage.Marshal(b, m, deterministic)
}
func (m *Example2_InMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Example2_InMessage.Merge(m, src)
}
func (m *Example2_InMessage) XXX_Size() int {
	return xxx_messageInfo_Example2_InMessage.Size(m)
}
func (m *Example2_InMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_Example2_InMessage.DiscardUnknown(m)
}

var xxx_messageInfo_Example2_InMessage proto.InternalMessageInfo

func (m *Example2_InMessage) GetFuga() string {
	if m != nil {
		return m.Fuga
	}
	return ""
}

func init() {
	proto.RegisterType((*SayRequest)(nil), "echo.SayRequest")
	proto.RegisterType((*SayResponse)(nil), "echo.SayResponse")
	proto.RegisterType((*Example1)(nil), "echo.Example1")
	proto.RegisterType((*Example1_InMessage)(nil), "echo.Example1.InMessage")
	proto.RegisterType((*Example2)(nil), "echo.Example2")
	proto.RegisterType((*Example2_InMessage)(nil), "echo.Example2.InMessage")
}

func init() { proto.RegisterFile("echo.proto", fileDescriptor_08134aea513e0001) }

var fileDescriptor_08134aea513e0001 = []byte{
	// 484 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x52, 0xb1, 0x8e, 0xd3, 0x4c,
	0x10, 0x96, 0x93, 0xe8, 0xfe, 0x64, 0x92, 0x1f, 0xf9, 0x56, 0x14, 0x2b, 0x8b, 0xe8, 0x0e, 0x93,
	0x22, 0x44, 0x9c, 0x57, 0x31, 0x12, 0x48, 0x47, 0x17, 0xe9, 0x8a, 0x14, 0x27, 0x90, 0x03, 0x42,
	0x50, 0x80, 0xd6, 0xf6, 0xc6, 0xb1, 0xb0, 0x3d, 0x3e, 0x7b, 0x1d, 0x9d, 0x29, 0xaf, 0xa4, 0x84,
	0x8e, 0xd7, 0x40, 0x79, 0x12, 0x5e, 0x81, 0x07, 0x41, 0xde, 0x75, 0x42, 0x0a, 0x5a, 0x1a, 0x7b,
	0x76, 0xe6, 0x9b, 0x6f, 0x66, 0xbf, 0x6f, 0x01, 0x44, 0xb0, 0x41, 0x27, 0x2f, 0x50, 0x22, 0xe9,
	0x35, 0xb1, 0xf5, 0x20, 0x42, 0x8c, 0x12, 0xc1, 0x78, 0x1e, 0x33, 0x9e, 0x65, 0x28, 0xb9, 0x8c,
	0x31, 0x2b, 0x35, 0xc6, 0x3a, 0x6b, 0xab, 0xea, 0xe4, 0x57, 0x6b, 0x26, 0xe3, 0x54, 0x94, 0x92,
	0xa7, 0x79, 0x0b, 0x18, 0xab, 0xdf, 0x85, 0xb8, 0x95, 0x22, 0x53, 0x8d, 0x2c, 0xba, 0x49, 0x22,
	0x91, 0xe9, 0xb2, 0xfd, 0x19, 0x60, 0xc5, 0x6b, 0x4f, 0xdc, 0x54, 0xa2, 0x94, 0xe4, 0x39, 0x40,
	0x2a, 0xca, 0x92, 0x47, 0xe2, 0x63, 0x1c, 0x52, 0xe3, 0xdc, 0x98, 0x0e, 0x16, 0xf4, 0xfb, 0x8e,
	0xde, 0x07, 0x33, 0x48, 0x62, 0x91, 0xc9, 0xeb, 0x4a, 0x8f, 0x5f, 0x86, 0xa6, 0xe1, 0x0d, 0x5a,
	0xec, 0x32, 0x24, 0x0f, 0x61, 0xb4, 0x6f, 0xf4, 0x31, 0xac, 0x69, 0xa7, 0x69, 0xf5, 0x86, 0x6d,
	0x6e, 0x81, 0x61, 0x7d, 0x69, 0x7e, 0xdd, 0xd1, 0x11, 0xf4, 0x57, 0xbc, 0x5e, 0x66, 0x79, 0x25,
	0xcd, 0x8e, 0xfd, 0xc3, 0x80, 0xa1, 0x1a, 0x5e, 0xe6, 0x98, 0x95, 0xe2, 0x5f, 0x4e, 0x27, 0xcf,
	0xa0, 0x5f, 0x88, 0x40, 0xc4, 0x5b, 0x11, 0xd2, 0xee, 0xb9, 0x31, 0x1d, 0xba, 0x96, 0xa3, 0xa5,
	0x73, 0xf6, 0xd2, 0x39, 0xaf, 0xf7, 0xd2, 0x79, 0x07, 0x6c, 0xbb, 0x75, 0xa3, 0xd1, 0x2b, 0x5e,
	0x27, 0xc8, 0x43, 0xfb, 0x0d, 0xf4, 0xaf, 0x6e, 0x79, 0x9a, 0x27, 0x62, 0x4e, 0x66, 0xd0, 0x5d,
	0x23, 0xaa, 0x55, 0x87, 0x2e, 0x75, 0x94, 0x77, 0xfb, 0xa2, 0xb3, 0xcc, 0xae, 0xf5, 0x02, 0x5e,
	0x03, 0xb2, 0xc6, 0x30, 0x38, 0x64, 0x88, 0x09, 0x5d, 0x9f, 0x17, 0xfa, 0x8e, 0x5e, 0x13, 0xda,
	0xef, 0x0e, 0xb4, 0x2e, 0x79, 0x02, 0xbd, 0x0d, 0x46, 0xe2, 0xaf, 0xbc, 0xee, 0x11, 0xaf, 0x42,
	0x59, 0x67, 0xc7, 0xc4, 0x04, 0x7a, 0xeb, 0x2a, 0xe2, 0x2d, 0xb3, 0x8a, 0xdd, 0xb7, 0xd0, 0xbb,
	0x0a, 0x36, 0x48, 0x5e, 0x42, 0x77, 0xc5, 0x6b, 0x62, 0x6a, 0xbe, 0x3f, 0xb6, 0x5b, 0xa7, 0x47,
	0x19, 0xed, 0x85, 0xfd, 0xe8, 0xcb, 0x8e, 0xf6, 0xc9, 0xc9, 0x1a, 0xd1, 0xe7, 0xc5, 0xdd, 0xcf,
	0x5f, 0xdf, 0x3a, 0xa7, 0xf6, 0x88, 0x6d, 0xe7, 0xac, 0x41, 0xb2, 0x92, 0xd7, 0x97, 0xc6, 0x6c,
	0x81, 0x77, 0x3b, 0xfa, 0x02, 0x2c, 0xf8, 0xff, 0xc3, 0xd4, 0x99, 0x3d, 0x6e, 0xd9, 0x26, 0xe4,
	0xbf, 0xc9, 0x5c, 0x59, 0x4c, 0x3b, 0x30, 0x86, 0x7b, 0x6d, 0x4d, 0xf3, 0x4e, 0xc8, 0x60, 0x32,
	0x6f, 0x95, 0xa4, 0xc6, 0x7b, 0x27, 0x8a, 0xe5, 0xa6, 0xf2, 0x9d, 0x00, 0x53, 0xb6, 0xdd, 0xf2,
	0x4f, 0x3c, 0x15, 0x4c, 0xc6, 0x09, 0x8b, 0x8a, 0x3c, 0x50, 0x9f, 0x0b, 0xfd, 0x58, 0xd5, 0xd4,
	0xdc, 0xf7, 0x4f, 0x94, 0x57, 0x4f, 0x7f, 0x07, 0x00, 0x00, 0xff, 0xff, 0x17, 0x5a, 0x4a, 0x61,
	0x26, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// EchoClient is the client API for Echo service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type EchoClient interface {
	Say(ctx context.Context, in *SayRequest, opts ...grpc.CallOption) (*SayResponse, error)
}

type echoClient struct {
	cc *grpc.ClientConn
}

func NewEchoClient(cc *grpc.ClientConn) EchoClient {
	return &echoClient{cc}
}

func (c *echoClient) Say(ctx context.Context, in *SayRequest, opts ...grpc.CallOption) (*SayResponse, error) {
	out := new(SayResponse)
	err := c.cc.Invoke(ctx, "/echo.Echo/Say", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EchoServer is the server API for Echo service.
type EchoServer interface {
	Say(context.Context, *SayRequest) (*SayResponse, error)
}

func RegisterEchoServer(s *grpc.Server, srv EchoServer) {
	s.RegisterService(&_Echo_serviceDesc, srv)
}

func _Echo_Say_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SayRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EchoServer).Say(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/echo.Echo/Say",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EchoServer).Say(ctx, req.(*SayRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Echo_serviceDesc = grpc.ServiceDesc{
	ServiceName: "echo.Echo",
	HandlerType: (*EchoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Say",
			Handler:    _Echo_Say_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "echo.proto",
}