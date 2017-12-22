#!/bin/bash -eux

GOARCH=wasm GOOS=js go build -o test.wasm main.go
cp $GOROOT/misc/wasm/* ./
echo open http://localhost:8888/wasm_exec.html
go run server.go
