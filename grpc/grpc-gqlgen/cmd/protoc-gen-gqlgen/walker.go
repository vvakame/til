package main

import (
	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	descriptor "github.com/jhump/protoreflect/desc"
	gqlgen_proto "github.com/vvakame/til/grpc/grpc-gqlgen/gqlgen-proto"
	"golang.org/x/xerrors"
	"reflect"
)

type Visitor interface {
	VisitFileDescriptor(w *Walker, fd *descriptor.FileDescriptor, opts *gqlgen_proto.FileRule) error
	VisitServiceDescriptor(w *Walker, sd *descriptor.ServiceDescriptor) error
	VisitMethodDescriptor(w *Walker, md *descriptor.MethodDescriptor, opts *gqlgen_proto.SchemaRule) error
	VisitInputMessageDescriptor(w *Walker, md *descriptor.MessageDescriptor, opts *gqlgen_proto.MessageRule) error
	VisitOutputMessageDescriptor(w *Walker, md *descriptor.MessageDescriptor, opts *gqlgen_proto.MessageRule) error
	VisitFieldDescriptor(w *Walker, fd *descriptor.FieldDescriptor, opts *gqlgen_proto.FieldRule) error
}

func Visit(req *plugin.CodeGeneratorRequest, v Visitor) error {
	w := &Walker{
		CurrentRequest: req,
	}
	return w.start(v)
}

type Walker struct {
	CurrentRequest  *plugin.CodeGeneratorRequest
	CurrentFile     *descriptor.FileDescriptor
	CurrentFileRule *gqlgen_proto.FileRule
	CurrentService  *descriptor.ServiceDescriptor
	CurrentMethod   *descriptor.MethodDescriptor
	CurrentMessage  *descriptor.MessageDescriptor
	CurrentField    *descriptor.FieldDescriptor
}

func (w *Walker) start(v Visitor) error {

	req := w.CurrentRequest

	fdMap, err := descriptor.CreateFileDescriptors(req.GetProtoFile())
	if err != nil {
		return xerrors.Errorf("on descriptor.CreateFileDescriptors: %w", err)
	}

	for _, fname := range req.GetFileToGenerate() {
		f := fdMap[fname]

		err := w.visitFileDescriptor(v, f)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *Walker) visitFileDescriptor(v Visitor, req *descriptor.FileDescriptor) error {
	bk := w.CurrentFile
	w.CurrentFile = req
	defer func() {
		w.CurrentFile = bk
	}()

	var optVal *gqlgen_proto.FileRule
	opts := req.GetOptions()
	if opts != nil && !isNilPtr(opts) {
		ext, err := proto.GetExtension(opts, gqlgen_proto.E_Resolver)
		if xerrors.Is(err, proto.ErrMissingExtension) {
			// ok
		} else if err != nil {
			return xerrors.Errorf("%s on proto.GetExtension in visitFileDescriptor: %w", req.GetFullyQualifiedName(), err)
		} else {
			optVal = ext.(*gqlgen_proto.FileRule)
		}
	}

	err := v.VisitFileDescriptor(w, req, optVal)
	if err != nil {
		return err
	}

	for _, srvc := range req.GetServices() {
		err := w.visitServiceDescriptor(v, srvc)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *Walker) visitServiceDescriptor(v Visitor, req *descriptor.ServiceDescriptor) error {
	bk := w.CurrentService
	w.CurrentService = req
	defer func() {
		w.CurrentService = bk
	}()

	err := v.VisitServiceDescriptor(w, req)
	if err != nil {
		return err
	}

	for _, method := range req.GetMethods() {
		err := w.visitMethodDescriptor(v, method)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *Walker) visitMethodDescriptor(v Visitor, req *descriptor.MethodDescriptor) error {
	bk := w.CurrentMethod
	w.CurrentMethod = req
	defer func() {
		w.CurrentMethod = bk
	}()

	var optVal *gqlgen_proto.SchemaRule
	opts := req.GetOptions()
	if opts != nil && !isNilPtr(opts) {
		ext, err := proto.GetExtension(opts, gqlgen_proto.E_Schema)
		if xerrors.Is(err, proto.ErrMissingExtension) {
			// ok
		} else if err != nil {
			return xerrors.Errorf("%s on proto.GetExtension in visitMethodDescriptor: %w", req.GetFullyQualifiedName(), err)
		} else {
			optVal = ext.(*gqlgen_proto.SchemaRule)
		}
	}

	err := v.VisitMethodDescriptor(w, req, optVal)
	if err != nil {
		return err
	}

	{
		err := w.visitMessageDescriptor(v, req.GetInputType(), true)
		if err != nil {
			return err
		}
	}
	{
		err := w.visitMessageDescriptor(v, req.GetOutputType(), false)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *Walker) visitMessageDescriptor(v Visitor, req *descriptor.MessageDescriptor, isInput bool) error {
	bk := w.CurrentMessage
	w.CurrentMessage = req
	defer func() {
		w.CurrentMessage = bk
	}()

	var optVal *gqlgen_proto.MessageRule
	opts := req.GetOptions()
	if opts != nil && !isNilPtr(opts) {
		ext, err := proto.GetExtension(opts, gqlgen_proto.E_Type)
		if xerrors.Is(err, proto.ErrMissingExtension) {
			// ok
		} else if err != nil {
			return xerrors.Errorf("%s on proto.GetExtension in GenerateMessageInfo: %w", req.GetFullyQualifiedName(), err)
		} else {
			optVal = ext.(*gqlgen_proto.MessageRule)
		}
	}

	if isInput {
		err := v.VisitInputMessageDescriptor(w, req, optVal)
		if err != nil {
			return err
		}
	} else {
		err := v.VisitOutputMessageDescriptor(w, req, optVal)
		if err != nil {
			return err
		}
	}

	for _, f := range req.GetFields() {
		err := w.visitFieldDescriptor(v, f)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *Walker) visitFieldDescriptor(v Visitor, req *descriptor.FieldDescriptor) error {
	bk := w.CurrentField
	w.CurrentField = req
	defer func() {
		w.CurrentField = bk
	}()

	var optVal *gqlgen_proto.FieldRule
	opts := req.GetOptions()
	if opts != nil && !isNilPtr(opts) {
		ext, err := proto.GetExtension(opts, gqlgen_proto.E_Field)
		if xerrors.Is(err, proto.ErrMissingExtension) {
			// ok
		} else if err != nil {
			return xerrors.Errorf("%s on proto.GetExtension in GenerateFieldInfo: %w", req.GetFullyQualifiedName(), err)
		} else {
			optVal = ext.(*gqlgen_proto.FieldRule)
		}
	}

	err := v.VisitFieldDescriptor(w, req, optVal)
	if err != nil {
		return err
	}

	return nil
}

func isNilPtr(x interface{}) bool {
	v := reflect.ValueOf(x)
	return v.Kind() == reflect.Ptr && v.IsNil()
}
