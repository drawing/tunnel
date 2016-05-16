package engine

import (
	"log"
	"regexp"
	"sync"
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

	mutex sync.RWMutex
}

type RouterItem struct {
	// match use regexp
	Domains []string
	network Network
}

func NewRouterItem(domains []string, network Network) *RouterItem {
	r := &RouterItem{}
	r.Domains = domains
	r.network = network
	return r
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
	r.mutex.Lock()
	defer r.mutex.Unlock()

	log.Println("add router:(", item.network, ")", item.Domains)

	r.dynamic = append(r.dynamic, item)
}

func (r *Router) RemoveRouter(id uint64) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	index := 0
	find := false

	for k, v := range r.dynamic {
		if v.network.ID() == id {
			find = true
			index = k
			break
		}
	}

	if find {
		log.Println("remove router(", r.dynamic[index].network, "):", r.dynamic[index].Domains)
		r.dynamic = append(r.dynamic[0:index], r.dynamic[index+1:]...)
	}
}

func (r *Router) Match(loc Location) Network {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// log.Println("Match", loc, len(r.dynamic))
	for _, item := range r.dynamic {
		if loc.Domain != "" {
			// match domain
			// log.Println("Match 1", item.Domains)
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
