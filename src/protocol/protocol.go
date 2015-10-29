package protocol

import (
	//"fmt"
	"net"
	"packet"
	//"crypto/md5"
	//"encoding/hex"
	//"encoding/json"	
	//"strconv"
	//"strings"
	//"db"
	//"time"
	//"github.com/anachronistic/apns"
	//"math/rand"
	)

type Protocol interface {
	Process(conn net.Conn, pkt packet.Packet, data map[string]interface{}) bool;
}

type ProtocolEcho struct {

}
func( protocol *ProtocolEcho) Process(conn net.Conn, pkt packet.Packet, data map[string]interface{}) bool {
	return packet.SendMapData(conn, data)
}

type ProtocolDispatcher struct {
	handlerMap map[string] Protocol ;
}
	
func (dispatcher *ProtocolDispatcher) RegisterProtocol () {
	dispatcher.handlerMap = make(map[string] Protocol)
	
	dispatcher.handlerMap["echo"] = new(ProtocolEcho)
}

func (dispatcher *ProtocolDispatcher)Dispatch(conn net.Conn, pkt packet.Packet, data map[string]interface{} ) {


	msgType := data["type"].(string)
	//println("Dispatch type: ", msgType)
	var handler Protocol = dispatcher.handlerMap[msgType]
	if handler != nil {
		handler.Process(conn, pkt ,data)
	}
}
