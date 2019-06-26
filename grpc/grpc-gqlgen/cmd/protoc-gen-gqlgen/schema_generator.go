package main

import (
	"bytes"
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	basedescriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/vektah/gqlparser/ast"
	"github.com/vektah/gqlparser/formatter"
	gqlgen_proto "github.com/vvakame/til/grpc/grpc-gqlgen/gqlgen-proto"
)

func schemaGenerate(ctx context.Context, b *Builder) ([]*plugin.CodeGeneratorResponse_File, error) {
	g := &schemaGenerator{b: b}
	return g.Generate(ctx, b.FileInfos)
}

type schemaGenerator struct {
	b *Builder
}

func (g *schemaGenerator) Generate(ctx context.Context, fileInfos []*FileInfo) ([]*plugin.CodeGeneratorResponse_File, error) {
	var files []*plugin.CodeGeneratorResponse_File
	for _, fileInfo := range fileInfos {
		doc := &ast.SchemaDocument{}

		for _, service := range fileInfo.Services {
			for _, method := range service.Methods {
				def := &ast.Definition{
					Kind: ast.Object,
				}

				switch method.GraphQLOperationType {
				case GraphQLQuery:
					def.Name = "Query" // TODO カスタム名対応

				case GraphQLMutation:
					def.Name = "Mutation" // TODO カスタム名対応

				case GraphQLSubscription:
					def.Name = "Subscription" // TODO カスタム名対応
				}

				field := &ast.FieldDefinition{
					Name: method.GraphQLName,
					Arguments: ast.ArgumentDefinitionList{ // Request
						{
							Name: "input",
							Type: g.messageInfoToType(method.RequestMessage),
						},
					},
					Type: g.messageInfoToType(method.ResponseMessage),
				}
				def.Fields = ast.FieldList{
					field,
				}

				doc.Extensions = append(doc.Extensions, def)
			}
		}

		for _, message := range fileInfo.MessageInfos {

			def := &ast.Definition{
				Name: message.GraphQLName(),
			}
			switch message.GraphQLMessageType {
			case gqlgen_proto.MessageType_UNKNOWN,
				gqlgen_proto.MessageType_TYPE:
				def.Kind = ast.Object

			case gqlgen_proto.MessageType_INPUT:
				def.Kind = ast.InputObject
			}

			for _, fieldInfo := range message.Fields {
				field := &ast.FieldDefinition{
					Name:      fieldInfo.GraphQLName(),
					Arguments: nil,
					Type:      g.fieldInfoToType(fieldInfo),
				}
				def.Fields = append(def.Fields, field)
			}

			doc.Definitions = append(doc.Definitions, def)
		}

		for _, enum := range fileInfo.EnumInfos {

			def := &ast.Definition{
				Name: enum.GraphQLName(),
				Kind: ast.Enum,
			}

			for _, valueInfo := range enum.Values {
				value := &ast.EnumValueDefinition{
					Name: valueInfo.GraphQLName(),
				}
				def.EnumValues = append(def.EnumValues, value)
			}

			doc.Definitions = append(doc.Definitions, def)
		}

		var buf bytes.Buffer
		err := formatter.NewFormatter(&buf).FormatSchemaDocument(doc)
		if err != nil {
			return nil, err
		}

		files = append(files, &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(fmt.Sprintf("%s.graphql", fileInfo.PackageName)),
			Content: proto.String(buf.String()),
		})
	}

	return files, nil
}

func (g *schemaGenerator) messageInfoToType(message *MessageInfo) *ast.Type {
	return &ast.Type{
		NamedType: message.GraphQLName(),
		NonNull:   true,
	}
}

func (g *schemaGenerator) fieldInfoToType(field *FieldInfo) *ast.Type {
	t := &ast.Type{
		NonNull: !field.GraphQLOptional,
	}

outer:
	switch field.Type {
	case basedescriptor.FieldDescriptorProto_TYPE_STRING:
		t.NamedType = "String"

	case basedescriptor.FieldDescriptorProto_TYPE_BOOL:
		t.NamedType = "Boolean"

	case basedescriptor.FieldDescriptorProto_TYPE_INT32,
		basedescriptor.FieldDescriptorProto_TYPE_UINT32:
		t.NamedType = "Int"

	case basedescriptor.FieldDescriptorProto_TYPE_FLOAT,
		basedescriptor.FieldDescriptorProto_TYPE_FIXED32:
		t.NamedType = "Float"

	case basedescriptor.FieldDescriptorProto_TYPE_ENUM:
		fqn := field.TypeEnum.GetFullyQualifiedName()
		switch fqn {
		case "google.protobuf.Timestamp":
			t.NamedType = "Timestamp"
			break outer
		}

		enumInfo := g.b.FindEnumInfo(fqn)
		if enumInfo == nil {
			panic(fmt.Sprintf("specified EnumInfo doesn't exists: %s", field.TypeEnum.GetFullyQualifiedName()))
		}
		t.NamedType = enumInfo.GraphQLName()

	case basedescriptor.FieldDescriptorProto_TYPE_MESSAGE:
		fqn := field.TypeMessage.GetFullyQualifiedName()
		switch fqn {
		case "google.protobuf.Timestamp":
			t.NamedType = "Timestamp"
			break outer
		}

		messageInfo := g.b.FindMessageInfo(fqn)
		if messageInfo == nil {
			panic(fmt.Sprintf("specified MessageInfo doesn't exists: %s", field.TypeMessage.GetFullyQualifiedName()))
		}
		t.NamedType = messageInfo.GraphQLName()

	default:
		panic(fmt.Sprintf("unknown type: %s", field.Type.String()))
	}

	if field.Repeated {
		t = &ast.Type{
			Elem:    t,
			NonNull: true,
		}
	}

	return t
}
