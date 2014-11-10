package tools

type ServerConfig struct {
	Router     RouterConfig
	MessageInf string
	DataInf    string
	Routers    []string
}

type RouterConfig struct {
	IPNetList  []string
	DomainList []string
}
