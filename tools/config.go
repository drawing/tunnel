package tools

type ServerConfig struct {
	Router     RouterConfig
	MessageInf string
	DataInf    string
}

type RouterConfig struct {
	IPNetList  []string
	DomainList []string
}
