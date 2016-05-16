package engine

import (
	// "errors"
	"bufio"
	"log"
	"net"
	"net/http"
	"strings"
)

type HttpSource struct {
	address string
	stream  chan FromConn
}

func (s *HttpSource) SetAddress(addr string) {
	s.address = addr
}

func (s *HttpSource) Run(stream chan FromConn) {
	log.Println("HTTP Running:", s.address)

	ln, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Println("http listen failed", s.address, err)
		return
	}

	s.stream = stream

	for {
		conn, err := ln.Accept()
		if err == nil {
			go s.ReadHttp(conn)
		}
	}
}

func (s *HttpSource) ReadHttp(http_client net.Conn) error {
	const (
		IP_KIND     = 1
		DOMAIN_KIND = 3
	)

	var loc Location

	client := NewPeekConn(http_client, http_client)

	req, err := http.ReadRequest(bufio.NewReader(client))
	if err != nil {
		return err
	}

	log.Println("HTTP Proxy:", req.URL.Host)

	client.(*PeekConn).PeekMode = false

	loc.Network = "tcp"
	loc.Address = ""

	strs := strings.Split(req.URL.Host, ":")
	if len(strs) <= 1 {
		loc.Domain = req.URL.Host
		loc.Port = "80"
	} else {
		loc.Domain = strs[0]
		loc.Port = strs[1]
	}

	if req.Method == "CONNECT" {
		client.(*PeekConn).Clear()
		client.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
		log.Println("Connect Method")
	}

	/*
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

	*/
	var from FromConn
	from.Loc = loc
	from.Conn = client

	// send to engine
	s.stream <- from
	return nil
}
