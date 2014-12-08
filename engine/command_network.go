package engine

import (
	"log"
	"net"
)

type CommandNetwork struct {
}

func (n *CommandNetwork) Dial(loc Location) (net.Conn, error) {
	log.Println("CommandNet:", loc.Network+"@"+loc.String())
	return NewCommand(), nil
}

func (n *CommandNetwork) ID() uint64 {
	return 0
}
