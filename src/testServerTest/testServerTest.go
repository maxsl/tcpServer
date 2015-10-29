package main

import (
	"net"
	"flag"
	"strconv"
	"fmt"
	"packet"
	"runtime"
	)
var testPacketCount int

var packetCount int
var connectionCount int

func main() {
	runtime.GOMAXPROCS(16)
	var serverIP = ""
	var serverPort = 0
	testPacketCount = 1000
	packetCount = 0
	flag.Parse()	
	println("flag.NArg():",flag.NArg())
	if flag.NArg() >= 2 {

		serverIP = flag.Arg(0)
		serverPort,_ = strconv.Atoi(flag.Arg(1))

	}
	if flag.NArg() >= 3  {
		connectionCount,_ = strconv.Atoi(flag.Arg(2))
		testPacketCount,_ = strconv.Atoi(flag.Arg(3))	
	}
	hostAndPort := fmt.Sprintf("%s:%d", serverIP, serverPort)
		
	for i:=1; i < connectionCount; i++ {
		conn,err := net.Dial("tcp",hostAndPort)
		if(err == nil) {
				go connectionHandler(conn)
		} else {
			println("connect failed")
		
		}	
	}
	conn,err := net.Dial("tcp",hostAndPort)
	if(err == nil) {
			connectionHandler(conn)
	} else {
		println("connect failed")
	
	}
}

func recvPacket(conn net.Conn) ( packet.Packet, bool) {
	pkt := packet.Packet{}
	var sizeByte []byte = make([]byte, packet.HeadSize)

	length, err := conn.Read(sizeByte)
	if err != nil {
		println("connection disconnected!")
		return pkt,false
	}

	pkt.Size = uint32(sizeByte[0])
	pkt.Size |= uint32(sizeByte[1]) << 8
	pkt.Size |= uint32(sizeByte[2]) << 16
	pkt.Size |= uint32(sizeByte[3]) << 24
	if pkt.Size > 1024*64 -1 {
		return pkt,false
	}
	var flagsByte []byte = make([]byte, packet.FlagsSize)

	length, err = conn.Read(flagsByte)
	if err != nil {
		println("connection disconnected!")
		return pkt,false
	}
	pkt.Flags = uint32(flagsByte[0])
	pkt.Flags |= uint32(flagsByte[1]) << 8
	pkt.Flags |= uint32(flagsByte[2]) << 16
	pkt.Flags |= uint32(flagsByte[3]) << 24
	var contentSize uint32
	contentSize = pkt.Size - uint32(packet.HeadSize) - uint32(packet.FlagsSize)
	pkt.Data = make([]byte, int(contentSize))
	var recvBytes uint32
	recvBytes = 0
	for uint32(recvBytes) < contentSize {
		length, err = conn.Read(pkt.Data[recvBytes:])
		if err != nil {
			return pkt,false
		}
		recvBytes += uint32(length)
	}
	return pkt,true
}

func sendEchoPacket(conn net.Conn) bool{
	m := make(map[string] interface {})
	m["type"] = "echo"
	m["msg"] = "hello world! hello world! hello world! hello world! hello world! hello world! hello world!hello world! hello world! hello world! hello world! hello world! hello world! hello world!hello world! hello world! hello world! hello world! hello world! hello world! hello world!hello world! hello world! hello world! hello world! hello world! hello world! hello world!hello world! hello world! hello world! hello world! hello world! hello world! hello world!hello world! hello world! hello world! hello world! hello world! hello world! hello world!hello world! hello world! hello world! hello world! hello world! hello world! hello world!hello world! hello world! hello world! hello world! hello world! hello world! hello world!hello world! hello world! hello world! hello world! hello world! hello world! hello world!hello world! hello world! hello world! hello world! hello world! hello world! hello world!hello world! hello world! hello world! hello world! hello world! hello world! hello world!hello world! hello world! hello world! hello world! hello world! hello world! hello world!hello world! hello world! hello world! hello world! hello world! hello world! hello world!hello world! hello world! hello world! hello world! hello world! hello world! hello world!hello world! hello world! hello world! hello world! hello world! hello world! hello world!hello world! hello world! hello world! hello world! hello world! hello world! hello world!hello world! hello world! hello world! hello world! hello world! hello world! hello world!hello world! hello world! hello world! hello world! hello world! hello world! hello world!hello world! hello world! hello world! hello world! hello world! hello world! hello world!hello world! hello world! hello world! hello world! hello world! hello world! hello world!hello world! hello world! hello world! hello world! hello world! hello world! hello world!hello world! hello world! hello world! hello world! hello world! hello world! hello world!hello world! hello world! hello world! hello world! hello world! hello world! hello world!hello world! hello world! hello world! hello world! hello world! hello world! hello world!hello world! hello world! hello world! hello world! hello world! hello world! hello world!hello world! hello world! hello world! hello world! hello world! hello world! hello world!hello world! hello world! hello world! hello world! hello world! hello world! hello world!hello world! hello world! hello world! hello world! hello world! hello world! hello world!hello world! hello world! hello world! hello world! hello world! hello world! hello world!hello world! hello world! hello world! hello world! hello world! hello world! hello world! hello world! hello world! hello world! hello world! hello world! hello world! hello world! hello world! hello world! hello world! hello world! hello world! hello world! hello world! hello world! hello world! hello world! hello world! hello world! hello world! hello world! hello world! hello world! "
	return packet.SendMapData(conn,m) 
}

func connectionHandler(conn net.Conn) {
	connFrom := conn.RemoteAddr().String()
	println("Connection to: ", connFrom)
	for {

		if(!sendEchoPacket(conn)){
			goto DISCONNECT
		}
		_,ok :=recvPacket(conn)
		if ( !ok ) {
			goto DISCONNECT
		}
		packetCount++;
		if(packetCount % 10 == 0) {
			println("packet count" ,packetCount)
		}
		if(packetCount > testPacketCount) {
			break
		}
	}	
	println("packet count" ,packetCount)
		
DISCONNECT:
	
	conn.Close()

	println("Closed connection:", connFrom)
}
	
