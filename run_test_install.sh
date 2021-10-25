#!/bin/bash
export GO111MODULE="auto"
export GOPROXY="https://goproxy.io,direct"

source run_test.sh

source run_install.sh