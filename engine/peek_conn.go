package engine

import (
	"bytes"
	"net"
)

type PeekConn struct {
	reader net.Conn
	writer net.Conn

	buf *bytes.Buffer

	PeekMode bool

	EmpytConn
}

func NewPeekConn(reader net.Conn, writer net.Conn) net.Conn {
	return &PeekConn{reader, writer, bytes.NewBuffer(nil), true, EmpytConn{}}
}

func (conn *PeekConn) Clear() {
	conn.buf.Reset()
}

func (conn *PeekConn) Read(b []byte) (n int, err error) {
	if conn.PeekMode {
		n, err = conn.reader.Read(b)
		if err == nil {
			conn.buf.Write(b[:n])
		}
		return n, err
	} else {
		if conn.buf.Len() != 0 {
			return conn.buf.Read(b)
		}
		return conn.reader.Read(b)
	}
}

func (conn *PeekConn) Write(b []byte) (n int, err error) {
	return conn.writer.Write(b)
}

func (conn *PeekConn) Close() error {
	conn.reader.Close()
	conn.writer.Close()
	return nil
}
