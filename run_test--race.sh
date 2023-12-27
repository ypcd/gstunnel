#!/bin/bash
export GO111MODULE="auto"
#export GOPROXY="https://goproxy.io,direct"
export gorace="log_path=."

go mod tidy
go get -u ./...
go test -race -timeout 0 -p 1 -cover ./...