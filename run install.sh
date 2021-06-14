#!/bin/bash

export GOPROXY="https://goproxy.io,direct"

cd gstunnel_server
go mod tidy
go install
cd ..

cd gstunnel_client
go mod tidy
go install
cd ..