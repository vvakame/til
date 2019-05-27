package main

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func Test(t *testing.T) {
	b, err := ioutil.ReadFile("../../protoc-gen-gqlgen.input.dump")
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	err = run(bytes.NewBuffer(b), &buf)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(buf.String())
}
