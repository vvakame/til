package main

import (
	"context"
	"reflect"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	descriptor "github.com/jhump/protoreflect/desc"
	"github.com/k0kubun/pp"
	proto_extentions "github.com/vvakame/til/grpc/grpc-gqlgen/gqlgen-proto"
	"golang.org/x/xerrors"
)

type GraphQLOperationType int

const (
	GraphQLQuery GraphQLOperationType = iota
	GraphQLMutation
	GraphQLSubscription
)

type Builder struct {
	CurrentRequest  *plugin.CodeGeneratorRequest
	CurrentFile     *descriptor.FileDescriptor
	CurrentFileRule *proto_extentions.FileRule
	CurrentService  *descriptor.ServiceDescriptor
	CurrentMethod   *descriptor.MethodDescriptor
	CurrentMessage  *descriptor.MessageDescriptor
	CurrentField    *descriptor.FieldDescriptor

	FileInfos []*FileInfo
}

type FileInfo struct {
	PackageName    string
	ProtoGoPackage string

	InferrenceRules []*InferrenceRule

	Services []*ServiceInfo
}

type InferrenceRule struct {
	Src         string
	Dest        string
	MessageType proto_extentions.MessageType
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
	Name string

	GraphQLAlias       string
	GraphQLMessageType proto_extentions.MessageType

	Fields []*FieldInfo
}

type FieldInfo struct {
	Name string

	GraphQLAlias    string
	GraphQLOptional bool
	GraphQLID       bool
}

