package engine

import (
	"encoding/gob"
	"errors"
	"log"
	"net"
	"time"
)

const MaxBufferSize = 4092

const (
	PkgCommandConnect = iota
	PkgCommandData
	PkgCommandClose
	PkgCommandRegister
)

type Package struct {
	Command int
	Id      uint64
	Data    []byte
	Loc     *Location
	Router  *RouterItem
}

type TunConn struct {
	id        uint64
	c         net.Conn
	dec       *gob.Decoder
	enc       *gob.Encoder
	available bool
}

func NewTunConn(conn net.Conn, id uint64) *TunConn {
	dec := gob.NewDecoder(conn)
	enc := gob.NewEncoder(conn)
	return &TunConn{id, conn, dec, enc, true}
}

func (conn *TunConn) SetID(id uint64) {
	conn.id = id
}

func (conn *TunConn) Clone() *TunConn {
	n := new(TunConn)
	n.id = conn.id
	n.c = conn.c
	n.dec = conn.dec
	n.enc = conn.enc
	n.available = conn.available
	return n
}

func (conn *TunConn) ReadPackage(pkg *Package) (err error) {
	if !conn.available {
		return errors.New("EOF")
	}

	return conn.dec.Decode(pkg)
}

// ensure it 's atomic op
func (conn *TunConn) WritePackage(pkg *Package) (err error) {
	if !conn.available {
		return errors.New("EOF")
	}
	return conn.enc.Encode(pkg)
}

// write normal data
func (conn *TunConn) Write(b []byte) (n int, err error) {
	if !conn.available {
		return 0, errors.New("EOF")
	}

	pkg := Package{}
	pkg.Id = conn.id
	pkg.Command = PkgCommandData
	pkg.Data = b

	err = conn.WritePackage(&pkg)
	if err != nil {
		return 0, err
	}

	return len(pkg.Data), err
}

func (conn *TunConn) Close() error {
	if conn.available {
		conn.Write([]byte{})
		conn.available = false
	}
	return nil
}

// not implement

// read normal data
func (conn *TunConn) Read(b []byte) (n int, err error) {
	log.Fatalln("TunConn LocalAddr not impement")
	return n, err
}
func (conn *TunConn) LocalAddr() net.Addr {
	log.Fatalln("TunConn LocalAddr not impement")
	return nil
}
func (conn *TunConn) RemoteAddr() net.Addr {
	log.Fatalln("TunConn RemoteAddr not impement")
	return nil
}
func (conn *TunConn) SetDeadline(t time.Time) error {
	log.Fatalln("TunConn SetDeadline not impement")
	return nil
}
func (conn *TunConn) SetReadDeadline(t time.Time) error {
	log.Fatalln("TunConn SetReadDeadline not impement")
	return nil
}
func (conn *TunConn) SetWriteDeadline(t time.Time) error {
	log.Fatalln("TunConn SetWriteDeadline not impement")
	return nil
}
