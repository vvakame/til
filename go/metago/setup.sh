#!/bin/bash -eux

cd "$(dirname "$0")"

# build tools
rm -rf build-cmd/
mkdir build-cmd

GOBIN=$(pwd -P)/build-cmd
export GOBIN
go install golang.org/x/tools/cmd/goimports
go install golang.org/x/lint/golint
go install honnef.co/go/tools/cmd/staticcheck
