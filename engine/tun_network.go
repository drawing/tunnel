package engine

import (
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
	return n.loop.Connect(loc)
}
