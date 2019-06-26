package main

import (
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
	FileInfos  []*FileInfo
	SchemaDocs []*ast.SchemaDocument

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

	InferrenceRules []*InferrenceRule

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

func (fi *FileInfo) GoPackageName() string {
	ss := strings.SplitN(fi.ProtoGoPackage, ";", 2)
	if len(ss) == 2 {
		return ss[1]
	}

	ss = strings.Split(fi.ProtoGoPackage, "/")
	return ss[len(ss)-1]
}

type InferrenceRule struct {
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
	GraphQLName          string
}

type MessageInfo struct {
	Proto *descriptor.MessageDescriptor

	Name string

	GraphQLAlias       string
	GraphQLMessageType gqlgen_proto.MessageType

	Fields []*FieldInfo
}

func (m *MessageInfo) GraphQLName() string {
	name := m.GraphQLAlias
	if name == "" {
		name = m.Name
	}
	return templates.ToGo(name)
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
	name := f.GraphQLAlias
	if name == "" {
		name = f.Name
	}
	return templates.ToGoPrivate(name)
}

type EnumInfo struct {
	Proto *descriptor.EnumDescriptor

	Name string

	GraphQLAlias string

	Values []*EnumValueInfo
}

func (e *EnumInfo) GraphQLName() string {
	name := e.GraphQLAlias
	if name == "" {
		name = e.Name
	}
	return templates.ToGo(name)
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
	for _, fileInfo := range b.FileInfos {
		for _, messageInfo := range fileInfo.MessageInfos {
			if messageInfo.Proto.GetFullyQualifiedName() == name {
				return messageInfo
			}
		}
	}

	return nil
}

func (b *Builder) FindEnumInfo(name string) *EnumInfo {
	for _, fileInfo := range b.FileInfos {
		for _, enumInfo := range fileInfo.EnumInfos {
			if enumInfo.Proto.GetFullyQualifiedName() == name {
				return enumInfo
			}
		}
	}

	return nil
}

func (b *Builder) VisitFileDescriptor(w *Walker, req *descriptor.FileDescriptor, opts *gqlgen_proto.FileRule) error {

	fileInfo := &FileInfo{
		Proto:          req,
		PackageName:    req.GetPackage(),
		ProtoGoPackage: req.GetFileOptions().GetGoPackage(),
	}
	b.CurrentFileInfo = fileInfo

	for _, v := range opts.GetTypeInference() {
		src, err := regexp.Compile(v.GetSrc())
		if err != nil {
			return err
		}
		fileInfo.InferrenceRules = append(fileInfo.InferrenceRules, &InferrenceRule{
			Src:         src,
			Dest:        v.GetDest(),
			MessageType: v.GetType(),
		})
	}

	b.FileInfos = append(b.FileInfos, fileInfo)

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
	if opts != nil {
		if opts.GetQuery() != "" {
			method.GraphQLOperationType = GraphQLQuery
			method.GraphQLName = opts.GetQuery()
		} else if opts.GetMutation() != "" {
			method.GraphQLOperationType = GraphQLMutation
			method.GraphQLName = opts.GetMutation()
		} else if opts.GetSubscription() != "" {
			method.GraphQLOperationType = GraphQLSubscription
			method.GraphQLName = opts.GetSubscription()
		}
	}

	b.CurrentServiceInfo.Methods = append(b.CurrentServiceInfo.Methods, method)

	return nil
}

func (b *Builder) VisitMessageDescriptor(w *Walker, req *descriptor.MessageDescriptor, opts *gqlgen_proto.MessageRule, info *VisitMessageInfo) error {
	messageInfo := &MessageInfo{
		Proto: req,
		Name:  req.GetName(),
	}
	b.CurrentMessageInfo = messageInfo

	for _, rule := range b.CurrentFileInfo.InferrenceRules {
		ss := rule.Src.FindStringSubmatch(req.GetName())
		if len(ss) == 0 {
			continue
		}

		name := rule.Dest
		for idx, s := range ss[1:] {
			name = strings.Replace(name, fmt.Sprintf("$%d", idx+1), s, 1)
		}
		messageInfo.GraphQLAlias = name

		switch rule.MessageType {
		case gqlgen_proto.MessageType_UNKNOWN:
		// ignore
		case gqlgen_proto.MessageType_TYPE,
			gqlgen_proto.MessageType_INPUT:
			messageInfo.GraphQLMessageType = rule.MessageType
		}
	}

	if opts != nil {
		if v := opts.GetAlias(); v != "" {
			messageInfo.GraphQLAlias = v
		}
		if v := opts.GetType(); v != gqlgen_proto.MessageType_UNKNOWN {
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
		Proto: req,
		Name:  req.GetName(),
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
