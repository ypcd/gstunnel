#!/bin/bash
export GO111MODULE="auto"
export GOPROXY="https://goproxy.io,direct"

sh ./run_test.sh
echo "test done."

sh ./run_install.sh
echo "install done."