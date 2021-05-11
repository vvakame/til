// Code generated by protoc-gen-go. DO NOT EDIT.
// source: gqlgen-proto/options.proto

package gqlgen_proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	descriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
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

type MethodType int32

const (
	MethodType_OPERATION_UNKNOWN      MethodType = 0
	MethodType_OPERATION_QUERY        MethodType = 1
	MethodType_OPERATION_MUTATION     MethodType = 2
	MethodType_OPERATION_SUBSCRIPTION MethodType = 3
)

var MethodType_name = map[int32]string{
	0: "OPERATION_UNKNOWN",
	1: "OPERATION_QUERY",
	2: "OPERATION_MUTATION",
	3: "OPERATION_SUBSCRIPTION",
}

var MethodType_value = map[string]int32{
	"OPERATION_UNKNOWN":      0,
	"OPERATION_QUERY":        1,
	"OPERATION_MUTATION":     2,
	"OPERATION_SUBSCRIPTION": 3,
}

func (x MethodType) String() string {
	return proto.EnumName(MethodType_name, int32(x))
}

func (MethodType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_03469d44191f4235, []int{0}
}

type MessageType int32

const (
	MessageType_TYPE_UNKNOWN MessageType = 0
	MessageType_TYPE_TYPE    MessageType = 1
	MessageType_TYPE_INPUT   MessageType = 2
)

var MessageType_name = map[int32]string{
	0: "TYPE_UNKNOWN",
	1: "TYPE_TYPE",
	2: "TYPE_INPUT",
}

var MessageType_value = map[string]int32{
	"TYPE_UNKNOWN": 0,
	"TYPE_TYPE":    1,
	"TYPE_INPUT":   2,
}

func (x MessageType) String() string {
	return proto.EnumName(MessageType_name, int32(x))
}

func (MessageType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_03469d44191f4235, []int{1}
}

type FileRule struct {
	MethodRule           []*MethodInferenceRule  `protobuf:"bytes,1,rep,name=method_rule,json=methodRule,proto3" json:"method_rule,omitempty"`
	MessageRule          []*MessageInferenceRule `protobuf:"bytes,2,rep,name=message_rule,json=messageRule,proto3" json:"message_rule,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *FileRule) Reset()         { *m = FileRule{} }
func (m *FileRule) String() string { return proto.CompactTextString(m) }
func (*FileRule) ProtoMessage()    {}
func (*FileRule) Descriptor() ([]byte, []int) {
	return fileDescriptor_03469d44191f4235, []int{0}
}

func (m *FileRule) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FileRule.Unmarshal(m, b)
}
func (m *FileRule) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FileRule.Marshal(b, m, deterministic)
}
func (m *FileRule) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FileRule.Merge(m, src)
}
func (m *FileRule) XXX_Size() int {
	return xxx_messageInfo_FileRule.Size(m)
}
func (m *FileRule) XXX_DiscardUnknown() {
	xxx_messageInfo_FileRule.DiscardUnknown(m)
}

var xxx_messageInfo_FileRule proto.InternalMessageInfo

func (m *FileRule) GetMethodRule() []*MethodInferenceRule {
	if m != nil {
		return m.MethodRule
	}
	return nil
}

func (m *FileRule) GetMessageRule() []*MessageInferenceRule {
	if m != nil {
		return m.MessageRule
	}
	return nil
}

type MethodInferenceRule struct {
	Src                  string     `protobuf:"bytes,1,opt,name=src,proto3" json:"src,omitempty"`
	Dest                 string     `protobuf:"bytes,2,opt,name=dest,proto3" json:"dest,omitempty"`
	Type                 MethodType `protobuf:"varint,3,opt,name=type,proto3,enum=gqlgen.api.MethodType" json:"type,omitempty"`
	Skip                 bool       `protobuf:"varint,4,opt,name=skip,proto3" json:"skip,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *MethodInferenceRule) Reset()         { *m = MethodInferenceRule{} }
func (m *MethodInferenceRule) String() string { return proto.CompactTextString(m) }
func (*MethodInferenceRule) ProtoMessage()    {}
func (*MethodInferenceRule) Descriptor() ([]byte, []int) {
	return fileDescriptor_03469d44191f4235, []int{1}
}

func (m *MethodInferenceRule) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MethodInferenceRule.Unmarshal(m, b)
}
func (m *MethodInferenceRule) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MethodInferenceRule.Marshal(b, m, deterministic)
}
func (m *MethodInferenceRule) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MethodInferenceRule.Merge(m, src)
}
func (m *MethodInferenceRule) XXX_Size() int {
	return xxx_messageInfo_MethodInferenceRule.Size(m)
}
func (m *MethodInferenceRule) XXX_DiscardUnknown() {
	xxx_messageInfo_MethodInferenceRule.DiscardUnknown(m)
}

