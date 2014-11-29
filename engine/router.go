package engine

import (
	"log"
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
	// match use regexp
	Domains []string
	network Network
}

func NewRouter() *Router {
	r := &Router{}
	r.dynamic = make([]RouterItem, 0, 10)
	return r
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
		}
	}

	return r.other
}
