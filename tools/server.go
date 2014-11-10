package tools

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"strconv"
)

type Server struct {
	self    Router
	routers []Router

	conf ServerConfig
}

func ReadNetworkFromSocks5(socks5 net.Conn) (loc Location, err error) {
	const (
		IP_KIND     = 1
		DOMAIN_KIND = 3
	)

	const max_socks5_header_len = 255
	buffer := make([]byte, max_socks5_header_len)

	// VER NMETHODS METHODS
	_, err = io.ReadFull(socks5, buffer[0:2])
	if err != nil || buffer[0] != 0x5 {
		return loc, errors.New("socks5 version error")
	}

	_, err = io.ReadFull(socks5, buffer[0:buffer[1]])
	if err != nil {
		return loc, errors.New("socks5 methods error")
	}

	// REPLY: VER METHOD
	socks5.Write([]byte{0x05, 0x00})

	// VER CMD RSV ATYP
	_, err = io.ReadFull(socks5, buffer[0:4])
	if err != nil || buffer[0] != 0x5 || buffer[1] != 0x1 {
		return loc, errors.New("cmd error")
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
			return loc, errors.New("domain len error")
		}
		length = int(buffer[0])
	default:
		return loc, errors.New("unknown address type")
	}

	_, err = io.ReadFull(socks5, buffer[0:length+2])
	if err != nil {
		return loc, errors.New("read address failed")
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
	return loc, nil
}

func (s *Server) Init(conf ServerConfig) {
	s.self.Init(conf.Router)
	s.conf = conf
}

func (s *Server) TransitSocks5(client net.Conn) {
	defer client.Close()

	// socks5 read, target address
	loc, err := ReadNetworkFromSocks5(client)
	if err != nil {
		log.Println(err)
		return
	}

	// select one server
	conn, err := s.CreateRouterSock(loc)
	if err != nil {
		log.Println(err)
		return
	}

	// trans from server
	go Transmit(conn, client)
	Transmit(client, conn)
}

func (s *Server) CreateRouterSock(loc Location) (net.Conn, error) {
	// target is self, direct trans
	if s.self.Match(loc) {
		return net.Dial(loc.Network, loc.String())
	}

	// target is other side, wrapper it
	for _, v := range s.routers {
		if v.Match(loc) {
			// create full-duplex pipe
			return v.Dial(loc.String()), nil
		}
	}
	return nil, errors.New("no such router")
}

func (s *Server) AddRouter(conn net.Conn) {
	var r Router
	var conf RouterConfig

	r.sock = conn
	r.Init(conf)
	r.SendRouterConfig()
	r.Work()
}

func (s *Server) ConnectRouters() {
	for _, v := range s.conf.Routers {
		conn, err := net.Dial("tcp", v)
		if err != nil {
			log.Println("connect error:", v)
		}
		s.AddRouter(conn)
	}
}
