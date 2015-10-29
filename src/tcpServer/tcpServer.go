package tcpServer

import (
	"fmt"
	"net"
	"packet"
	"protocol"
	"encoding/json"	
	"compress/gzip"
	"bytes"
	"io/ioutil"
	"time"
)

type Server struct {
	ip string
	port int
	connectionCount int
	dispatcher protocol.ProtocolDispatcher
}

func (server *Server ) Run(serverIP string, serverPort int) int32 {

	println("server is running\n")
	server.ip =  serverIP;
	server.port = serverPort;
	server.dispatcher.RegisterProtocol()
	server.connectionCount = 0;
	go server.PrintStates()
	for {
		hostAndPort := fmt.Sprintf("%s:%d", serverIP, serverPort)
		listener := server.initServer(hostAndPort)
		for {
			conn, err := listener.Accept()
			server.connectionCount ++;
			//println("accpet")
			if err == nil {
				//println("accpet 2")
				go server.connectionHandler(conn)
			} else {
				break
			}
		}
		listener.Close()
	}
	return 0;
}

func (server *Server) PrintStates() {
	for {
		println("connectionCount:", server.connectionCount)
		time.Sleep(1000000000)
	}
}

func (server *Server) initServer(hostAndPort string) *net.TCPListener {
	serverAddr, err := net.ResolveTCPAddr("tcp", hostAndPort)
	server.checkError(err, "Resolving address:port failed: '"+hostAndPort+"'")
	listener, err := net.ListenTCP("tcp", serverAddr)
	server.checkError(err, "ListenTCP: ")
	println("Listening to: ", listener.Addr().String())
	return listener
}

func (server *Server )connectionHandler(conn net.Conn) {
	//connFrom := conn.RemoteAddr().String()
	//println("Connection from: ", connFrom)

	for {
		pkt := packet.Packet{}
		var sizeByte []byte = make([]byte, packet.HeadSize)

		length, err := conn.Read(sizeByte)
		if err != nil {
			goto DISCONNECT
		}

		pkt.Size = uint32(sizeByte[0])
		pkt.Size |= uint32(sizeByte[1]) << 8
		pkt.Size |= uint32(sizeByte[2]) << 16
		pkt.Size |= uint32(sizeByte[3]) << 24
		//println("recv:", pkt.Size)
		if pkt.Size > 1024*64 -1 {
			goto DISCONNECT
		}
		var flagsByte []byte = make([]byte, packet.FlagsSize)

		length, err = conn.Read(flagsByte)
		if err != nil {
			goto DISCONNECT
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
				goto DISCONNECT
			}
			recvBytes += uint32(length)
		}

		switch err {
		case nil:
			server.parse(conn, pkt, err)
		default:
			goto DISCONNECT
		}
	}
DISCONNECT:
	
	conn.Close()

	server.connectionCount--;
}

func (server *Server )parse(conn net.Conn, pkt packet.Packet, err error) {

	var f interface{}
	var json_err error
	if pkt.Flags & packet.FlagsGZip != 0 {
		gBuf := bytes.NewBuffer(pkt.Data)
		reader,_ := gzip.NewReader(gBuf)
		data,_ :=ioutil.ReadAll(reader)
		reader.Close()

		json_err = json.Unmarshal(data, &f)
		
	} else {
		json_err = json.Unmarshal(pkt.Data, &f)
	}
	if json_err != nil {
			println("msg not json: ", json_err.Error(), pkt.Size)
			println(string(pkt.Data))
		return
	}
	//println("msg is json: ", pkt.Size)
	//println(string(data))

	m := f.(map[string]interface{})
	if(m != nil) {
		server.dispatcher.Dispatch(conn, pkt, m)
	}

}


func (server *Server )checkError(error error, info string) {
	if error != nil {
		panic("ERROR: " + info + " " + error.Error()) // terminate program
	}
}