var xxx_messageInfo_MethodInferenceRule proto.InternalMessageInfo

func (m *MethodInferenceRule) GetSrc() string {
	if m != nil {
		return m.Src
	}
	return ""
}

func (m *MethodInferenceRule) GetDest() string {
	if m != nil {
		return m.Dest
	}
	return ""
}

func (m *MethodInferenceRule) GetType() MethodType {
	if m != nil {
		return m.Type
	}
	return MethodType_OPERATION_UNKNOWN
}

func (m *MethodInferenceRule) GetSkip() bool {
	if m != nil {
		return m.Skip
	}
	return false
}

type MessageInferenceRule struct {
	Src                  string      `protobuf:"bytes,1,opt,name=src,proto3" json:"src,omitempty"`
	Dest                 string      `protobuf:"bytes,2,opt,name=dest,proto3" json:"dest,omitempty"`
	Type                 MessageType `protobuf:"varint,3,opt,name=type,proto3,enum=gqlgen.api.MessageType" json:"type,omitempty"`
	Skip                 bool        `protobuf:"varint,4,opt,name=skip,proto3" json:"skip,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *MessageInferenceRule) Reset()         { *m = MessageInferenceRule{} }
func (m *MessageInferenceRule) String() string { return proto.CompactTextString(m) }
func (*MessageInferenceRule) ProtoMessage()    {}
func (*MessageInferenceRule) Descriptor() ([]byte, []int) {
	return fileDescriptor_03469d44191f4235, []int{2}
}

func (m *MessageInferenceRule) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MessageInferenceRule.Unmarshal(m, b)
}
func (m *MessageInferenceRule) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MessageInferenceRule.Marshal(b, m, deterministic)
}
func (m *MessageInferenceRule) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MessageInferenceRule.Merge(m, src)
}
func (m *MessageInferenceRule) XXX_Size() int {
	return xxx_messageInfo_MessageInferenceRule.Size(m)
}
func (m *MessageInferenceRule) XXX_DiscardUnknown() {
	xxx_messageInfo_MessageInferenceRule.DiscardUnknown(m)
}

var xxx_messageInfo_MessageInferenceRule proto.InternalMessageInfo

func (m *MessageInferenceRule) GetSrc() string {
	if m != nil {
		return m.Src
	}
	return ""
}

func (m *MessageInferenceRule) GetDest() string {
	if m != nil {
		return m.Dest
	}
	return ""
}

func (m *MessageInferenceRule) GetType() MessageType {
	if m != nil {
		return m.Type
	}
	return MessageType_TYPE_UNKNOWN
}

func (m *MessageInferenceRule) GetSkip() bool {
	if m != nil {
		return m.Skip
	}
	return false
}

type SchemaRule struct {
	// Types that are valid to be assigned to Pattern:
	//	*SchemaRule_Query
	//	*SchemaRule_Mutation
	//	*SchemaRule_Subscription
	Pattern              isSchemaRule_Pattern `protobuf_oneof:"pattern"`
	Skip                 bool                 `protobuf:"varint,4,opt,name=skip,proto3" json:"skip,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *SchemaRule) Reset()         { *m = SchemaRule{} }
func (m *SchemaRule) String() string { return proto.CompactTextString(m) }
func (*SchemaRule) ProtoMessage()    {}
func (*SchemaRule) Descriptor() ([]byte, []int) {
	return fileDescriptor_03469d44191f4235, []int{3}
}

func (m *SchemaRule) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SchemaRule.Unmarshal(m, b)
}
func (m *SchemaRule) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SchemaRule.Marshal(b, m, deterministic)
}
func (m *SchemaRule) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SchemaRule.Merge(m, src)
}
func (m *SchemaRule) XXX_Size() int {
	return xxx_messageInfo_SchemaRule.Size(m)
}
func (m *SchemaRule) XXX_DiscardUnknown() {
	xxx_messageInfo_SchemaRule.DiscardUnknown(m)
}

