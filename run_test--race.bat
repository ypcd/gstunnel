SET GO111MODULE=auto
SET GOPROXY="https://goproxy.io,direct"

go mod tidy
go test -race -vet=off -timeout 0 -p 1 ./... 1>out.log 2>err.log