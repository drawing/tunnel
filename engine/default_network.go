package engine

import (
	"errors"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type dnsRecord struct {
	IP         string
	LastUpdate time.Time
}

type DefaultNetwork struct {
	dns   map[string]dnsRecord
	mutex sync.RWMutex
}

func NewDefaultNetwork() *DefaultNetwork {
	n := &DefaultNetwork{}
	n.dns = map[string]dnsRecord{}
	return n
}

func (n *DefaultNetwork) Dial(loc Location) (net.Conn, error) {
	log.Println("Connect:", loc.Network+"@"+loc.String())

	http_proxy := os.Getenv("http_proxy")
	if len(http_proxy) != 0 {
		strs := strings.Split(http_proxy, ":")
		if len(strs) > 1 {
			loc.Address = strs[0]
			loc.Port = strs[1]
			log.Println("Use Proxy:", loc.Address, loc.Port)
		}
	}

	address := loc.Address

	if len(address) == 0 {
		var record dnsRecord
		n.mutex.RLock()
		record, ok := n.dns[loc.Address]
		n.mutex.RUnlock()

		// 30 minute cache
		if !ok || time.Since(record.LastUpdate).Minutes() > 30 {
			list, err := net.LookupIP(loc.Domain)
			if err != nil || len(list) < 1 {
				log.Println("Host not found:", loc.Domain, err)
				return nil, errors.New("host not found")
			}

			record.IP = list[0].String()
			record.LastUpdate = time.Now()
		}

		address = record.IP

		if len(n.dns) > 10000 {
			// clear cache
			n.dns = map[string]dnsRecord{}
		}
	}

	return net.Dial(loc.Network, address+":"+loc.Port)
}

func (n *DefaultNetwork) ID() uint64 {
	return 0
}
