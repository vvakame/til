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
	gqlgen_proto "github.com/vvakame/til/grpc/grpc-gqlgen/gqlgen-proto"
	"golang.org/x/xerrors"
	"io/ioutil"
	"os"
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
	FileInfos  []*FileInfo

	CurrentFileInfo    *FileInfo
	CurrentServiceInfo *ServiceInfo
	CurrentMethodInfo  *MethodInfo
	CurrentMessageInfo *MessageInfo
	CurrentFieldInfo   *FieldInfo
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
	Name string

	GraphQLAlias       string
	GraphQLMessageType gqlgen_proto.MessageType

	Fields []*FieldInfo
}

type FieldInfo struct {
	Name string

	GraphQLAlias    string
	GraphQLOptional bool
	GraphQLID       bool
}

func (b *Builder) Process(ctx context.Context, req *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {

	err := Visit(req, b)
	if err != nil {
		return nil, err
	}

	resp := &plugin.CodeGeneratorResponse{}

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

func (b *Builder) VisitFileDescriptor(w *Walker, req *descriptor.FileDescriptor, opts *gqlgen_proto.FileRule) error {

	fileInfo := &FileInfo{
		PackageName:    req.GetPackage(),
		ProtoGoPackage: req.GetFileOptions().GetGoPackage(),
	}
	b.CurrentFileInfo = fileInfo

	for _, v := range opts.GetTypeInference() {
		fileInfo.InferrenceRules = append(fileInfo.InferrenceRules, &InferrenceRule{
			Src:         v.GetSrc(),
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

func (b *Builder) VisitInputMessageDescriptor(w *Walker, req *descriptor.MessageDescriptor, opts *gqlgen_proto.MessageRule) error {
	err := b.visitMessageDescriptor(w, req, opts)
	if err != nil {
		return err
	}

	b.CurrentMethodInfo.RequestMessage = b.CurrentMessageInfo

	return nil
}

func (b *Builder) VisitOutputMessageDescriptor(w *Walker, req *descriptor.MessageDescriptor, opts *gqlgen_proto.MessageRule) error {
	err := b.visitMessageDescriptor(w, req, opts)
	if err != nil {
		return err
	}

	b.CurrentMethodInfo.ResponseMessage = b.CurrentMessageInfo

	return nil
}

func (b *Builder) visitMessageDescriptor(w *Walker, req *descriptor.MessageDescriptor, opts *gqlgen_proto.MessageRule) error {
	messageInfo := &MessageInfo{
		Name: req.GetName(),
	}
	b.CurrentMessageInfo = messageInfo
	if opts != nil {
		messageInfo.GraphQLAlias = opts.GetAlias()
		messageInfo.GraphQLMessageType = opts.GetType()
	}

	return nil
}

func (b *Builder) VisitFieldDescriptor(w *Walker, req *descriptor.FieldDescriptor, opts *gqlgen_proto.FieldRule) error {
	fieldInfo := &FieldInfo{
		Name: req.GetName(),
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
