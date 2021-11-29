#!/bin/sh
export GO111MODULE="auto"
export GOPROXY="https://goproxy.io,direct"

go mod tidy
go install -race ./...
