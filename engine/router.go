package engine

import (
	"log"
	"net"
	"regexp"
)

type Location struct {
	Network string
	Domain  string
	Address string
	Port    string
}

func (l *Location) String() string {
	if l.Address == "" {
		return l.Domain + ":" + l.Port
	}
	return l.Address + ":" + l.Port
}

type Router struct {
	dynamic []RouterItem
	other   Network
}

type RouterItem struct {
	//match
	Domains []string
	iplist  []*net.IPNet

	network Network
}

func (r *Router) Init() {
	r.dynamic = make([]RouterItem, 0, 10)
}

func (r *Router) SetDefault(network Network) {
	r.other = network
}

func (r *Router) AddRouter(item RouterItem) {
	log.Println("AddRouter", item)
	r.dynamic = append(r.dynamic, item)
}

func (r *Router) Match(loc Location) Network {
	log.Println("Match", loc, len(r.dynamic))
	for _, item := range r.dynamic {
		if loc.Domain != "" {
			// match domain
			log.Println("Match 1", item.Domains)
			for _, v := range item.Domains {
				matched, _ := regexp.MatchString(v, loc.Domain)
				if matched {
					return item.network
				}
			}
		} else {
			// match address
			for _, v := range item.Domains {
				matched, _ := regexp.MatchString(v, loc.Address)
				if matched {
					return item.network
				}
			}
			/*
				ip := net.ParseIP(loc.Address)
				if ip == nil {
					log.Println("loc address:", loc.Address)
					return false
				}
				for _, v := range r.iplist {
					if v.Contains(ip) {
						return true
					}
				}
			*/
		}
	}

	return r.other
}
