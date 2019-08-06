#!/bin/bash -eux

cd "$(dirname "$0")"

targets=$(find . -type f \( -name '*.go' -and -not -iwholename '*vendor*' -and -not -iwholename '*node_modules*' \))
packages=$(go list ./...)

PATH=$(pwd)/build-cmd:$PATH
export PATH
command -v goimports golint staticcheck
# shellcheck disable=SC2086
go generate $packages
# shellcheck disable=SC2086
goimports -w $targets
# go vet ./...
# shellcheck disable=SC2086
# golint -min_confidence 0.6 -set_exit_status $packages
# shellcheck disable=SC2086
# staticcheck $packages

# shellcheck disable=SC2086
go test $packages -p 1 -coverpkg="$(go list -m)/..." -covermode=atomic -coverprofile=coverage.txt "$@"
go tool cover -html=./coverage.txt -o cover.html

go build ./internal/testbed/...
go build -tags metago ./internal/testbed/...
