#!/bin/sh
export GO111MODULE="auto"

sh ./run_test--race.sh
echo "test done."

sh ./run_install--race.sh
echo "install--race done."