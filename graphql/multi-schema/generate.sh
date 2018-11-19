#!/bin/bash -eu

cd `dirname $0`

targets=`find . -type f \( -name '*.go' -and -not -iwholename '*vendor*'  -and -not -iwholename '*node_modules*' \)`
packages=`go list ./...`

set -x

# Apply tools
export PATH=$(pwd)/bin:$PATH
GO111MODULE=off go generate $packages
