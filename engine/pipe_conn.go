package engine

import (
	"net"
)

type PipeConn struct {
	reader net.Conn
	writer net.Conn

	EmpytConn
}

func NewPipeConn(reader net.Conn, writer net.Conn) net.Conn {
	return &PipeConn{reader, writer, EmpytConn{}}
}

func (conn *PipeConn) Read(b []byte) (n int, err error) {
	return conn.reader.Read(b)
}

func (conn *PipeConn) Write(b []byte) (n int, err error) {
	return conn.writer.Write(b)
}

func (conn *PipeConn) Close() error {
	conn.reader.Close()
	conn.writer.Close()
	return nil
}
