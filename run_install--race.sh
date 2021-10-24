#!/bin/sh
export GO111MODULE="auto"
export GOPROXY="https://goproxy.io,direct"

cd gstunnellib
go mod tidy
cd ..

cd gstunnel_client
go mod tidy
go install -race
cd ..

cd gstunnel_server
go mod tidy
go install -race
cd ..
