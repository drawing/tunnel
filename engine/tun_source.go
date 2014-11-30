package engine

import (
	"log"
	"net"
	"sync/atomic"
)

type Tun struct {
	stream  chan FromConn
	address string
	mode    string
	item    RouterItem
	router  *Router
	def     RouterItem
}

func (t *Tun) SetAddress(mode string, address string) {
	t.mode = mode
	t.address = address
}

func (t *Tun) SetRouterItem(item RouterItem) {
	t.item = item
}
func (t *Tun) SetDefault(def RouterItem) {
	log.Println("SetDefault:", def)
	t.def = def
}

func (t *Tun) SetRouter(router *Router) {
	t.router = router
}

func (t *Tun) Run(stream chan FromConn) {
	t.stream = stream

	if t.mode == "Client" {
		conn, err := net.Dial("tcp", t.address)
		if err != nil {
			log.Println("Tun Dial:", err, t.address)
			return
		}
		loop := NewTunLoop(conn, t.stream, t.router)

		if len(t.item.Domains) > 0 {
			item := t.item
			item.network = NewTunNetwork(loop)
			t.router.AddRouter(t.item)
		}

		if len(t.def.Domains) > 0 {
			loop.Register(t.def)
		}

		go loop.Run()
	} else {
		ln, err := net.Listen("tcp", t.address)
		if err != nil {
			return
		}

		for {
			conn, err := ln.Accept()
			log.Println("One Client UP", err)
			if err == nil {
				go NewTunLoop(conn, t.stream, t.router).Run()
			}
		}
	}
}

type TunLoop struct {
	conn   net.Conn
	stream chan FromConn
	id     uint64
	ctx    map[uint64]net.Conn
	tunnel *TunConn
	router *Router

	diff uint64
}

func NewTunLoop(conn net.Conn, stream chan FromConn, router *Router) *TunLoop {
	loop := &TunLoop{}
	loop.conn = conn
	loop.stream = stream
	loop.ctx = map[uint64]net.Conn{}
	loop.tunnel = NewTunConn(loop.conn, 0)
	loop.router = router

	s1 := conn.LocalAddr().String()
	s2 := conn.RemoteAddr().String()
	if len(s1) > len(s2) {
		loop.diff = 0
	} else if len(s1) < len(s2) {
		loop.diff = uint64(s2[len(s1)])
	} else {
		for i := 0; i < len(s1); i++ {
			if s1[i] != s2[i] {
				loop.diff = uint64(s1[i])
				break
			}
		}
	}

	log.Println(s1, s2, loop.diff)
	loop.diff = loop.diff << 32
	loop.id = loop.diff
	log.Println(s1, s2, loop.diff)

	return loop
}

func (t *TunLoop) Connect(loc Location) (net.Conn, error) {
	// TODO: forbid newID out of bound
	newID := atomic.AddUint64(&t.id, 1)

	ch := NewChannelConn()
	tu := t.tunnel.Clone()
	tu.SetID(newID)

	var pkg Package
	pkg.Command = PkgCommandConnect
	pkg.Id = newID
	pkg.Loc = &loc

	t.tunnel.WritePackage(&pkg)

	t.ctx[newID] = ch

	conn := NewPipeConn(ch, tu)
	return conn, nil
}

func (t *TunLoop) Register(item RouterItem) {
	var pkg Package
	pkg.Command = PkgCommandRegister
	pkg.Router = &item

	t.tunnel.WritePackage(&pkg)
}

func (t *TunLoop) Run() {
	for {
		var pkg Package
		err := t.tunnel.ReadPackage(&pkg)
		if err != nil {
			break
		}

		switch pkg.Command {
		case PkgCommandConnect:
			log.Println("Connect", pkg)
			// new connection
			var from FromConn
			if pkg.Loc == nil {
				continue
			}

			from.Loc = *pkg.Loc

			ch := NewChannelConn()
			tu := t.tunnel.Clone()
			tu.SetID(pkg.Id)

			t.ctx[pkg.Id] = ch
			from.Conn = NewPipeConn(ch, tu)

			t.stream <- from
		case PkgCommandData:
			to, present := t.ctx[pkg.Id]
			if !present {
				continue
			}
			_, err := to.Write(pkg.Data)
			if err != nil {
				// delete the ctx item
				continue
			}
		case PkgCommandRegister:
			if pkg.Router == nil || pkg.Router.Domains == nil {
				log.Println("Register: router is nil")
				continue
			}

			var item RouterItem
			item.Domains = pkg.Router.Domains
			item.network = NewTunNetwork(t)
			t.router.AddRouter(item)
			log.Println("Add Client Router:", item)
		}
	}
}
