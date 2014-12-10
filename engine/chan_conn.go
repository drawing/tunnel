package engine

import (
	"errors"
	"log"
	"net"
)

type ChannelConn struct {
	c         chan []byte
	available bool
	data      []byte

	EmpytConn
}

func NewChannelConn() net.Conn {
	return &ChannelConn{make(chan []byte, 30), true, []byte{}, EmpytConn{}}
}

func (conn *ChannelConn) Read(b []byte) (n int, err error) {
	if !conn.available {
		return 0, errors.New("EOF")
	}

	if len(b) == 0 {
		return 0, errors.New("zero []byte")
	}

	if len(conn.data) == 0 {
		var ok bool
		conn.data, ok = <-conn.c
		if !ok {
			return 0, errors.New("EOF")
		}
	}

	length := 0
	for length < len(b) && length < len(conn.data) {
		b[length] = conn.data[length]
		length++
	}

	if length < len(conn.data) {
		conn.data = conn.data[length:]
	} else {
		conn.data = []byte{}
	}

	return length, nil
}

func (conn *ChannelConn) Write(b []byte) (n int, err error) {
	defer func() {
		msg := recover()
		if msg != nil {
			log.Println("Channel Write", msg)
			n = 0
			err = errors.New("EOF")
		}
	}()

	if !conn.available {
		return 0, errors.New("EOF")
	}

	data := make([]byte, len(b))
	for i := 0; i < len(b); i++ {
		data[i] = b[i]
	}
	conn.c <- data

	return len(b), nil
}

func (conn *ChannelConn) Close() error {
	defer func() { recover() }()

	if conn.available {
		conn.available = false
		close(conn.c)
	}

	return nil
}
