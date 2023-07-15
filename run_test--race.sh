#!/bin/bash
export GO111MODULE="auto"
export GOPROXY="https://goproxy.io,direct"
export gorace="log_path=."

go mod tidy
go test -race ./...