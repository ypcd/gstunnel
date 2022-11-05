SET GOPROXY="https://goproxy.io,direct"

SET GO111MODULE=on
SET GOOS=darwin
SET GOARCH=arm64

go install ./...