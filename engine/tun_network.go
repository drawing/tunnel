package engine

import (
	"log"
	"net"
)

type TunNetwork struct {
	loop *TunLoop
}

func NewTunNetwork(loop *TunLoop) Network {
	n := &TunNetwork{loop}
	return n
}

func (n *TunNetwork) Dial(loc Location) (net.Conn, error) {
	addr := n.loop.RemoteAddr()
	log.Println("Redirect:", loc.Network+"@"+loc.String(), "->",
		addr.Network()+"@"+addr.String())
	return n.loop.Connect(loc)
}

func (n *TunNetwork) ID() uint64 {
	return n.loop.UniID
}
