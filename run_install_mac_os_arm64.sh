#!/bin/bash

export GOPROXY="https://goproxy.io,direct"

export GO111MODULE="on"
export GOOS="darwin"
export GOARCH="arm64"

source run_install.sh