#!/bin/sh -eux

cd `dirname $0`

go get -u golang.org/x/tools/cmd/goimports
# go get -u golang.org/x/tools/cmd/vet

set +e # ubuntu uses too old go version...
go get -u github.com/golang/lint/golint
set -e

go get -u github.com/favclip/jwg/cmd/jwg
go get -u github.com/favclip/qbg/cmd/qbg

go get -u github.com/constabulary/gb/...
go get -u github.com/PalmStoneGames/gb-gae

gb vendor restore
