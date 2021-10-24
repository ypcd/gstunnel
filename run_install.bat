set GO111MODULE="auto"
set GOPROXY="https://goproxy.io,direct"

cd gstunnellib
go mod tidy
go test
cd ..

cd gstunnel_client
go mod tidy
go install
cd ..

cd gstunnel_server
go mod tidy
go install
cd ..

timeout 10000