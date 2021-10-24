#!/bin/bash
export GO111MODULE="auto"
export GOPROXY="https://goproxy.io,direct"

cd gstunnellib
go mod tidy
go test
cd ..

cd gstunnel_server
go mod tidy
go install
cd ..

cd gstunnel_client
go mod tidy
go install
cd ..