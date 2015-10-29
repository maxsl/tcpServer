rm ./bin/testServer
export GOPATH="/Users/lavrock/code/golang/tcpServer"
go install testServer
./bin/testServer 127.0.0.1 8008
