package main

import (
	"reflect"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	descriptor "github.com/jhump/protoreflect/desc"
	gqlgen_proto "github.com/vvakame/til/grpc/grpc-gqlgen/gqlgen-proto"
	"golang.org/x/xerrors"
)

type Visitor interface {
	VisitFileDescriptor(w *Walker, fd *descriptor.FileDescriptor, opts *gqlgen_proto.FileRule, info *VisitFileInfo) error
	VisitServiceDescriptor(w *Walker, sd *descriptor.ServiceDescriptor) error
	VisitMethodDescriptor(w *Walker, md *descriptor.MethodDescriptor, opts *gqlgen_proto.SchemaRule) error
	VisitMessageDescriptor(w *Walker, md *descriptor.MessageDescriptor, opts *gqlgen_proto.MessageRule, info *VisitMessageInfo) error
	VisitFieldDescriptor(w *Walker, fd *descriptor.FieldDescriptor, opts *gqlgen_proto.FieldRule) error
	VisitEnumDescriptor(w *Walker, ed *descriptor.EnumDescriptor, opts *gqlgen_proto.EnumRule) error
	VisitEnumValueDescriptor(w *Walker, enumValueDescriptor *descriptor.EnumValueDescriptor, opts *gqlgen_proto.EnumValueRule) error
}

type VisitFileInfo struct {
	IsGenerate bool
}

type VisitMessageInfo struct {
	IsInput  bool
	IsOutput bool
}

func Visit(req *plugin.CodeGeneratorRequest, v Visitor) error {
	w := &Walker{
		CurrentRequest: req,
	}
	return w.start(v)
}

type Walker struct {
	CurrentRequest   *plugin.CodeGeneratorRequest
	CurrentFile      *descriptor.FileDescriptor
	CurrentFileRule  *gqlgen_proto.FileRule
	CurrentService   *descriptor.ServiceDescriptor
	CurrentMethod    *descriptor.MethodDescriptor
	CurrentMessage   *descriptor.MessageDescriptor
	CurrentField     *descriptor.FieldDescriptor
	CurrentEnum      *descriptor.EnumDescriptor
	CurrentEnumValue *descriptor.EnumValueDescriptor
}

func (w *Walker) start(v Visitor) error {

	req := w.CurrentRequest

	fdMap, err := descriptor.CreateFileDescriptors(req.GetProtoFile())
	if err != nil {
		return xerrors.Errorf("on descriptor.CreateFileDescriptors: %w", err)
	}

	nonTargetFdMap := make(map[string]*descriptor.FileDescriptor)
	for k, v := range fdMap {
		nonTargetFdMap[k] = v
	}
	for _, fname := range req.GetFileToGenerate() {
		delete(nonTargetFdMap, fname)
	}
	for _, f := range req.GetProtoFile() {
		f, ok := nonTargetFdMap[f.GetName()]
		if !ok {
			// FileToGenerateじゃないものだけ先に処理する
			continue
		}

		err := w.visitFileDescriptor(v, f, &VisitFileInfo{IsGenerate: false})
		if err != nil {
			return err
		}
	}

	for _, fname := range req.GetFileToGenerate() {
		f := fdMap[fname]

		err := w.visitFileDescriptor(v, f, &VisitFileInfo{IsGenerate: true})
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *Walker) visitFileDescriptor(v Visitor, req *descriptor.FileDescriptor, info *VisitFileInfo) error {
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

	err := v.VisitFileDescriptor(w, req, optVal, info)
	if err != nil {
		return err
	}

	for _, srvc := range req.GetServices() {
		err := w.visitServiceDescriptor(v, srvc)
		if err != nil {
			return err
		}
	}

	for _, message := range req.GetMessageTypes() {
		err := w.visitMessageDescriptor(v, message, &VisitMessageInfo{})
		if err != nil {
			return err
		}
	}

	for _, enum := range req.GetEnumTypes() {
		err := w.visitEnumDescriptor(v, enum)
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
		err := w.visitMessageDescriptor(v, req.GetInputType(), &VisitMessageInfo{IsInput: true})
		if err != nil {
			return err
		}
	}
	{
		err := w.visitMessageDescriptor(v, req.GetOutputType(), &VisitMessageInfo{IsOutput: true})
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *Walker) visitMessageDescriptor(v Visitor, req *descriptor.MessageDescriptor, info *VisitMessageInfo) error {
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

	err := v.VisitMessageDescriptor(w, req, optVal, info)
	if err != nil {
		return err
	}

	for _, f := range req.GetFields() {
		err := w.visitFieldDescriptor(v, f)
		if err != nil {
			return err
		}
	}

	if !info.IsInput && !info.IsOutput {
		for _, m := range req.GetNestedMessageTypes() {
			// TODO nested らしくする
			err := w.visitMessageDescriptor(v, m, &VisitMessageInfo{})
			if err != nil {
				return err
			}
		}

		for _, e := range req.GetNestedEnumTypes() {
			// TODO nested らしくする
			err := w.visitEnumDescriptor(v, e)
			if err != nil {
				return err
			}
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

func (w *Walker) visitEnumDescriptor(v Visitor, req *descriptor.EnumDescriptor) error {
	bk := w.CurrentEnum
	w.CurrentEnum = req
	defer func() {
		w.CurrentEnum = bk
	}()

	var optVal *gqlgen_proto.EnumRule
	opts := req.GetOptions()
	if opts != nil && !isNilPtr(opts) {
		ext, err := proto.GetExtension(opts, gqlgen_proto.E_Enum)
		if xerrors.Is(err, proto.ErrMissingExtension) {
			// ok
		} else if err != nil {
			return xerrors.Errorf("%s on proto.GetExtension in visitEnumDescriptor: %w", req.GetFullyQualifiedName(), err)
		} else {
			optVal = ext.(*gqlgen_proto.EnumRule)
		}
	}

	err := v.VisitEnumDescriptor(w, req, optVal)
	if err != nil {
		return err
	}

	for _, ev := range req.GetValues() {
		err := w.visitEnumValueDescriptor(v, ev)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *Walker) visitEnumValueDescriptor(v Visitor, req *descriptor.EnumValueDescriptor) error {
	bk := w.CurrentEnumValue
	w.CurrentEnumValue = req
	defer func() {
		w.CurrentEnumValue = bk
	}()

	var optVal *gqlgen_proto.EnumValueRule
	opts := req.GetOptions()
	if opts != nil && !isNilPtr(opts) {
		ext, err := proto.GetExtension(opts, gqlgen_proto.E_EnumValue)
		if xerrors.Is(err, proto.ErrMissingExtension) {
			// ok
		} else if err != nil {
			return xerrors.Errorf("%s on proto.GetExtension in visitEnumValueDescriptor: %w", req.GetFullyQualifiedName(), err)
		} else {
			optVal = ext.(*gqlgen_proto.EnumValueRule)
		}
	}

	err := v.VisitEnumValueDescriptor(w, req, optVal)
	if err != nil {
		return err
	}

	return nil
}

func isNilPtr(x interface{}) bool {
	v := reflect.ValueOf(x)
	return v.Kind() == reflect.Ptr && v.IsNil()
}
