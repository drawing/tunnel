package engine

import (
	"errors"
	"net"
	"time"
)

type EmpytConn struct{}

var NotImplErr error = errors.New("method not implement")

func (conn *EmpytConn) Read(b []byte) (n int, err error) {
	return 0, NotImplErr
}

func (conn *EmpytConn) Write(b []byte) (n int, err error) {
	return 0, NotImplErr
}

func (conn *EmpytConn) Close() error {
	return NotImplErr
}

func (conn *EmpytConn) LocalAddr() net.Addr {
	return nil
}

func (conn *EmpytConn) RemoteAddr() net.Addr {
	return nil
}

func (conn *EmpytConn) SetDeadline(t time.Time) error {
	return NotImplErr
}

func (conn *EmpytConn) SetReadDeadline(t time.Time) error {
	return NotImplErr
}

func (conn *EmpytConn) SetWriteDeadline(t time.Time) error {
	return NotImplErr
}
