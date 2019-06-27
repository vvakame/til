package main

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/99designs/gqlgen/codegen/templates"
	basedescriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	descriptor "github.com/jhump/protoreflect/desc"
	"github.com/vektah/gqlparser/ast"
	gqlgen_proto "github.com/vvakame/til/grpc/grpc-gqlgen/gqlgen-proto"
	"golang.org/x/xerrors"
)

type GraphQLOperationType int

const (
	GraphQLQuery GraphQLOperationType = iota
	GraphQLMutation
	GraphQLSubscription
)

type Builder struct {
	AllFileInfos      []*FileInfo
	GenerateFileInfos []*FileInfo
	SchemaDocs        []*ast.SchemaDocument

	CurrentFileInfo      *FileInfo
	CurrentServiceInfo   *ServiceInfo
	CurrentMethodInfo    *MethodInfo
	CurrentMessageInfo   *MessageInfo
	CurrentFieldInfo     *FieldInfo
	CurrentEnumInfo      *EnumInfo
	CurrentEnumValueInfo *EnumValueInfo
}

type FileInfo struct {
	Proto *descriptor.FileDescriptor

	PackageName    string
	ProtoGoPackage string

	MethodRules  []*MethodRule
	MessageRules []*MessageRule

	Services     []*ServiceInfo
	MessageInfos []*MessageInfo
	EnumInfos    []*EnumInfo
}

func (fi *FileInfo) Prepare() error {
	if fi.ProtoGoPackage == "" {
		return xerrors.New("ProtoGoPackage is nil")
	}

	return nil
}

func (fi *FileInfo) GoImportPath() string {
	ss := strings.SplitN(fi.ProtoGoPackage, ";", 2)
	if len(ss) == 2 {
		return ss[0]
	}

	return fi.ProtoGoPackage
}

func (fi *FileInfo) GoPackageName() string {
	ss := strings.SplitN(fi.ProtoGoPackage, ";", 2)
	if len(ss) == 2 {
		return ss[1]
	}

	ss = strings.Split(fi.ProtoGoPackage, "/")
	return ss[len(ss)-1]
}

type MethodRule struct {
	Src        *regexp.Regexp
	Dest       string
	MethodType gqlgen_proto.MethodType
}

type MessageRule struct {
	Src         *regexp.Regexp
	Dest        string
	MessageType gqlgen_proto.MessageType
}

type ServiceInfo struct {
	Name string

	Methods []*MethodInfo
}

type MethodInfo struct {
	Name            string
	RequestMessage  *MessageInfo
	ResponseMessage *MessageInfo

	GraphQLOperationType GraphQLOperationType
	GraphQLAlias         string
}

func (m *MethodInfo) GraphQLName() string {
	if m.GraphQLAlias != "" {
		return m.GraphQLAlias
	}

	return templates.ToGoPrivate(m.Name)
}

type MessageInfo struct {
	ParentMessage *MessageInfo
	Proto         *descriptor.MessageDescriptor

	Name string

	GraphQLAlias       string
	GraphQLMessageType gqlgen_proto.MessageType

	Fields []*FieldInfo
}

func (m *MessageInfo) GraphQLName() string {
	if m.GraphQLAlias != "" {
		return m.GraphQLAlias
	}

	var buf bytes.Buffer

	var printName func(parent *MessageInfo)
	printName = func(parent *MessageInfo) {
		if parent == nil {
			return
		}
		printName(parent.ParentMessage)
		buf.WriteString(parent.GraphQLName())
	}
	printName(m.ParentMessage)
	buf.WriteString(templates.ToGo(m.Name))

	return buf.String()
}

func (m *MessageInfo) GoName() string {
	var buf bytes.Buffer

	var printName func(desc descriptor.Descriptor) bool
	printName = func(desc descriptor.Descriptor) bool {
		switch v := desc.(type) {
		case *descriptor.EnumDescriptor:
			if printName(v.GetParent()) {
				buf.WriteString("_")
			}
			buf.WriteString(v.GetName())
			return true

		case *descriptor.MessageDescriptor:
			if printName(v.GetParent()) {
				buf.WriteString("_")
			}
			buf.WriteString(v.GetName())
			return true

		default:
			return false
		}
	}
	printName(m.Proto)

	return buf.String()
}

type FieldInfo struct {
	Name        string
	Type        basedescriptor.FieldDescriptorProto_Type
	Repeated    bool
	TypeMessage *descriptor.MessageDescriptor
	TypeEnum    *descriptor.EnumDescriptor

	GraphQLAlias    string
	GraphQLOptional bool
	GraphQLID       bool
}

func (f *FieldInfo) GraphQLName() string {
	if f.GraphQLAlias != "" {
		return f.GraphQLAlias
	}

	return templates.ToGoPrivate(f.Name)
}

type EnumInfo struct {
	ParentMessage *MessageInfo
	Proto         *descriptor.EnumDescriptor

	Name string

	GraphQLAlias string

	Values []*EnumValueInfo
}

