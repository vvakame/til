package main

import (
	"bytes"
	"context"
	"io/ioutil"
	"testing"
)

func Test(t *testing.T) {
	b, err := ioutil.ReadFile("../../protoc-gen-gqlgen.input.dump")
	if err != nil {
		t.Fatal(err)
	}

	req, err := parseReq(bytes.NewBuffer(b))
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	bldr := &Builder{}
	resp, err := bldr.Process(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range resp.GetFile() {
		t.Log(*f.Name)
		t.Log(*f.Content)
	}
}
