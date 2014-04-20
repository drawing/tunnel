package util

import (
	"log"
	"net"
)

import (
	"../protocol"
)

func ClientReadTunnel(server net.Conn, manager *ConnMgr) {
	defer server.Close()

	wrapper := protocol.CreateTrans(server, uint16(0))

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
			continue
		}

		_, err = side.Write(data[0:length])
		if err != nil {
			manager.Remove(wrapper.Id)
			continue
		}
	}
	log.Println("remote closed")
}
