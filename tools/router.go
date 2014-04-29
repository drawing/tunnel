package tools

import (
	"code.google.com/p/go.net/proxy"
	"net"
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
	domains []string
	sock    net.Conn
	writer  net.Conn
	runner  map[int64]net.Conn
	current int64
}

func (r *Router) Match(loc Location) bool {
	return true
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
