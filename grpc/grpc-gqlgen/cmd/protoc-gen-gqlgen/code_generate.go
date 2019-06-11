package main

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/k0kubun/pp"
	proto_extentions "github.com/vvakame/til/grpc/grpc-gqlgen/gqlgen-proto"
)

type GraphQLOperationType int

const (
	GraphQLQuery GraphQLOperationType = iota
	GraphQLMutation
	GraphQLSubscription
)

type Builder struct {
	CurrentRequest  *plugin.CodeGeneratorRequest
	CurrentFile     *descriptor.FileDescriptorProto
	CurrentFileRule *proto_extentions.FileRule
	CurrentService  *descriptor.ServiceDescriptorProto
	CurrentMethod   *descriptor.MethodDescriptorProto
	CurrentMessage  *descriptor.DescriptorProto
	CurrentField    *descriptor.FieldDescriptorProto

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

	files := make(map[string]*descriptor.FileDescriptorProto)
	for _, f := range req.ProtoFile {
		files[f.GetName()] = f
	}

	resp := &plugin.CodeGeneratorResponse{}

	for _, fname := range req.FileToGenerate {
		f := files[fname]

		fileInfo, err := b.GenerateFileInfo(ctx, f)
		if err != nil {
			return nil, err
		}

		b.FileInfos = append(b.FileInfos, fileInfo)
	}

	pp.Println(b.FileInfos)

	return resp, nil
}

func (b *Builder) GenerateFileInfo(ctx context.Context, f *descriptor.FileDescriptorProto) (*FileInfo, error) {
	b.CurrentFile = f
	defer func() {
		b.CurrentFile = nil
	}()

	fileInfo := &FileInfo{
		PackageName:    f.GetPackage(),
		ProtoGoPackage: f.GetOptions().GetGoPackage(),
	}

	if opts := f.GetOptions(); opts != nil {
		ext, err := proto.GetExtension(opts, proto_extentions.E_Resolver)
		if err == proto.ErrMissingExtension {
			// ok
		} else if err != nil {
			return nil, err
		} else {
			v := ext.(*proto_extentions.FileRule)
			rules, err := b.GenerateInferenceRules(ctx, v)
			if err != nil {
				return nil, err
			}
			fileInfo.InferrenceRules = rules
		}
	}

	for _, srvc := range f.GetService() {
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

func (b *Builder) GenerateServiceInfo(ctx context.Context, req *descriptor.ServiceDescriptorProto) (*ServiceInfo, error) {
	b.CurrentService = req
	defer func() {
		b.CurrentService = nil
	}()

	service := &ServiceInfo{
		Name: req.GetName(),
	}

	for _, m := range req.GetMethod() {
		methodInfo, err := b.GenerateMethodInfo(ctx, m)
		if err != nil {
			return nil, err
		}

		service.Methods = append(service.Methods, methodInfo)
	}

	return service, nil
}

func (b *Builder) GenerateMethodInfo(ctx context.Context, req *descriptor.MethodDescriptorProto) (*MethodInfo, error) {
	b.CurrentMethod = req
	defer func() {
		b.CurrentMethod = nil
	}()

	method := &MethodInfo{
		Name: req.GetName(),
	}

	if opts := req.GetOptions(); opts != nil {
		ext, err := proto.GetExtension(opts, proto_extentions.E_Schema)
		if err == proto.ErrMissingExtension {
			// ok
		} else if err != nil {
			return nil, err
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

	for _, m := range b.CurrentFile.GetMessageType() {
		messageName := fmt.Sprintf(".%s.%s", b.CurrentFile.GetPackage(), m.GetName())
		if messageName == req.GetInputType() {
			messageInfo, err := b.GenerateMessageInfo(ctx, m)
			if err != nil {
				return nil, err
			}

			method.RequestMessage = messageInfo
		}
		if messageName == req.GetOutputType() {
			messageInfo, err := b.GenerateMessageInfo(ctx, m)
			if err != nil {
				return nil, err
			}

			method.ResponseMessage = messageInfo
		}
	}

	if method.RequestMessage == nil {
		return nil, fmt.Errorf("request message doesn't lookup in method %s", req.GetName())
	}
	if method.ResponseMessage == nil {
		return nil, fmt.Errorf("response message doesn't lookup in method %s", req.GetName())
	}

	return method, nil
}

func (b *Builder) GenerateMessageInfo(ctx context.Context, req *descriptor.DescriptorProto) (*MessageInfo, error) {
	b.CurrentMessage = req
	defer func() {
		b.CurrentMessage = nil
	}()

	messageInfo := &MessageInfo{
		Name: req.GetName(),
	}

	if opts := req.GetOptions(); opts != nil {
		ext, err := proto.GetExtension(opts, proto_extentions.E_Type)
		if err == proto.ErrMissingExtension {
			// ok
		} else if err != nil {
			return nil, err
		} else {
			v := ext.(*proto_extentions.MessageRule)
			messageInfo.GraphQLAlias = v.GetAlias()
			messageInfo.GraphQLMessageType = v.GetType()
		}
	}

	for _, f := range req.GetField() {
		fieldInfo, err := b.GenerateFieldInfo(ctx, f)
		if err != nil {
			return nil, err
		}
		messageInfo.Fields = append(messageInfo.Fields, fieldInfo)
	}

	return messageInfo, nil
}

func (b *Builder) GenerateFieldInfo(ctx context.Context, req *descriptor.FieldDescriptorProto) (*FieldInfo, error) {
	b.CurrentField = req
	defer func() {
		b.CurrentField = nil
	}()

	fieldInfo := &FieldInfo{
		Name: req.GetName(),
	}

	if opts := req.GetOptions(); opts != nil {
		ext, err := proto.GetExtension(opts, proto_extentions.E_Field)
		if err == proto.ErrMissingExtension {
			// ok
		} else if err != nil {
			return nil, err
		} else {
			v := ext.(*proto_extentions.FieldRule)
			fieldInfo.GraphQLID = v.GetId()
			fieldInfo.GraphQLAlias = v.GetAlias()
			fieldInfo.GraphQLOptional = v.GetOptional()
		}
	}

	return fieldInfo, nil
}
