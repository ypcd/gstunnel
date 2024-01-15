#!/bin/bash
export GO111MODULE="auto"

go mod tidy
go get -u ./...
go test ./... -timeout 0 -p 1 -cover
