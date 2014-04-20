package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"time"
)

const MaxBufferSize = 4092

type TransConn struct {
	c         net.Conn
	Id        uint16
	available bool
}

func CreateTrans(conn net.Conn, id uint16) *TransConn {
	return &TransConn{conn, id, true}
}

func (conn *TransConn) Read(b []byte) (n int, err error) {
	if !conn.available {
		return 0, errors.New("EOF")
	}

	if len(b) < MaxBufferSize {
		return 0, errors.New("byte size error")
	}

	_, err = io.ReadFull(conn.c, b[:4])
	if err != nil {
		return 0, err
	}

	var length uint16
	var id uint16

	binary.Read(bytes.NewBuffer(b[:2]), binary.BigEndian, &length)
	binary.Read(bytes.NewBuffer(b[2:4]), binary.BigEndian, &id)
	conn.Id = id

	if length == 0 {
		return 0, nil
	}
	if int(length) > len(b) {
		return 0, errors.New("recv length error")
	}

	_, err = io.ReadFull(conn.c, b[0:length])
	if err != nil {
		return 0, err
	}

	return int(length), nil
}

/*
func (conn *TransConn) Read(b []byte) (n int, err error) {
	if !conn.available {
		return 0, errors.New("EOF")
	}

	if len(b) < MaxBufferSize {
		return 0, errors.New("byte size error")
	}

	pkg := make([]byte, 4096*4)

	_, err = io.ReadFull(conn.c, pkg[:8])
	if err != nil {
		return 0, err
	}

	head, err := base64.StdEncoding.DecodeString(string(pkg[:8]))
	if err != nil {
		return 0, err
	}

	var length uint16
	var id uint16

	binary.Read(bytes.NewBuffer(head[:2]), binary.BigEndian, &length)
	binary.Read(bytes.NewBuffer(head[2:4]), binary.BigEndian, &id)
	conn.Id = id

	// log.Println("rr:", length, id)

	if length == 0 {
		return 0, nil
	}

	if int(length) > len(pkg) {
		log.Println("length:", length, b[0:4])
		return 0, errors.New("recv length error")
	}

	// log.Println("goo")
	_, err = io.ReadFull(conn.c, pkg[0:length])
	if err != nil {
		log.Println("1")
		return 0, err
	}
	// log.Println("dd")

	content, err := base64.StdEncoding.DecodeString(string(pkg[:length]))
	if err != nil {
		log.Println("2")
		return 0, err
	}

	if len(content) > len(b) {
		log.Println("too large", len(content), len(b))
		return 0, errors.New("too large")
	}

	for length = 0; length < uint16(len(content)); length++ {
		b[length] = content[length]
	}

	// log.Println("header:", id, length)
	return int(length), nil
}
*/

func (conn *TransConn) Write(b []byte) (n int, err error) {
	if !conn.available {
		return 0, errors.New("EOF")
	}

	pkg := new(bytes.Buffer)

	binary.Write(pkg, binary.BigEndian, uint16(len(b)))
	binary.Write(pkg, binary.BigEndian, conn.Id)
	binary.Write(pkg, binary.BigEndian, b)

	_, err = conn.c.Write(pkg.Bytes())

	return len(b), err
}

/*
func (conn *TransConn) Write(b []byte) (n int, err error) {
	if !conn.available {
		return 0, errors.New("EOF")
	}

	// log.Println("header:", conn.Id, len(b))

	pkg := new(bytes.Buffer)
	header := new(bytes.Buffer)

	data_pkg := []byte(base64.StdEncoding.EncodeToString(b))

	binary.Write(pkg, binary.BigEndian, b)
	binary.Write(header, binary.BigEndian, uint16(len(data_pkg)))
	binary.Write(header, binary.BigEndian, conn.Id)

	data_header := []byte(base64.StdEncoding.EncodeToString(header.Bytes()))

	// log.Println("header:", string(data_pkg))

	// log.Println("rr:", len(data_pkg), conn.Id)

	_, err = conn.c.Write(data_header)
	if err != nil {
		return 0, err
	}

	_, err = conn.c.Write(data_pkg)

	return len(b), err
}
*/

func (conn *TransConn) Close() error {
	if conn.available {
		conn.Write([]byte{})
		conn.available = false
	}
	return nil
}

// not implement
func (conn *TransConn) LocalAddr() net.Addr {
	log.Fatalln("TransConn LocalAddr not impement")
	return nil
}
func (conn *TransConn) RemoteAddr() net.Addr {
	log.Fatalln("TransConn RemoteAddr not impement")
	return nil
}
func (conn *TransConn) SetDeadline(t time.Time) error {
	log.Fatalln("TransConn SetDeadline not impement")
	return nil
}
func (conn *TransConn) SetReadDeadline(t time.Time) error {
	log.Fatalln("TransConn SetReadDeadline not impement")
	return nil
}
func (conn *TransConn) SetWriteDeadline(t time.Time) error {
	log.Fatalln("TransConn SetWriteDeadline not impement")
	return nil
}
