package util

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
	"strconv"
)

import (
	"../protocol"
)

func ReadFromSock5(sock5 net.Conn, client net.Conn) {
	defer client.Close()
	defer sock5.Close()

	// VER NMETHODS METHODS
	buffer := make([]byte, protocol.MaxBufferSize)

	_, err := io.ReadFull(sock5, buffer[0:2])
	if err != nil || buffer[0] != 0x5 {
		log.Println("sock5 read header failed:", err)
		return
	}

	_, err = io.ReadFull(sock5, buffer[0:buffer[1]])
	if err != nil {
		log.Println("read NMethods failed:", err)
		return
	}

	// VER METHOD
	client.Write([]byte{0x05, 0x00})

	_, err = io.ReadFull(sock5, buffer[0:4])
	if err != nil || buffer[0] != 0x5 || buffer[1] != 0x1 {
		log.Println("read CMD:", err, buffer[0:4])
		return
	}

	addr_len := 4
	addr_type := buffer[3]
	if addr_type == 1 {
		// IPv4
	} else if addr_type == 3 {
		// Domain
		_, err = io.ReadFull(sock5, buffer[0:1])
		if err != nil {
			log.Println("Read CMD", err)
			return
		}
		addr_len = int(buffer[0])
	} else {
		log.Println("read addr type error", buffer[0:4])
		return
	}

	_, err = io.ReadFull(sock5, buffer[0:addr_len+2])
	if err != nil {
		log.Println("Read Addr", err)
		return
	}

	addr_string := ""
	var addr_port uint16 = 80
	binary.Read(bytes.NewBuffer(buffer[addr_len:addr_len+2]), binary.BigEndian, &addr_port)
	if addr_type == 1 {
		addr_string = net.IPv4(buffer[0], buffer[1], buffer[2], buffer[3]).String() +
			":" + strconv.Itoa(int(addr_port))

	} else {
		addr_string = string(buffer[0:addr_len]) + ":" + strconv.Itoa(int(addr_port))
	}

	log.Println("CONNECT:", addr_string)

	target, err := net.Dial("tcp", addr_string)
	if err != nil {
		log.Println("Connect", err, addr_string)
		return
	}

	defer target.Close()

	client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0xca, 0x0a, 0x0a, 0xca, 0x33, 0x33})

	go Transmit(target, client)

	for {
		length, err := sock5.Read(buffer)
		if err != nil {
			break
		}

		_, err = target.Write(buffer[0:length])
		if err != nil {
			break
		}
	}
}

func ServerReadTunnel(client net.Conn) {
	defer client.Close()

	manager := CreateConnMgr()

	wrapper := protocol.CreateTrans(client, uint16(0))

	data := make([]byte, protocol.MaxBufferSize)

	for {
		length, err := wrapper.Read(data)
		if err != nil {
			break
		}

		if length == 0 {
			manager.Remove(wrapper.Id)
			continue
		}

		side := manager.Get(wrapper.Id)
		if side == nil {
			side = protocol.CreateChannel()
			manager.Add(wrapper.Id, side)
			go ReadFromSock5(side, protocol.CreateTrans(client, wrapper.Id))
		}

		_, err = side.Write(data[0:length])
		if err != nil {
			manager.Remove(wrapper.Id)
		}
	}

	log.Println("connection close")
}
