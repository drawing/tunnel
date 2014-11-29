package engine

import (
	"log"
	"net"
	"sync/atomic"
)

const (
	CommandConnect = 1
)

type Tun struct {
	stream  chan FromConn
	address string
	mode    string
	item    RouterItem
	router  *Router
}

func (t *Tun) SetAddress(mode string, address string) {
	t.mode = mode
	t.address = address
}

func (t *Tun) SetRouterItem(item RouterItem) {
	t.item = item
}

func (t *Tun) SetRouter(router *Router) {
	t.router = router
}

func (t *Tun) Run(stream chan FromConn) {
	if t.mode == "Client" {
		conn, err := net.Dial("tcp", t.address)
		if err != nil {
			log.Println("Tun Dial:", err, t.address)
			return
		}
		loop := NewTunnLoop(conn, t.stream)
		item := t.item
		item.network = NewTunNetwork(loop)
		t.router.AddRouter(item)
		go loop.Run()
	} else {
		ln, err := net.Listen("tcp", t.address)
		if err != nil {
			return
		}

		t.stream = stream

		for {
			conn, err := ln.Accept()
			log.Println("One Client UP", err)
			if err == nil {
				go NewTunnLoop(conn, t.stream).Run()
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
}

func NewTunnLoop(conn net.Conn, stream chan FromConn) *TunLoop {
	loop := &TunLoop{}
	loop.conn = conn
	loop.stream = stream
	loop.ctx = map[uint64]net.Conn{}
	loop.tunnel = NewTunConn(loop.conn, 0)
	return loop
}

func (t *TunLoop) Connect(loc Location) (net.Conn, error) {
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
		case PkgCommandRouter:
			/*if pkg.Router == nil {
				continue
			}
			Router*/
		}
	}
}
