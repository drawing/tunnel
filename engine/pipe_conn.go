package engine

import (
	"log"
	"net"
	"time"
)

type PipeConn struct {
	reader net.Conn
	writer net.Conn
}

func NewPipeConn(reader net.Conn, writer net.Conn) net.Conn {
	return &PipeConn{reader, writer}
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

// not implement
func (conn *PipeConn) LocalAddr() net.Addr {
	log.Fatalln("PipeConn LocalAddr not impement")
	return nil
}
func (conn *PipeConn) RemoteAddr() net.Addr {
	log.Fatalln("PipeConn RemoteAddr not impement")
	return nil
}
func (conn *PipeConn) SetDeadline(t time.Time) error {
	log.Fatalln("PipeConn SetDeadline not impement")
	return nil
}
func (conn *PipeConn) SetReadDeadline(t time.Time) error {
	log.Fatalln("PipeConn SetReadDeadline not impement")
	return nil
}
func (conn *PipeConn) SetWriteDeadline(t time.Time) error {
	log.Fatalln("PipeConn SetWriteDeadline not impement")
	return nil
}
