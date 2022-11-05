SET GO111MODULE=auto
SET GOPROXY="https://goproxy.io,direct"

go mod tidy
go test -timeout 0 -count 1 ./...