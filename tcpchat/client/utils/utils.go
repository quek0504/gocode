package utils

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"gocode/chatroom/common/message"
	"net"
)

// utils package to send and read data to and from server
type Transfer struct {
	Conn   net.Conn
	Buffer [8096]byte
}

func (this *Transfer) ReadPkg() (mes message.Message, err error) {
	// waiting for data content length
	_, err = this.Conn.Read(this.Buffer[:4])
	if err != nil {
		//err = errors.New("read pkg header error")
		return
	}

	// fmt.Println("Read buffer=", this.Buffer[:4])

	// get data content length (convert buffer[:4] to uinit32)
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(this.Buffer[:4])

	// read actual data up to pkgLen
	n, err := this.Conn.Read(this.Buffer[:pkgLen])
	// incomplete data read
	if n != int(pkgLen) || err != nil {
		//err = errors.New("read pkg body error")
		return
	}

	// deserialize data from server(buffer[:pkgLen]) -> message.Message
	err = json.Unmarshal(this.Buffer[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarshal(buffer) err=", err)
		return
	}
	return
}

func (this *Transfer) WritePkg(data []byte) (err error) {
	// send data content length
	var pkgLen uint32          // 0 - 4294967295
	pkgLen = uint32(len(data)) // cast int to uint32
	// max 4 * 8 = 32 bytes
	binary.BigEndian.PutUint32(this.Buffer[0:4], pkgLen)

	n, err := this.Conn.Write(this.Buffer[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(buffer) fail", err)
		return
	}

	// send actual data
	n, err = this.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write(buffer) fail", err)
		return
	}
	return
}
