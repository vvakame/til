#!/bin/bash -eux

cd `dirname $0`

targets=`find . -type f \( -name '*.go' -and -not -iwholename '*vendor*'  -and -not -iwholename '*node_modules*' \)`
packages=`go list ./...`

# Apply tools
export PATH=$(pwd)/bin:$PATH
goimports -w $targets
go tool vet $targets
# golint -min_confidence 0.6 -set_exit_status $packages

go test -race ./... $@
