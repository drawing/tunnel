package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

import (
	"./engine"
)

/*
TODO:
1. fix. change id rule
2. fix. mutex
3. conn quit
4. sec tunnel
5. fix. client quit, remove router
6. fix recycle
7. format struct
*/

type SourceConfig struct {
	Category string
	Location string
	Protocol string
	SecFile  string
}

type RouterConfig struct {
	Domains []string
}

type ItemConifg struct {
	Source *SourceConfig
	Router *RouterConfig
}

type TunConfig struct {
	Sources []ItemConifg
	Default *RouterConfig
}

var eng engine.Engine
var config TunConfig

func RunServerMode(config TunConfig) {
	router := engine.NewRouter()
	eng.SetRouter(router)

	log.Println("prepare...", config)

	for _, v := range config.Sources {
		if v.Source == nil {
			log.Println("CONF:", "source is nil", v)
			continue
		}
		switch v.Source.Category {
		case "Socks5":
			sour := &engine.Socks5Source{}
			sour.SetAddress(v.Source.Location)
			eng.AddSource(sour)
		case "ConnectTunnel":
			sour := &engine.Tun{}
			sour.SetAddress("Client", v.Source.Location)
			sour.SetRouter(router)
			if v.Router != nil {
				var item engine.RouterItem
				item.Domains = v.Router.Domains
				sour.SetRouterItem(item)
			}

			if config.Default != nil {
				var def engine.RouterItem
				def.Domains = config.Default.Domains
				sour.SetDefault(def)
			}

			eng.AddSource(sour)
		case "ListenTunnel":
			sour := &engine.Tun{}
			sour.SetAddress("Server", v.Source.Location)
			sour.SetRouter(router)
			eng.AddSource(sour)
		}
	}

	router.SetDefault(&engine.DefaultNetwork{})

	log.Println("running...")

	eng.Run()
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("usage:", os.Args[0], "config.json")
	}

	buffer, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln(os.Args[1], "file not found")
	}

	err = json.Unmarshal(buffer, &config)
	if err != nil {
		log.Fatalln("config:", err)
	}

	eng.Init()
	RunServerMode(config)
}