func (b *Builder) Process(ctx context.Context, req *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {
	b.CurrentRequest = req
	defer func() {
		b.CurrentRequest = nil
	}()

	resp := &plugin.CodeGeneratorResponse{}

	fdMap, err := descriptor.CreateFileDescriptors(req.GetProtoFile())
	if err != nil {
		return nil, xerrors.Errorf("on descriptor.CreateFileDescriptors: %w", err)
	}

	for _, fname := range req.FileToGenerate {
		f := fdMap[fname]

		fileInfo, err := b.GenerateFileInfo(ctx, f)
		if err != nil {
			return nil, xerrors.Errorf("%s on Builder.GenerateFileInfo: %w", fname, err)
		}

		b.FileInfos = append(b.FileInfos, fileInfo)
	}

	pp.Println(b.FileInfos)

	return resp, nil
}

func (b *Builder) GenerateFileInfo(ctx context.Context, req *descriptor.FileDescriptor) (*FileInfo, error) {
	b.CurrentFile = req
	defer func() {
		b.CurrentFile = nil
	}()

	fileInfo := &FileInfo{
		PackageName:    req.GetPackage(),
		ProtoGoPackage: req.GetFileOptions().GetGoPackage(),
	}

	if opts := req.GetOptions(); opts != nil && !isNilPtr(opts) {
		ext, err := proto.GetExtension(opts, proto_extentions.E_Resolver)
		if xerrors.Is(err, proto.ErrMissingExtension) {
			// ok
		} else if err != nil {
			return nil, xerrors.Errorf("%s on proto.GetExtension in GenerateFileInfo: %w", req.GetFullyQualifiedName(), err)
		} else {
			v := ext.(*proto_extentions.FileRule)
			rules, err := b.GenerateInferenceRules(ctx, v)
			if err != nil {
				return nil, err
			}
			fileInfo.InferrenceRules = rules
		}
	}

	for _, srvc := range req.GetServices() {
		serviceInfo, err := b.GenerateServiceInfo(ctx, srvc)
		if err != nil {
			return nil, err
		}

		fileInfo.Services = append(fileInfo.Services, serviceInfo)
	}

	return fileInfo, nil
}

func (b *Builder) GenerateInferenceRules(ctx context.Context, req *proto_extentions.FileRule) ([]*InferrenceRule, error) {
	b.CurrentFileRule = req
	defer func() {
		b.CurrentFileRule = nil
	}()

	var rules []*InferrenceRule
	for _, v := range req.GetTypeInference() {
		rules = append(rules, &InferrenceRule{
			Src:         v.GetSrc(),
			Dest:        v.GetDest(),
			MessageType: v.GetType(),
		})
	}

	return rules, nil
}

func (b *Builder) GenerateServiceInfo(ctx context.Context, req *descriptor.ServiceDescriptor) (*ServiceInfo, error) {
	b.CurrentService = req
	defer func() {
		b.CurrentService = nil
	}()

	service := &ServiceInfo{
		Name: req.GetName(),
	}

	for _, m := range req.GetMethods() {
		methodInfo, err := b.GenerateMethodInfo(ctx, m)
		if err != nil {
			return nil, err
		}

		service.Methods = append(service.Methods, methodInfo)
	}

	return service, nil
}

func (b *Builder) GenerateMethodInfo(ctx context.Context, req *descriptor.MethodDescriptor) (*MethodInfo, error) {
	b.CurrentMethod = req
	defer func() {
		b.CurrentMethod = nil
	}()

	method := &MethodInfo{
		Name: req.GetName(),
	}

	if opts := req.GetOptions(); opts != nil && !isNilPtr(opts) {
		ext, err := proto.GetExtension(opts, proto_extentions.E_Schema)
		if xerrors.Is(err, proto.ErrMissingExtension) {
			// ok
		} else if err != nil {
			return nil, xerrors.Errorf("%s on proto.GetExtension in GenerateMethodInfo: %w", req.GetFullyQualifiedName(), err)
		} else {
			v := ext.(*proto_extentions.SchemaRule)
			v.GetPattern()
			if v.GetQuery() != "" {
				method.GraphQLOperationType = GraphQLQuery
				method.Name = v.GetQuery()
			} else if v.GetMutation() != "" {
				method.GraphQLOperationType = GraphQLMutation
				method.Name = v.GetMutation()
			} else if v.GetSubscription() != "" {
				method.GraphQLOperationType = GraphQLSubscription
				method.Name = v.GetSubscription()
			}
		}
	}

	{
		messageInfo, err := b.GenerateMessageInfo(ctx, req.GetInputType())
		if err != nil {
			return nil, err
		}

		method.RequestMessage = messageInfo
	}
	{
		messageInfo, err := b.GenerateMessageInfo(ctx, req.GetOutputType())
		if err != nil {
			return nil, err
		}

		method.ResponseMessage = messageInfo
	}

	return method, nil
}

func (b *Builder) GenerateMessageInfo(ctx context.Context, req *descriptor.MessageDescriptor) (*MessageInfo, error) {
	b.CurrentMessage = req
	defer func() {
		b.CurrentMessage = nil
	}()

	messageInfo := &MessageInfo{
		Name: req.GetName(),
	}

	if opts := req.GetOptions(); opts != nil && !isNilPtr(opts) {
		ext, err := proto.GetExtension(opts, proto_extentions.E_Type)
		if xerrors.Is(err, proto.ErrMissingExtension) {
			// ok
		} else if err != nil {
			return nil, xerrors.Errorf("%s on proto.GetExtension in GenerateMessageInfo: %w", req.GetFullyQualifiedName(), err)
		} else {
			v := ext.(*proto_extentions.MessageRule)
			messageInfo.GraphQLAlias = v.GetAlias()
			messageInfo.GraphQLMessageType = v.GetType()
		}
	}

	for _, f := range req.GetFields() {
		fieldInfo, err := b.GenerateFieldInfo(ctx, f)
		if err != nil {
			return nil, err
		}
		messageInfo.Fields = append(messageInfo.Fields, fieldInfo)
	}

	return messageInfo, nil
}

func (b *Builder) GenerateFieldInfo(ctx context.Context, req *descriptor.FieldDescriptor) (*FieldInfo, error) {
	b.CurrentField = req
	defer func() {
		b.CurrentField = nil
	}()

	fieldInfo := &FieldInfo{
		Name: req.GetName(),
	}

	if opts := req.GetOptions(); opts != nil && !isNilPtr(opts) {
		ext, err := proto.GetExtension(opts, proto_extentions.E_Field)
		if xerrors.Is(err, proto.ErrMissingExtension) {
			// ok
		} else if err != nil {
			return nil, xerrors.Errorf("%s on proto.GetExtension in GenerateFieldInfo: %w", req.GetFullyQualifiedName(), err)
		} else {
			v := ext.(*proto_extentions.FieldRule)
			fieldInfo.GraphQLID = v.GetId()
			fieldInfo.GraphQLAlias = v.GetAlias()
			fieldInfo.GraphQLOptional = v.GetOptional()
		}
	}

	return fieldInfo, nil
}

func isNilPtr(x interface{}) bool {
	v := reflect.ValueOf(x)
	return v.Kind() == reflect.Ptr && v.IsNil()
}
