package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
)

import (
	"./protocol"
	"./util"
)

type Config struct {
	RunMode string
	Server  string
	Client  string
	Network string
}

func RunServerMode(config Config) {
	var listener net.Listener
	var err error
	if config.Network == "tcp" {
		listener, err = net.Listen(config.Network, config.Server)
		if err != nil {
			log.Fatalln("listen failed:", err)
		}
	} else {
		listener, err = protocol.WebSocketListen(config.Network, config.Server)
		if err != nil {
			log.Fatalln("listen failed:", err)
		}
	}

	log.Println("running...")

	for {
		client, err := listener.Accept()
		if err != nil {
			log.Fatalln("accept failed:", err)
		}

		log.Println("new connection:", client.RemoteAddr().String())

		go util.ServerReadTunnel(client)
	}
}

func RunClientMode(config Config) {
	var server net.Conn

	conn, err := net.Listen("tcp", config.Client)
	if err != nil {
		log.Fatalln("listen failed:", err)
	}

	if config.Network == "tcp" {
		server, err = net.Dial(config.Network, config.Server)
		if err != nil {
			log.Fatalln("dial failed:", err)
		}
	} else {
		server, err = protocol.WebSocketDial(config.Network, config.Server)
		if err != nil {
			log.Fatalln("dial failed:", err)
		}
	}

	log.Println("running...")

	manager := util.CreateConnMgr()

	go func() {
		for {
			client, err := conn.Accept()
			if err != nil {
				log.Fatalln("accept failed:", err)
			}

			_, port, _ := net.SplitHostPort(client.RemoteAddr().String())
			id, _ := strconv.ParseUint(port, 10, 16)

			manager.Add(uint16(id), client)

			wrapper := protocol.CreateTrans(server, uint16(id))
			go util.Transmit(client, wrapper)
		}
	}()

	util.ClientReadTunnel(server, manager)
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("usage:", os.Args[0], "config.json")
	}

	buffer, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln(os.Args[1], "file not found")
	}

	var config Config
	err = json.Unmarshal(buffer, &config)
	if err != nil {
		log.Fatalln("config:", err)
	}

	if config.RunMode == "server" {
		RunServerMode(config)
	} else {
		RunClientMode(config)
	}
}
