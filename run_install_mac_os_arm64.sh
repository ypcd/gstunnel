#!/bin/bash

export GOPROXY="https://goproxy.io,direct"

export GO111MODULE="on"
export GOOS="darwin"
export GOARCH="arm64"

cd gstunnellib
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