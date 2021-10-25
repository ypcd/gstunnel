#!/bin/bash
export GO111MODULE="auto"
export GOPROXY="https://goproxy.io,direct"

export GO111MODULE="on"
export GOOS="darwin"
export GOARCH="arm64"

source run_test.sh
source run_install.sh