var xxx_messageInfo_SchemaRule proto.InternalMessageInfo

type isSchemaRule_Pattern interface {
	isSchemaRule_Pattern()
}

type SchemaRule_Query struct {
	Query string `protobuf:"bytes,1,opt,name=query,proto3,oneof"`
}

type SchemaRule_Mutation struct {
	Mutation string `protobuf:"bytes,2,opt,name=mutation,proto3,oneof"`
}

type SchemaRule_Subscription struct {
	Subscription string `protobuf:"bytes,3,opt,name=subscription,proto3,oneof"`
}

func (*SchemaRule_Query) isSchemaRule_Pattern() {}

func (*SchemaRule_Mutation) isSchemaRule_Pattern() {}

func (*SchemaRule_Subscription) isSchemaRule_Pattern() {}

func (m *SchemaRule) GetPattern() isSchemaRule_Pattern {
	if m != nil {
		return m.Pattern
	}
	return nil
}

func (m *SchemaRule) GetQuery() string {
	if x, ok := m.GetPattern().(*SchemaRule_Query); ok {
		return x.Query
	}
	return ""
}

func (m *SchemaRule) GetMutation() string {
	if x, ok := m.GetPattern().(*SchemaRule_Mutation); ok {
		return x.Mutation
	}
	return ""
}

func (m *SchemaRule) GetSubscription() string {
	if x, ok := m.GetPattern().(*SchemaRule_Subscription); ok {
		return x.Subscription
	}
	return ""
}

func (m *SchemaRule) GetSkip() bool {
	if m != nil {
		return m.Skip
	}
	return false
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*SchemaRule) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*SchemaRule_Query)(nil),
		(*SchemaRule_Mutation)(nil),
		(*SchemaRule_Subscription)(nil),
	}
}

type MessageRule struct {
	Alias                string      `protobuf:"bytes,1,opt,name=alias,proto3" json:"alias,omitempty"`
	Type                 MessageType `protobuf:"varint,2,opt,name=type,proto3,enum=gqlgen.api.MessageType" json:"type,omitempty"`
	Skip                 bool        `protobuf:"varint,4,opt,name=skip,proto3" json:"skip,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *MessageRule) Reset()         { *m = MessageRule{} }
func (m *MessageRule) String() string { return proto.CompactTextString(m) }
func (*MessageRule) ProtoMessage()    {}
func (*MessageRule) Descriptor() ([]byte, []int) {
	return fileDescriptor_03469d44191f4235, []int{4}
}

func (m *MessageRule) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MessageRule.Unmarshal(m, b)
}
func (m *MessageRule) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MessageRule.Marshal(b, m, deterministic)
}
func (m *MessageRule) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MessageRule.Merge(m, src)
}
func (m *MessageRule) XXX_Size() int {
	return xxx_messageInfo_MessageRule.Size(m)
}
func (m *MessageRule) XXX_DiscardUnknown() {
	xxx_messageInfo_MessageRule.DiscardUnknown(m)
}

var xxx_messageInfo_MessageRule proto.InternalMessageInfo

func (m *MessageRule) GetAlias() string {
	if m != nil {
		return m.Alias
	}
	return ""
}

func (m *MessageRule) GetType() MessageType {
	if m != nil {
		return m.Type
	}
	return MessageType_TYPE_UNKNOWN
}

func (m *MessageRule) GetSkip() bool {
	if m != nil {
		return m.Skip
	}
	return false
}

type FieldRule struct {
	Alias                string   `protobuf:"bytes,1,opt,name=alias,proto3" json:"alias,omitempty"`
	Optional             bool     `protobuf:"varint,2,opt,name=optional,proto3" json:"optional,omitempty"`
	Id                   bool     `protobuf:"varint,3,opt,name=id,proto3" json:"id,omitempty"`
	Skip                 bool     `protobuf:"varint,4,opt,name=skip,proto3" json:"skip,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FieldRule) Reset()         { *m = FieldRule{} }
func (m *FieldRule) String() string { return proto.CompactTextString(m) }
func (*FieldRule) ProtoMessage()    {}
func (*FieldRule) Descriptor() ([]byte, []int) {
	return fileDescriptor_03469d44191f4235, []int{5}
}

func (m *FieldRule) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FieldRule.Unmarshal(m, b)
}
func (m *FieldRule) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FieldRule.Marshal(b, m, deterministic)
}
func (m *FieldRule) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FieldRule.Merge(m, src)
}
func (m *FieldRule) XXX_Size() int {
	return xxx_messageInfo_FieldRule.Size(m)
}
func (m *FieldRule) XXX_DiscardUnknown() {
	xxx_messageInfo_FieldRule.DiscardUnknown(m)
}

