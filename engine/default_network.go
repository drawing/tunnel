package engine

import (
	"log"
	"net"
)

type DefaultNetwork struct {
}

func (n *DefaultNetwork) Dial(loc Location) (net.Conn, error) {
	log.Println("Default:", loc)
	return net.Dial(loc.Network, loc.String())
}

func (n *DefaultNetwork) ID() uint64 {
	return 0
}
