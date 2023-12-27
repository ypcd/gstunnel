#!/bin/bash
export GO111MODULE="auto"
#export GOPROXY="https://goproxy.io,direct"

go mod tidy
go get -u ./...
go test -timeout 0 -p 1 -cover ./...
