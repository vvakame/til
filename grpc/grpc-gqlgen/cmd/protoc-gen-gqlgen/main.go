//go:generate statik -src ./tmpls

package main

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/protobuf/proto"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	_ "github.com/golang/protobuf/ptypes/timestamp"
	_ "github.com/vvakame/til/grpc/grpc-gqlgen/cmd/protoc-gen-gqlgen/statik"
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

	ctx := context.Background()

	bldr := &Builder{}
	resp, err := bldr.Process(ctx, req)
	if err != nil {
		return err
	}

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

func emitResp(w io.Writer, resp *plugin.CodeGeneratorResponse) error {
	buf, err := proto.Marshal(resp)
	if err != nil {
		return err
	}
	_, err = w.Write(buf)
	return err
}
