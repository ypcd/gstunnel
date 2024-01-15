#!/bin/bash
export GO111MODULE="auto"

go mod tidy
go install ./...
