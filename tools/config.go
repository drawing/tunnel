package tools

type Network struct {
	Protocol string
	Location string
	Domains  []string
}

type TunConfig struct {
	Outgoing []Network
	Incoming []Network
}
