package engine

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"strconv"
)

type Socks5Source struct {
	address string
	stream  chan FromConn
}

func (s *Socks5Source) SetAddress(addr string) {
	s.address = addr
}

func (s *Socks5Source) Run(stream chan FromConn) {
	log.Println("Socks5Source:", s.address)
	ln, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Println("Socks5Source:", "listen failed", s.address, err)
		return
	}

	s.stream = stream

	for {
		conn, err := ln.Accept()
		if err == nil {
			go s.ReadSocks5(conn)
		}
	}
}

func (s *Socks5Source) ReadSocks5(socks5 net.Conn) error {
	const (
		IP_KIND     = 1
		DOMAIN_KIND = 3
	)

	var loc Location

	const max_socks5_header_len = 255
	buffer := make([]byte, max_socks5_header_len)

	// VER NMETHODS METHODS
	_, err := io.ReadFull(socks5, buffer[0:2])
	if err != nil || buffer[0] != 0x5 {
		return errors.New("socks5 version error")
	}

	_, err = io.ReadFull(socks5, buffer[0:buffer[1]])
	if err != nil {
		return errors.New("socks5 methods error")
	}

	// REPLY: VER METHOD
	socks5.Write([]byte{0x05, 0x00})

	// VER CMD RSV ATYP
	_, err = io.ReadFull(socks5, buffer[0:4])
	if err != nil || buffer[0] != 0x5 || buffer[1] != 0x1 {
		return errors.New("cmd error")
	}

	// only support tcp now
	loc.Network = "tcp"

	length := 4
	kind := buffer[3]
	switch kind {
	case IP_KIND:
	case DOMAIN_KIND:
		_, err = io.ReadFull(socks5, buffer[0:1])
		if err != nil {
			return errors.New("domain len error")
		}
		length = int(buffer[0])
	default:
		return errors.New("unknown address type")
	}

	_, err = io.ReadFull(socks5, buffer[0:length+2])
	if err != nil {
		return errors.New("read address failed")
	}

	var addr_port uint16 = 80
	binary.Read(bytes.NewBuffer(buffer[length:length+2]), binary.BigEndian, &addr_port)
	switch kind {
	case IP_KIND:
		loc.Address = net.IPv4(buffer[0], buffer[1], buffer[2], buffer[3]).String()
		loc.Port = strconv.Itoa(int(addr_port))
		loc.Domain = ""
	case DOMAIN_KIND:
		loc.Address = ""
		loc.Port = strconv.Itoa(int(addr_port))
		loc.Domain = string(buffer[0:length])
	}

	socks5.Write([]byte{0x05, 0x00, 0x00, 0x01, 0xca, 0x0a, 0x0a, 0xca, 0x33, 0x33})

	var from FromConn
	from.Loc = loc
	from.Conn = socks5

	log.Println("socks5 in:", from.Loc)

	// send to engine
	s.stream <- from
	return nil
}
