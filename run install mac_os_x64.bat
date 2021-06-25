set goos=darwin
set goarch=amd64

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