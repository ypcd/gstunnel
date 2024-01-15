#!/bin/bash
export GO111MODULE="auto"

sh ./run_test.sh
echo "test done."

sh ./run_install.sh
echo "install done."