func (e *EnumInfo) GraphQLName() string {
	if e.GraphQLAlias != "" {
		return e.GraphQLAlias
	}

	var buf bytes.Buffer

	var printName func(parent *MessageInfo)
	printName = func(parent *MessageInfo) {
		if parent == nil {
			return
		}
		printName(parent.ParentMessage)
		buf.WriteString(parent.GraphQLName())
	}
	printName(e.ParentMessage)
	buf.WriteString(templates.ToGo(e.Name))

	return buf.String()
}

func (e *EnumInfo) GoName() string {
	var buf bytes.Buffer

	var printName func(desc descriptor.Descriptor) bool
	printName = func(desc descriptor.Descriptor) bool {
		switch v := desc.(type) {
		case *descriptor.EnumDescriptor:
			if printName(v.GetParent()) {
				buf.WriteString("_")
			}
			buf.WriteString(v.GetName())
			return true

		case *descriptor.MessageDescriptor:
			if printName(v.GetParent()) {
				buf.WriteString("_")
			}
			buf.WriteString(v.GetName())
			return true

		default:
			return false
		}
	}
	printName(e.Proto)

	return buf.String()
}

type EnumValueInfo struct {
	Name string

	GraphQLAlias string
}

func (v *EnumValueInfo) GraphQLName() string {
	name := v.GraphQLAlias
	if name == "" {
		name = v.Name
	}
	// protoとGraphQLのenumのnaming ruleがまぁ同じだよなという仮定
	return name
}

