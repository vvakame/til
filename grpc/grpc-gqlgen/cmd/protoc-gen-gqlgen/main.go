package main

import (
	"bytes"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	proto_extentions "github.com/vvakame/til/grpc/grpc-gqlgen/proto-extentions"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
)

func main() {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	_ = ioutil.WriteFile("./protoc-gen-gqlgen.input.dump", b, 0666)

	err = run(bytes.NewBuffer(b), os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}

func run(r io.Reader, w io.Writer) error {
	req, err := parseReq(r)
	if err != nil {
		return err
	}

	resp := processReq(req)

	return emitResp(w, resp)
}

func parseReq(r io.Reader) (*plugin.CodeGeneratorRequest, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var req plugin.CodeGeneratorRequest
	if err = proto.Unmarshal(buf, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

func processReq(req *plugin.CodeGeneratorRequest) *plugin.CodeGeneratorResponse {
	files := make(map[string]*descriptor.FileDescriptorProto)
	for _, f := range req.ProtoFile {
		files[f.GetName()] = f
	}

	resp := &plugin.CodeGeneratorResponse{}

	for _, fname := range req.FileToGenerate {
		f := files[fname]

		var buf bytes.Buffer
		walkFileDescriptor(&buf, f)

		genFileName := fmt.Sprintf("%s.info", path.Base(fname))
		resp.File = append(resp.File, &plugin.CodeGeneratorResponse_File{
			Name:    proto.String(genFileName),
			Content: proto.String(buf.String()),
		})
	}

	return resp
}

func emitResp(w io.Writer, resp *plugin.CodeGeneratorResponse) error {
	buf, err := proto.Marshal(resp)
	if err != nil {
		return err
	}
	_, err = w.Write(buf)
	return err
}

func walkFileDescriptor(w io.Writer, f *descriptor.FileDescriptorProto) {
	_, _ = fmt.Fprintf(w, "file: %s\n", f.GetName())

	opts := f.GetOptions()
	if opts != nil {
		ext, err := proto.GetExtension(opts, proto_extentions.E_Resolver)
		if err == proto.ErrMissingExtension {
			// ok
		} else if err != nil {
			log.Fatal(err)
		} else {
			opt := ext.(*proto_extentions.FileRule)
			for _, v := range opt.GetTypeInference() {
				_, _ = fmt.Fprintf(w, "fileRule: %s %s %s\n", v.GetSrc(), v.GetDest(), v.GetType())
			}
		}
	}

	for _, srv := range f.GetService() {
		_, _ = fmt.Fprintf(w, "service: %s\n", srv.GetName())

		for _, mt := range srv.GetMethod() {
			_, _ = fmt.Fprintf(w, "method: %s\n", mt.GetName())

			opts := mt.GetOptions()
			if opts != nil {
				ext, err := proto.GetExtension(opts, proto_extentions.E_Schema)
				if err == proto.ErrMissingExtension {
					// ok
				} else if err != nil {
					log.Fatal(err)
				} else {
					v := ext.(*proto_extentions.SchemaRule)
					_, _ = fmt.Fprintf(w, "schemaRule: %s %s %s %s\n", v.GetPattern(), v.GetQuery(), v.GetMutation(), v.GetSubscription())
				}
			}
		}
	}

	for _, msg := range f.GetMessageType() {
		_, _ = fmt.Fprintf(w, "message: %s\n", msg.GetName())

		opts := msg.GetOptions()
		if opts != nil {
			ext, err := proto.GetExtension(opts, proto_extentions.E_Type)
			if err == proto.ErrMissingExtension {
				// ok
			} else if err != nil {
				log.Fatal(err)
			} else {
				v := ext.(*proto_extentions.MessageRule)
				_, _ = fmt.Fprintf(w, "messageRule: %s %s\n", v.GetType(), v.GetAlias())
			}
		}

		for _, f := range msg.GetField() {
			_, _ = fmt.Fprintf(w, "field: %s\n", f.GetName())

			opts := f.GetOptions()
			if opts != nil {
				ext, err := proto.GetExtension(opts, proto_extentions.E_Field)
				if err == proto.ErrMissingExtension {
					// ok
				} else if err != nil {
					log.Fatal(err)
				} else {
					v := ext.(*proto_extentions.FieldRule)
					_, _ = fmt.Fprintf(w, "fieldRule: %s\n", v.GetAlias())
				}
			}
		}
	}
}
