rm ./bin/testServer
export GOPATH="/home/pi/tcpServer"
go install testServer
./bin/testServer 192.168.0.66 8008
