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

func RunMessage(control net.Listener) {
	for {
		client, err := control.Accept()
		if err != nil {
			log.Fatalln("control accept failed:", err)
		}

		log.Println("new server:", client.RemoteAddr().String())

		go server.AddRouter(client)
	}
}

func RunServerMode(config tools.TunConfig) {
	var listener net.Listener
	var err error

	listener, err = net.Listen("tcp", config.DataInf)
	if err != nil {
		log.Fatalln("listen failed:", err)
	}

	control, err := net.Listen("tcp", config.MessageInf)
	if err != nil {
		log.Fatalln("control failed:", err)
	}

	log.Println("connent...")

	server.ConnectRouters()

	log.Println("running...")

	go RunMessage(control)

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

	var config tools.TunConfig
	err = json.Unmarshal(buffer, &config)
	if err != nil {
		log.Fatalln("config:", err)
	}

	RunServerMode(config)
}
