#!/bin/bash
export GO111MODULE="auto"
export GOPROXY="https://goproxy.io,direct"

export GO111MODULE="on"
export GOOS="darwin"
export GOARCH="arm64"

cd timerm
go test
go mod tidy
cd ..

cd gstunnellib
go test
go mod tidy
go install
cd ..

cd gstunnel_server
go mod tidy
go install
cd ..
echo "gstunnel_server installed."

cd gstunnel_client
go mod tidy
go install
cd ..
echo "gstunnel_client installed."