package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
)

func ReadSocks5(conn net.Conn) (ip, port string, err error) {
	const (
		IP_KIND     = 1
		DOMAIN_KIND = 3
	)

	const max_socks5_header_len = 255
	buffer := make([]byte, max_socks5_header_len)

	// VER NMETHODS METHODS
	_, err = io.ReadFull(conn, buffer[0:2])
	if err != nil || buffer[0] != 0x5 {
		return ip, port, errors.New("socks5 version error")
	}

	_, err = io.ReadFull(conn, buffer[0:buffer[1]])
	if err != nil {
		return ip, port, errors.New("socks5 methods error")
	}

	// REPLY: VER METHOD
	conn.Write([]byte{0x05, 0x00})

	// VER CMD RSV ATYP
	_, err = io.ReadFull(conn, buffer[0:4])
	if err != nil || buffer[0] != 0x5 || buffer[1] != 0x1 {
		return ip, port, errors.New("cmd error")
	}

	// only support tcp now
	length := 4
	kind := buffer[3]
	switch kind {
	case IP_KIND:
	case DOMAIN_KIND:
		_, err = io.ReadFull(conn, buffer[0:1])
		if err != nil {
			return ip, port, errors.New("domain len error")
		}
		length = int(buffer[0])
	default:
		return ip, port, errors.New("unknown address type")
	}

	_, err = io.ReadFull(conn, buffer[0:length+2])
	if err != nil {
		return ip, port, errors.New("read address failed")
	}

	var addr_port uint16 = 80
	binary.Read(bytes.NewBuffer(buffer[length:length+2]), binary.BigEndian, &addr_port)
	switch kind {
	case IP_KIND:
		ip = net.IPv4(buffer[0], buffer[1], buffer[2], buffer[3]).String()
		port = strconv.Itoa(int(addr_port))
	case DOMAIN_KIND:
		port = strconv.Itoa(int(addr_port))
		ip = string(buffer[0:length])
	}

	conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0xca, 0x0a, 0x0a, 0xca, 0x33, 0x33})

	return ip, port, nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	ip, port, err := ReadSocks5(conn)
	if err != nil {
		log.Println("socks5 failed:", err)
		return
	}

	log.Println(ip, port)

	if port == "80" {
		doHTTPRequest(conn, 1, ip)
	} else if port == "443" {
		doHTTPSRequest(conn, ip)
	} else {
		log.Println("unkown protocol", ip, port)
	}
}

func doHTTPSRequest(conn net.Conn, ip string) {
	// var config tls.Config
	config := &tls.Config{Certificates: []tls.Certificate{ /*cert*/ }}
	sec := tls.Server(conn, config)
	doHTTPRequest(sec, 2, ip)
}

func doHTTPRequest(conn net.Conn, protocol int, ip string) {
	reader := bufio.NewReader(conn)
	req, err := http.ReadRequest(reader)
	if err != nil {
		log.Println("http failed:", err)
		return
	}

	var url string
	if protocol == 1 {
		url = "http://" + ip + req.RequestURI
	} else {
		url = "https://" + ip + req.RequestURI
	}

	log.Println("REQ:", url)

	client := &http.Client{}

	// construct request
	send_req, err := http.NewRequest(req.Method, url, req.Body)
	if err != nil {
		log.Println("new req:", err)
		return
	}
	send_req.Header = req.Header

	// redirect header

	// redirect body

	resp, err := client.Do(send_req)
	if err != nil {
		log.Println("request failed:", err)
		return
	}

	// write back header

	// write back body

	err = resp.Write(conn)
	if err != nil {
		log.Println("wirte resp failed:", err)
		return
	}
}

func main() {
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalln("listen failed:", err)
	}

	log.Println("(9000) Running...")
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept failed:", err)
			continue
		}
		go handleConnection(conn)
	}
}
