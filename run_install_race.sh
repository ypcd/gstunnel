#!/bin/sh
export GO111MODULE="auto"

go mod tidy
go get -u ./...
#go install -race ./...

cd gstunnel_server
sh run_build_race.sh
mv gstunnel_server_race $HOME/go/bin
cd ..

cd gstunnel_client
sh run_build_race.sh
mv gstunnel_client_race $HOME/go/bin
cd ..
