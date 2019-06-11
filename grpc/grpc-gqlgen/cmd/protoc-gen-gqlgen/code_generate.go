package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	descriptor "github.com/jhump/protoreflect/desc"
	"github.com/rakyll/statik/fs"
	proto_extentions "github.com/vvakame/til/grpc/grpc-gqlgen/gqlgen-proto"
	"golang.org/x/xerrors"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"text/template"
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

	tmplBytes, err := ioutil.ReadFile("./tmpls/glue.gotmpl")
	if os.IsNotExist(err) {
		statikFS, err := fs.New()
		if err != nil {
			return nil, xerrors.Errorf("on fs.New: %w", err)
		}
		f, err := statikFS.Open("/glue.gotmpl")
		if err != nil {
			return nil, xerrors.Errorf("on statikFS.Open: %w", err)
		}
		tmplBytes, err = ioutil.ReadAll(f)
		if err != nil {
			return nil, xerrors.Errorf("on ioutil.ReadAll: %w", err)
		}

	} else if err != nil {
		return nil, xerrors.Errorf("on read glue.gotmpl: %w", err)
	}
	tmpl, err := template.
		New("glue").
		Funcs(map[string]interface{}{
			"first": func(ss ...string) string {
				for _, s := range ss {
					if s != "" {
						return s
					}
				}
				return ""
			},
			"goName": func(name string) string {
				return templates.ToGo(name)
			},
			"goNamePrivate": func(name string) string {
				return templates.ToGoPrivate(name)
			},
		}).
		Parse(string(tmplBytes))
	if err != nil {
		return nil, xerrors.Errorf("on parse template: %w", err)
	}
	for _, fileInfo := range b.FileInfos {
		err = fileInfo.Prepare()
		if err != nil {
			return nil, xerrors.Errorf("%s on fileInfo.Prepare: %w", fileInfo.PackageName, err)
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, fileInfo)
		if err != nil {
			return nil, xerrors.Errorf("%s on tmpl.Execute: %w", fileInfo.PackageName, err)
		}

		resp.File = append(resp.File, &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(fmt.Sprintf("%s.gql.go", fileInfo.PackageName)),
			Content: proto.String(buf.String()),
		})
	}

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
				method.GraphQLName = v.GetQuery()
			} else if v.GetMutation() != "" {
				method.GraphQLOperationType = GraphQLMutation
				method.GraphQLName = v.GetMutation()
			} else if v.GetSubscription() != "" {
				method.GraphQLOperationType = GraphQLSubscription
				method.GraphQLName = v.GetSubscription()
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
