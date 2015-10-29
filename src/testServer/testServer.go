package main

import (
	"tcpServer"
	"flag"
	"strconv"
	"runtime"
	)
func main() {
	runtime.GOMAXPROCS(4)
	var serverIP = ""
	var serverPort = 0
	flag.Parse()	
	if flag.NArg() >= 2 {

		serverIP = flag.Arg(0)
		serverPort,_ = strconv.Atoi(flag.Arg(1))
	}
	server := new(tcpServer.Server)

	server.Run(serverIP,serverPort)	
}
