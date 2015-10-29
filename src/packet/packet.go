package packet

import (
	"fmt"
	"net"
	"bytes"
	"compress/gzip"
	//"crypto/md5"
	"encoding/binary"
	//"encoding/hex"
	"encoding/json"	
	//"strconv"
	//"strings"
	//"db"
	//"time"
	//"github.com/anachronistic/apns"
	//"math/rand"
	)

const (
	GZipEnable = false
	)

const (
	FlagsGZip = 1
)
const (
	HeadSize = 4
	FlagsSize = 4
)
type Packet struct {
	Size uint32
	Flags uint32 
	Data []byte
}

func SendMapData(conn net.Conn, resp map[string]interface{}) bool {
	if conn == nil {
		return false
	}
	var length, flags uint32

	buf := new(bytes.Buffer)
	jsonData, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("json error", err.Error())

		return false
	}
	if(GZipEnable) {
		var gBuf bytes.Buffer
		w := gzip.NewWriter(&gBuf)
		w.Write(jsonData)
		//w.Flush()
		w.Close()
		length = uint32(len(gBuf.Bytes()) + 8)
		flags |= FlagsGZip
		binary.Write(buf, binary.LittleEndian, length)
		binary.Write(buf, binary.LittleEndian, flags)
	
		_, err = conn.Write(buf.Bytes())
		if err != nil {
			goto DISCONNECT
		}
		_, err = conn.Write(gBuf.Bytes())
		if err != nil {
			goto DISCONNECT
		}
	} else {

		length = uint32(len(jsonData) + 8)
		flags = 0
		binary.Write(buf, binary.LittleEndian, length)
		binary.Write(buf, binary.LittleEndian, flags)

		_, err = conn.Write(buf.Bytes())
		if err != nil {
			goto DISCONNECT
		}
		_, err = conn.Write(jsonData)
		if err != nil {
			goto DISCONNECT
		}
	}


	
	return true
DISCONNECT:

	conn.Close()
	return false
}