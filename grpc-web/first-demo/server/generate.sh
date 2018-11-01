#!/bin/bash -eu

cd `dirname $0`

protoc --proto_path=../ --go_out=plugins=grpc:./chat ../chat.proto
