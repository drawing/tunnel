package engine

import (
	"log"
	"net"
)

type FromConn struct {
	Conn net.Conn
	Loc  Location
}

type Source interface {
	Run(chan FromConn)
}

type Network interface {
	Dial(loc Location) (net.Conn, error)
}

type Engine struct {
	stream chan FromConn
	router *Router
}

func (e *Engine) Init() {
	e.stream = make(chan FromConn, 20)
}

func (e *Engine) SetRouter(router *Router) {
	e.router = router
}

func (e *Engine) AddSource(source Source) {
	go source.Run(e.stream)
}

func (e *Engine) Accept() (FromConn, error) {
	from := <-e.stream
	return from, nil
}

func (e *Engine) Run() error {
	log.Println("Engine Running...")
	for {
		from, err := e.Accept()
		if err != nil {
			return err
		}

		network := e.router.Match(from.Loc)

		log.Println(from.Loc)
		if network == nil {
			continue
		}

		go e.transform(from, network)
	}
}

func (e *Engine) transform(from FromConn, network Network) error {
	to, err := network.Dial(from.Loc)
	if err != nil {
		log.Println("Dail:", from.Loc, "failed,", err)
		return err
	}

	go e.transmit(to, from.Conn)
	e.transmit(from.Conn, to)

	return nil
}

func (e *Engine) transmit(from net.Conn, to net.Conn) {
	defer from.Close()
	defer to.Close()

	trans := make([]byte, 4092)

	for {
		length, err := from.Read(trans)
		if err != nil || length == 0 {
			break
		}

		l, err := to.Write(trans[0:length])
		if err != nil {
			break
		}

		if length != l {
			log.Println("write length error", length, l)
		}
	}
}
