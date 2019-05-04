#!/bin/bash -eux

cd `dirname $0`

targets=`find . -type f \( -name '*.go' -and -not -iwholename '*vendor*'  -and -not -iwholename '*node_modules*' \)`
packages=`go list ./...`

export PATH=$(pwd)/build-cmd:$PATH
which goimports wire
go generate $packages
goimports -w $targets

go test $packages -p 1 -coverpkg=`go list -m`/... -covermode=atomic -coverprofile=coverage.txt $@
