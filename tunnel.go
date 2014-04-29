package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"os"
)

import (
	"./tools"
)

var server tools.Server

func RunServerMode(config tools.ServerConfig) {
	var listener net.Listener
	var err error

	listener, err = net.Listen("tcp", config.DataInf)
	if err != nil {
		log.Fatalln("listen failed:", err)
	}

	log.Println("running...")

	for {
		client, err := listener.Accept()
		if err != nil {
			log.Fatalln("accept failed:", err)
		}

		log.Println("new connection:", client.RemoteAddr().String())

		go server.TransitSocks5(client)
	}
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("usage:", os.Args[0], "config.json")
	}

	buffer, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln(os.Args[1], "file not found")
	}

	var config tools.ServerConfig
	err = json.Unmarshal(buffer, &config)
	if err != nil {
		log.Fatalln("config:", err)
	}

	server.Init(config)
	RunServerMode(config)
}
