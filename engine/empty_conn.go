package engine

import (
	"errors"
	"net"
	"time"
)

var ErrorNotImpl error = errors.New("method not implement")
var ErrorEOF error = errors.New("EOF")

type EmpytConn struct{}

func (conn *EmpytConn) Read(b []byte) (n int, err error) {
	return 0, ErrorNotImpl
}

func (conn *EmpytConn) Write(b []byte) (n int, err error) {
	return 0, ErrorNotImpl
}

func (conn *EmpytConn) Close() error {
	return ErrorNotImpl
}

func (conn *EmpytConn) LocalAddr() net.Addr {
	return nil
}

func (conn *EmpytConn) RemoteAddr() net.Addr {
	return nil
}

func (conn *EmpytConn) SetDeadline(t time.Time) error {
	return ErrorNotImpl
}

func (conn *EmpytConn) SetReadDeadline(t time.Time) error {
	return ErrorNotImpl
}

func (conn *EmpytConn) SetWriteDeadline(t time.Time) error {
	return ErrorNotImpl
}
