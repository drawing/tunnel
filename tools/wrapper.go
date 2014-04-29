package tools

import (
	"bytes"
	"encoding/gob"
	"errors"
	"log"
	"net"
	"time"
)

const MaxBufferSize = 4092

const (
	PKG_KIND_CONNECT = iota
	PKG_KIND_DATA
	PKG_KIND_CLOSE
)

type Package struct {
	Kind int
	Id   int64
	Data []byte
	Addr *Location
}

type WrapperConn struct {
	id        int64
	c         net.Conn
	dec       *gob.Decoder
	available bool
}

func CreateWrapper(conn net.Conn, id int64) *WrapperConn {
	dec := gob.NewDecoder(conn)
	return &WrapperConn{id, conn, dec, true}
}

func (conn *WrapperConn) ReadPackage(pkg *Package) (err error) {
	if !conn.available {
		return errors.New("EOF")
	}
	return conn.dec.Decode(pkg)
}

// ensure it 's atomic op
func (conn *WrapperConn) WritePackage(pkg *Package) (err error) {
	if !conn.available {
		return errors.New("EOF")
	}
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err = enc.Encode(pkg)
	if err != nil {
		return err
	}
	_, err = conn.c.Write(network.Bytes())
	return err
}

// write normal data
func (conn *WrapperConn) Write(b []byte) (n int, err error) {
	if !conn.available {
		return 0, errors.New("EOF")
	}

	pkg := Package{}
	pkg.Id = conn.id
	pkg.Kind = 1
	pkg.Data = b

	err = conn.WritePackage(&pkg)
	if err != nil {
		return 0, err
	}

	return len(pkg.Data), err
}

func (conn *WrapperConn) Close() error {
	if conn.available {
		conn.Write([]byte{})
		conn.available = false
	}
	return nil
}

// not implement

// read normal data
func (conn *WrapperConn) Read(b []byte) (n int, err error) {
	log.Fatalln("WrapperConn LocalAddr not impement")
	return n, err
}
func (conn *WrapperConn) LocalAddr() net.Addr {
	log.Fatalln("WrapperConn LocalAddr not impement")
	return nil
}
func (conn *WrapperConn) RemoteAddr() net.Addr {
	log.Fatalln("WrapperConn RemoteAddr not impement")
	return nil
}
func (conn *WrapperConn) SetDeadline(t time.Time) error {
	log.Fatalln("WrapperConn SetDeadline not impement")
	return nil
}
func (conn *WrapperConn) SetReadDeadline(t time.Time) error {
	log.Fatalln("WrapperConn SetReadDeadline not impement")
	return nil
}
func (conn *WrapperConn) SetWriteDeadline(t time.Time) error {
	log.Fatalln("WrapperConn SetWriteDeadline not impement")
	return nil
}
