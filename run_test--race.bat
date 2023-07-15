SET GO111MODULE=auto
set gorace=log_path=.

go mod tidy
go test -race -vet=off -timeout 0 -p 1 -cover ./...