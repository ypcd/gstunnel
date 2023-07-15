SET GO111MODULE=auto

go mod tidy
go test -timeout 0 -count 1 -cover ./...
