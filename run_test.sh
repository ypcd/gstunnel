#!/bin/bash
export GO111MODULE="auto"
export GOPROXY="https://goproxy.io,direct"

cd gstunnellib
go mod tidy
go test
cd ..