func (b *Builder) Process(ctx context.Context, req *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {

	err := Visit(req, b)
	if err != nil {
		return nil, err
	}

	resp := &plugin.CodeGeneratorResponse{}

	{
		files, err := glueGenerate(ctx, b)
		if err != nil {
			return nil, xerrors.Errorf("on glueGenerate: %w", err)
		}
		resp.File = append(resp.File, files...)
	}
	{
		files, err := schemaGenerate(ctx, b)
		if err != nil {
			return nil, xerrors.Errorf("on schemaGenerate: %w", err)
		}
		resp.File = append(resp.File, files...)
	}

	return resp, nil
}

func (b *Builder) FindMessageInfo(name string) *MessageInfo {
	for _, fileInfo := range b.AllFileInfos {
		for _, messageInfo := range fileInfo.MessageInfos {
			if messageInfo.Proto.GetFullyQualifiedName() == name {
				return messageInfo
			}
		}
	}

	return nil
}

func (b *Builder) FindEnumInfo(name string) *EnumInfo {
	for _, fileInfo := range b.AllFileInfos {
		for _, enumInfo := range fileInfo.EnumInfos {
			if enumInfo.Proto.GetFullyQualifiedName() == name {
				return enumInfo
			}
		}
	}

	return nil
}

func (b *Builder) VisitFileDescriptor(w *Walker, req *descriptor.FileDescriptor, opts *gqlgen_proto.FileRule, info *VisitFileInfo) error {

	fileInfo := &FileInfo{
		Proto:          req,
		PackageName:    req.GetPackage(),
		ProtoGoPackage: req.GetFileOptions().GetGoPackage(),
	}
	b.CurrentFileInfo = fileInfo

	for _, v := range opts.GetMethodRule() {
		if v.GetSrc() == "" {
			return xerrors.New("src value is required in method_rule")
		}
		src, err := regexp.Compile(v.GetSrc())
		if err != nil {
			return err
		}
		fileInfo.MethodRules = append(fileInfo.MethodRules, &MethodRule{
			Src:        src,
			Dest:       v.GetDest(),
			MethodType: v.GetType(),
		})
	}

	for _, v := range opts.GetMessageRule() {
		if v.GetSrc() == "" {
			return xerrors.New("src value is required in message_rule")
		}
		src, err := regexp.Compile(v.GetSrc())
		if err != nil {
			return err
		}
		fileInfo.MessageRules = append(fileInfo.MessageRules, &MessageRule{
			Src:         src,
			Dest:        v.GetDest(),
			MessageType: v.GetType(),
		})
	}

	b.AllFileInfos = append(b.AllFileInfos, fileInfo)
	if info.IsGenerate {
		b.GenerateFileInfos = append(b.GenerateFileInfos, fileInfo)
	}

	return nil
}

func (b *Builder) VisitServiceDescriptor(w *Walker, req *descriptor.ServiceDescriptor) error {
	service := &ServiceInfo{
		Name: req.GetName(),
	}
	b.CurrentServiceInfo = service

	b.CurrentFileInfo.Services = append(b.CurrentFileInfo.Services, service)

	return nil
}

func (b *Builder) VisitMethodDescriptor(w *Walker, req *descriptor.MethodDescriptor, opts *gqlgen_proto.SchemaRule) error {
	method := &MethodInfo{
		Name: req.GetName(),
	}
	b.CurrentMethodInfo = method

	for _, rule := range b.CurrentFileInfo.MethodRules {
		ss := rule.Src.FindStringSubmatch(req.GetName())
		if len(ss) == 0 {
			continue
		}

		if name := rule.Dest; name != "" {
			for idx, s := range ss[1:] {
				name = strings.Replace(name, fmt.Sprintf("$%d", idx+1), s, 1)
			}
			method.GraphQLAlias = templates.ToGoPrivate(name)
		}

		switch rule.MethodType {
		case gqlgen_proto.MethodType_OPERATION_QUERY:
			method.GraphQLOperationType = GraphQLQuery
		case gqlgen_proto.MethodType_OPERATION_MUTATION:
			method.GraphQLOperationType = GraphQLMutation
		case gqlgen_proto.MethodType_OPERATION_SUBSCRIPTION:
			method.GraphQLOperationType = GraphQLSubscription
		}

		break
	}

	if opts != nil {
		if opts.GetQuery() != "" {
			method.GraphQLOperationType = GraphQLQuery
			method.GraphQLAlias = opts.GetQuery()
		} else if opts.GetMutation() != "" {
			method.GraphQLOperationType = GraphQLMutation
			method.GraphQLAlias = opts.GetMutation()
		} else if opts.GetSubscription() != "" {
			method.GraphQLOperationType = GraphQLSubscription
			method.GraphQLAlias = opts.GetSubscription()
		}
	}

	b.CurrentServiceInfo.Methods = append(b.CurrentServiceInfo.Methods, method)

	return nil
}

func (b *Builder) VisitMessageDescriptor(w *Walker, req *descriptor.MessageDescriptor, opts *gqlgen_proto.MessageRule, info *VisitMessageInfo) error {
	messageInfo := &MessageInfo{
		ParentMessage: b.FindMessageInfo(req.GetParent().GetFullyQualifiedName()),
		Proto:         req,
		Name:          req.GetName(),
	}
	b.CurrentMessageInfo = messageInfo

	for _, rule := range b.CurrentFileInfo.MessageRules {
		ss := rule.Src.FindStringSubmatch(req.GetName())
		if len(ss) == 0 {
			continue
		}

		if name := rule.Dest; name != "" {
			for idx, s := range ss[1:] {
				name = strings.Replace(name, fmt.Sprintf("$%d", idx+1), s, 1)
			}
			messageInfo.GraphQLAlias = templates.ToGo(name)
		}

		switch rule.MessageType {
		case gqlgen_proto.MessageType_TYPE_UNKNOWN:
		// ignore
		case gqlgen_proto.MessageType_TYPE_TYPE,
			gqlgen_proto.MessageType_TYPE_INPUT:
			messageInfo.GraphQLMessageType = rule.MessageType
		}

		break
	}

	if opts != nil {
		if v := opts.GetAlias(); v != "" {
			messageInfo.GraphQLAlias = v
		}
		if v := opts.GetType(); v != gqlgen_proto.MessageType_TYPE_UNKNOWN {
			messageInfo.GraphQLMessageType = v
		}
	}

	if info.IsInput {
		b.CurrentMethodInfo.RequestMessage = messageInfo
	} else if info.IsOutput {
		b.CurrentMethodInfo.ResponseMessage = messageInfo
	} else {
		b.CurrentFileInfo.MessageInfos = append(b.CurrentFileInfo.MessageInfos, messageInfo)
	}

	return nil
}

func (b *Builder) VisitFieldDescriptor(w *Walker, req *descriptor.FieldDescriptor, opts *gqlgen_proto.FieldRule) error {
	fieldInfo := &FieldInfo{
		Name:        req.GetName(),
		Type:        req.GetType(),
		Repeated:    req.GetLabel() == basedescriptor.FieldDescriptorProto_LABEL_REPEATED,
		TypeMessage: req.GetMessageType(),
		TypeEnum:    req.GetEnumType(),
	}
	b.CurrentFieldInfo = fieldInfo

	if opts != nil {
		fieldInfo.GraphQLID = opts.GetId()
		fieldInfo.GraphQLAlias = opts.GetAlias()
		fieldInfo.GraphQLOptional = opts.GetOptional()
	}

	b.CurrentMessageInfo.Fields = append(b.CurrentMessageInfo.Fields, fieldInfo)

	return nil
}

func (b *Builder) VisitEnumDescriptor(w *Walker, req *descriptor.EnumDescriptor, opts *gqlgen_proto.EnumRule) error {
	enumInfo := &EnumInfo{
		ParentMessage: b.FindMessageInfo(req.GetParent().GetFullyQualifiedName()),
		Proto:         req,
		Name:          req.GetName(),
	}
	b.CurrentEnumInfo = enumInfo

	if opts != nil {
		if v := opts.GetAlias(); v != "" {
			enumInfo.GraphQLAlias = v
		}
	}

	b.CurrentFileInfo.EnumInfos = append(b.CurrentFileInfo.EnumInfos, enumInfo)

	return nil
}

func (b *Builder) VisitEnumValueDescriptor(w *Walker, req *descriptor.EnumValueDescriptor, opts *gqlgen_proto.EnumValueRule) error {
	enumValueInfo := &EnumValueInfo{
		Name: req.GetName(),
	}
	b.CurrentEnumValueInfo = enumValueInfo

	if opts != nil {
		// none
	}

	b.CurrentEnumInfo.Values = append(b.CurrentEnumInfo.Values, enumValueInfo)

	return nil
}
