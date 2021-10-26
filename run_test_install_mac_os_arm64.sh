#!/bin/bash
export GO111MODULE="auto"
export GOPROXY="https://goproxy.io,direct"

export GO111MODULE="on"
export GOOS="darwin"
export GOARCH="arm64"

sh ./run_test.sh
sh ./run_install.sh