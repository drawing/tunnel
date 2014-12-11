package engine

import (
	"encoding/gob"
	"net"
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

	EmpytConn
}

func NewTunConn(conn net.Conn, id uint64) *TunConn {
	dec := gob.NewDecoder(conn)
	enc := gob.NewEncoder(conn)
	return &TunConn{id, conn, dec, enc, true, EmpytConn{}}
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
		return ErrorEOF
	}

	err = conn.dec.Decode(pkg)
	if err != nil {
		conn.id = 0
	}
	return err
}

// ensure it 's atomic op
func (conn *TunConn) WritePackage(pkg *Package) (err error) {
	if !conn.available {
		return ErrorEOF
	}

	err = conn.enc.Encode(pkg)
	if err != nil {
		conn.id = 0
	}

	return err
}

// write normal data
func (conn *TunConn) Write(b []byte) (n int, err error) {
	if !conn.available {
		return 0, ErrorEOF
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
	defer func() { recover() }()

	if conn.id != 0 {
		pkg := Package{}
		pkg.Id = conn.id
		pkg.Command = PkgCommandClose

		conn.WritePackage(&pkg)
	} else {
		conn.c.Close()
	}

	conn.available = false
	return nil
}