var xxx_messageInfo_FieldRule proto.InternalMessageInfo

func (m *FieldRule) GetAlias() string {
	if m != nil {
		return m.Alias
	}
	return ""
}

func (m *FieldRule) GetOptional() bool {
	if m != nil {
		return m.Optional
	}
	return false
}

func (m *FieldRule) GetId() bool {
	if m != nil {
		return m.Id
	}
	return false
}

func (m *FieldRule) GetSkip() bool {
	if m != nil {
		return m.Skip
	}
	return false
}

type EnumRule struct {
	Alias                string   `protobuf:"bytes,1,opt,name=alias,proto3" json:"alias,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EnumRule) Reset()         { *m = EnumRule{} }
func (m *EnumRule) String() string { return proto.CompactTextString(m) }
func (*EnumRule) ProtoMessage()    {}
func (*EnumRule) Descriptor() ([]byte, []int) {
	return fileDescriptor_03469d44191f4235, []int{6}
}

func (m *EnumRule) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EnumRule.Unmarshal(m, b)
}
func (m *EnumRule) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EnumRule.Marshal(b, m, deterministic)
}
func (m *EnumRule) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EnumRule.Merge(m, src)
}
func (m *EnumRule) XXX_Size() int {
	return xxx_messageInfo_EnumRule.Size(m)
}
func (m *EnumRule) XXX_DiscardUnknown() {
	xxx_messageInfo_EnumRule.DiscardUnknown(m)
}

var xxx_messageInfo_EnumRule proto.InternalMessageInfo

func (m *EnumRule) GetAlias() string {
	if m != nil {
		return m.Alias
	}
	return ""
}

type EnumValueRule struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EnumValueRule) Reset()         { *m = EnumValueRule{} }
func (m *EnumValueRule) String() string { return proto.CompactTextString(m) }
func (*EnumValueRule) ProtoMessage()    {}
func (*EnumValueRule) Descriptor() ([]byte, []int) {
	return fileDescriptor_03469d44191f4235, []int{7}
}

func (m *EnumValueRule) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EnumValueRule.Unmarshal(m, b)
}
func (m *EnumValueRule) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EnumValueRule.Marshal(b, m, deterministic)
}
func (m *EnumValueRule) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EnumValueRule.Merge(m, src)
}
func (m *EnumValueRule) XXX_Size() int {
	return xxx_messageInfo_EnumValueRule.Size(m)
}
func (m *EnumValueRule) XXX_DiscardUnknown() {
	xxx_messageInfo_EnumValueRule.DiscardUnknown(m)
}

var xxx_messageInfo_EnumValueRule proto.InternalMessageInfo

var E_Resolver = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FileOptions)(nil),
	ExtensionType: (*FileRule)(nil),
	Field:         50000,
	Name:          "gqlgen.api.resolver",
	Tag:           "bytes,50000,opt,name=resolver",
	Filename:      "gqlgen-proto/options.proto",
}

var E_Schema = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MethodOptions)(nil),
	ExtensionType: (*SchemaRule)(nil),
	Field:         50001,
	Name:          "gqlgen.api.schema",
	Tag:           "bytes,50001,opt,name=schema",
	Filename:      "gqlgen-proto/options.proto",
}

var E_Type = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MessageOptions)(nil),
	ExtensionType: (*MessageRule)(nil),
	Field:         50002,
	Name:          "gqlgen.api.type",
	Tag:           "bytes,50002,opt,name=type",
	Filename:      "gqlgen-proto/options.proto",
}

var E_Field = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*FieldRule)(nil),
	Field:         50003,
	Name:          "gqlgen.api.field",
	Tag:           "bytes,50003,opt,name=field",
	Filename:      "gqlgen-proto/options.proto",
}

var E_Enum = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.EnumOptions)(nil),
	ExtensionType: (*EnumRule)(nil),
	Field:         50004,
	Name:          "gqlgen.api.enum",
	Tag:           "bytes,50004,opt,name=enum",
	Filename:      "gqlgen-proto/options.proto",
}

var E_EnumValue = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.EnumValueOptions)(nil),
	ExtensionType: (*EnumValueRule)(nil),
	Field:         50005,
	Name:          "gqlgen.api.enum_value",
	Tag:           "bytes,50005,opt,name=enum_value",
	Filename:      "gqlgen-proto/options.proto",
}

func init() {
	proto.RegisterEnum("gqlgen.api.MethodType", MethodType_name, MethodType_value)
	proto.RegisterEnum("gqlgen.api.MessageType", MessageType_name, MessageType_value)
	proto.RegisterType((*FileRule)(nil), "gqlgen.api.FileRule")
	proto.RegisterType((*MethodInferenceRule)(nil), "gqlgen.api.MethodInferenceRule")
	proto.RegisterType((*MessageInferenceRule)(nil), "gqlgen.api.MessageInferenceRule")
	proto.RegisterType((*SchemaRule)(nil), "gqlgen.api.SchemaRule")
	proto.RegisterType((*MessageRule)(nil), "gqlgen.api.MessageRule")
	proto.RegisterType((*FieldRule)(nil), "gqlgen.api.FieldRule")
	proto.RegisterType((*EnumRule)(nil), "gqlgen.api.EnumRule")
	proto.RegisterType((*EnumValueRule)(nil), "gqlgen.api.EnumValueRule")
	proto.RegisterExtension(E_Resolver)
	proto.RegisterExtension(E_Schema)
	proto.RegisterExtension(E_Type)
	proto.RegisterExtension(E_Field)
	proto.RegisterExtension(E_Enum)
	proto.RegisterExtension(E_EnumValue)
}

func init() { proto.RegisterFile("gqlgen-proto/options.proto", fileDescriptor_03469d44191f4235) }

var fileDescriptor_03469d44191f4235 = []byte{
	// 706 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x95, 0x5b, 0x4f, 0x13, 0x41,
	0x14, 0xc7, 0xd9, 0xb6, 0xe0, 0xf6, 0x94, 0x4b, 0x1d, 0xa0, 0xd6, 0x06, 0x65, 0x6d, 0x7c, 0x68,
	0x30, 0xec, 0x26, 0xf8, 0x56, 0x13, 0x23, 0x90, 0x12, 0x1a, 0xa5, 0xad, 0x43, 0xab, 0xc1, 0xc4,
	0x34, 0x4b, 0x77, 0x58, 0x36, 0xec, 0x8d, 0xbd, 0x34, 0x21, 0xc6, 0x57, 0x7d, 0xf6, 0x9b, 0x79,
	0xfd, 0x3e, 0x66, 0xcf, 0xec, 0x85, 0xca, 0x7a, 0x89, 0x2f, 0xdd, 0x39, 0x67, 0x66, 0x7e, 0xff,
	0x39, 0x67, 0xfe, 0x03, 0xd0, 0xd0, 0x2f, 0x4d, 0x9d, 0xd9, 0xdb, 0xae, 0xe7, 0x04, 0x8e, 0xe2,
	0xb8, 0x81, 0xe1, 0xd8, 0xbe, 0x8c, 0x11, 0x01, 0x3e, 0x27, 0xab, 0xae, 0xd1, 0x90, 0x74, 0xc7,
	0xd1, 0x4d, 0xa6, 0xe0, 0xcc, 0x69, 0x78, 0xa6, 0x68, 0xcc, 0x9f, 0x78, 0x86, 0x1b, 0x38, 0x1e,
	0x5f, 0xdd, 0xfc, 0x24, 0x80, 0x78, 0x60, 0x98, 0x8c, 0x86, 0x26, 0x23, 0xcf, 0xa0, 0x62, 0xb1,
	0xe0, 0xdc, 0xd1, 0xc6, 0x5e, 0x68, 0xb2, 0xba, 0x20, 0x15, 0x5b, 0x95, 0x9d, 0x4d, 0x39, 0x03,
	0xca, 0x47, 0x38, 0xdd, 0xb5, 0xcf, 0x98, 0xc7, 0xec, 0x09, 0xee, 0xa2, 0xc0, 0xf7, 0x20, 0x61,
	0x1f, 0x16, 0x2d, 0xe6, 0xfb, 0xaa, 0xce, 0x38, 0xa2, 0x80, 0x08, 0x69, 0x16, 0x81, 0xf3, 0xb3,
	0x8c, 0x4a, 0xbc, 0x2b, 0x0a, 0x9a, 0xef, 0x60, 0x35, 0x47, 0x87, 0x54, 0xa1, 0xe8, 0x7b, 0x93,
	0xba, 0x20, 0x09, 0xad, 0x32, 0x8d, 0x86, 0x84, 0x40, 0x49, 0x63, 0x7e, 0x50, 0x2f, 0x60, 0x0a,
	0xc7, 0x64, 0x0b, 0x4a, 0xc1, 0x95, 0xcb, 0xea, 0x45, 0x49, 0x68, 0x2d, 0xef, 0xd4, 0x6e, 0x1e,
	0x7e, 0x78, 0xe5, 0x32, 0x8a, 0x6b, 0xa2, 0xfd, 0xfe, 0x85, 0xe1, 0xd6, 0x4b, 0x92, 0xd0, 0x12,
	0x29, 0x8e, 0x9b, 0xef, 0x61, 0x2d, 0xef, 0x84, 0xff, 0xa8, 0xfe, 0x68, 0x46, 0xfd, 0x4e, 0x4e,
	0xdd, 0x7f, 0x91, 0xff, 0x28, 0x00, 0x1c, 0x4f, 0xce, 0x99, 0xa5, 0xa2, 0x6a, 0x0d, 0xe6, 0x2f,
	0x43, 0xe6, 0x5d, 0x71, 0xdd, 0xc3, 0x39, 0xca, 0x43, 0xb2, 0x01, 0xa2, 0x15, 0x06, 0x6a, 0x74,
	0xef, 0x5c, 0xff, 0x70, 0x8e, 0xa6, 0x19, 0xf2, 0x10, 0x16, 0xfd, 0xf0, 0x94, 0xdf, 0x74, 0xb4,
	0xa2, 0x18, 0xaf, 0x98, 0xc9, 0xe6, 0xc9, 0xef, 0x95, 0xe1, 0x96, 0xab, 0x06, 0x01, 0xf3, 0xec,
	0xa6, 0x06, 0x95, 0xa3, 0xec, 0x52, 0xc8, 0x1a, 0xcc, 0xab, 0xa6, 0xa1, 0xfa, 0x71, 0x07, 0x78,
	0x90, 0xd6, 0x5b, 0xf8, 0xdf, 0x7a, 0x55, 0x28, 0x1f, 0x18, 0xcc, 0xd4, 0xfe, 0xa0, 0xd1, 0x00,
	0x91, 0x3b, 0x5c, 0x35, 0x51, 0x47, 0xa4, 0x69, 0x4c, 0x96, 0xa1, 0x60, 0x68, 0x58, 0x9f, 0x48,
	0x0b, 0x86, 0x96, 0x2b, 0x21, 0x81, 0xd8, 0xb1, 0x43, 0xeb, 0xf7, 0x0a, 0xcd, 0x15, 0x58, 0x8a,
	0x56, 0xbc, 0x52, 0xcd, 0x10, 0x8b, 0xdd, 0x32, 0x01, 0x32, 0xb3, 0x90, 0x75, 0xb8, 0xdd, 0x1f,
	0x74, 0xe8, 0xee, 0xb0, 0xdb, 0xef, 0x8d, 0x47, 0xbd, 0xe7, 0xbd, 0xfe, 0xeb, 0x5e, 0x75, 0x8e,
	0xac, 0xc2, 0x4a, 0x96, 0x7e, 0x39, 0xea, 0xd0, 0x93, 0xaa, 0x40, 0x6a, 0x40, 0xb2, 0xe4, 0xd1,
	0x68, 0x88, 0x83, 0x6a, 0x81, 0x34, 0xa0, 0x96, 0xe5, 0x8f, 0x47, 0x7b, 0xc7, 0xfb, 0xb4, 0x3b,
	0xc0, 0xb9, 0xe2, 0xd6, 0xd3, 0xb4, 0xd3, 0x28, 0x57, 0x85, 0xc5, 0xe1, 0xc9, 0xa0, 0x73, 0x4d,
	0x69, 0x09, 0xca, 0x98, 0x89, 0x7e, 0xaa, 0x02, 0x59, 0x06, 0xc0, 0xb0, 0xdb, 0x1b, 0x8c, 0x86,
	0xd5, 0x42, 0x7b, 0x00, 0xa2, 0xc7, 0x7c, 0xc7, 0x9c, 0x32, 0x8f, 0x6c, 0xc8, 0xfc, 0xc9, 0xcb,
	0xc9, 0x93, 0x97, 0xa3, 0xd7, 0xdd, 0xe7, 0x7f, 0x21, 0xea, 0x9f, 0x3f, 0x44, 0xad, 0xaa, 0xec,
	0xac, 0x5d, 0xbf, 0xa8, 0xe4, 0xf9, 0xd3, 0x94, 0xd2, 0x1e, 0xc0, 0x82, 0x8f, 0x26, 0x24, 0xf7,
	0x6f, 0xf0, 0x78, 0x63, 0x12, 0xe2, 0x97, 0x98, 0x38, 0xf3, 0xd0, 0x32, 0x03, 0xd3, 0x98, 0xd3,
	0xee, 0x71, 0xa3, 0x90, 0xcd, 0x1c, 0x1e, 0x96, 0x9e, 0x00, 0xbf, 0xc6, 0xc0, 0x3c, 0x2f, 0x21,
	0x11, 0x39, 0xed, 0x17, 0x30, 0x7f, 0x16, 0xf9, 0x86, 0xdc, 0xcb, 0x29, 0x98, 0x99, 0xe9, 0xf9,
	0xbe, 0xc5, 0xb8, 0xf5, 0xd9, 0x8a, 0x63, 0xc7, 0x51, 0x0e, 0x69, 0x1f, 0x42, 0x89, 0xd9, 0xa1,
	0x95, 0xd3, 0xbd, 0xc8, 0x17, 0x09, 0xeb, 0x7b, 0x5e, 0xf7, 0x12, 0x6b, 0x51, 0x24, 0xb4, 0xdf,
	0x02, 0x44, 0xdf, 0xf1, 0x34, 0xf2, 0x12, 0x79, 0x90, 0xcb, 0x43, 0x9f, 0x25, 0xd0, 0x1f, 0x31,
	0xf4, 0xee, 0xaf, 0xd0, 0xd4, 0x8d, 0xb4, 0xcc, 0x92, 0x70, 0x6f, 0xff, 0xcd, 0xae, 0x6e, 0x04,
	0xe7, 0xe1, 0xa9, 0x3c, 0x71, 0x2c, 0x65, 0x3a, 0x55, 0x2f, 0x54, 0x8b, 0x29, 0x81, 0x61, 0x2a,
	0xba, 0xe7, 0x4e, 0xf0, 0x67, 0x9b, 0x73, 0x94, 0xeb, 0xff, 0x22, 0x9e, 0xf0, 0x60, 0xcc, 0x8f,
	0xb1, 0x80, 0x9f, 0xc7, 0x3f, 0x03, 0x00, 0x00, 0xff, 0xff, 0xc1, 0x8f, 0x10, 0x6f, 0x46, 0x06,
	0x00, 0x00,
}
