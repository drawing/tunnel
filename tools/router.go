package tools

import (
	"code.google.com/p/go.net/proxy"
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
	sock    net.Conn
	writer  net.Conn
	runner  map[int64]net.Conn
	current int64

	//match
	namelist []string
	iplist   []*net.IPNet

	// conf
	conf RouterConfig
}

func (r *Router) Init(conf RouterConfig) {
	r.namelist = conf.DomainList
	r.iplist = []*net.IPNet{}
	r.runner = map[int64]net.Conn{}
	r.conf = conf

	for _, v := range conf.IPNetList {
		_, net, err := net.ParseCIDR(v)
		if err != nil {
			log.Println("IPNet invalid:", v)
			continue
		}
		r.iplist = append(r.iplist, net)
	}
}

func (r *Router) Match(loc Location) bool {
	if loc.Address == "" {
		// match domain
		for _, v := range r.namelist {
			matched, _ := regexp.MatchString(v, loc.Domain)
			if matched {
				return true
			}
		}
	} else {
		// match address
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
	}
	return false
}

func (r *Router) GenerateID() int64 {
	r.current += 1
	return r.current
}

func (r *Router) Work() {
	go r.ReadFromRemote()
}

func (r *Router) Dial(addr string) net.Conn {
	// generate unique id
	id := r.GenerateID()
	// build pipe
	reader := CreateChannel()
	// send new connection req

	// connect to server
	wrapper := CreateWrapper(r.sock, id)
	// wrapper.Connect()
	conn := CreatePipe(reader, wrapper)

	return conn
}

func (r *Router) ReadFromRemote() {
	defer r.sock.Close()
	wrapper := CreateWrapper(r.sock, 0)

	var pkg Package
	for {
		err := wrapper.ReadPackage(&pkg)
		if err == nil {
			break
		}

		switch pkg.Kind {
		case PKG_KIND_CONNECT:
			// connect
			dialer, _ := proxy.SOCKS5("tcp", "localhost", nil, new(net.Dialer))
			conn, _ := dialer.Dial("tcp", "addr")
			id := r.GenerateID()
			r.runner[id] = conn

			// conn -> remote
			wrapper := CreateWrapper(r.sock, id)
			go Transmit(conn, wrapper)
			// gateway
		case PKG_KIND_DATA:
			conn, present := r.runner[pkg.Id]
			if present {
				_, err := conn.Write(pkg.Data)
				if err == nil {
					continue
				}
			}
		case PKG_KIND_CLOSE:
			// close
			conn, present := r.runner[pkg.Id]
			if present {
				conn.Close()
			}
			delete(r.runner, pkg.Id)
		}
	}
}
