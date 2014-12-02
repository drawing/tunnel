package engine

import (
	"testing"
)

func TestChannelConn(t *testing.T) {
	conn := NewChannelConn()

	conn.Write([]byte("123456"))
	conn.Write([]byte("2"))

	v := make([]byte, 10)
	n, err := conn.Read(v[0:1])
	if err != nil || n != 1 || v[0] != '1' {
		t.Fail()
	}
	n, err = conn.Read(v[0:1])
	if err != nil || n != 1 || v[0] != '2' {
		t.Fail()
	}
	n, err = conn.Read(v[0:4])
	if err != nil || n != 4 || string(v[0:4]) != "3456" {
		t.Fail()
	}

	conn.Close()
	conn.Close()
}
