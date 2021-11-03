#!/bin/sh
export GO111MODULE="auto"
export GOPROXY="https://goproxy.io,direct"

sh ./run_test--race.sh
echo "test done."

sh ./run_install--race.sh
echo "install--race done."