package util

import (
	"log"
	"net"
	"sync"
)

type ConnMgr struct {
	mutex *sync.Mutex
	conns map[uint16]net.Conn
}

func CreateConnMgr() *ConnMgr {
	return &ConnMgr{new(sync.Mutex), map[uint16]net.Conn{}}
}

func (m *ConnMgr) Get(id uint16) net.Conn {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	c, p := m.conns[id]
	if p {
		return c
	}
	return nil
}

func (m *ConnMgr) Remove(id uint16) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	c, p := m.conns[id]
	if p {
		c.Close()
		delete(m.conns, id)
	}

	return nil
}

func (m *ConnMgr) Add(id uint16, conn net.Conn) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	c, p := m.conns[id]
	if p {
		c.Close()
	}
	m.conns[id] = conn

	return nil
}

func Transmit(one net.Conn, two net.Conn) {
	defer one.Close()
	defer two.Close()

	trans := make([]byte, 4092)

	for {
		length, err := one.Read(trans)
		if err != nil || length == 0 {
			break
		}

		l, err := two.Write(trans[0:length])
		if err != nil {
			break
		}

		if length != l {
			log.Println("write length error", length, l)
		}
	}
}
