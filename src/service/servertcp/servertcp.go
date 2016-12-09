package servertcp

import (
	"net"
)

import (
	dm "library/core/datamsg"
	"service"
)

const ServiceName = "servertcp"

type serverTCP struct {
	*service.Service
	output *dm.DataMsgPipe

	listener *net.TCPListener
	ip       string
	port     string
}

func NewServerTCP(name, ip, port string) *serverTCP {
	t := &serverTCP{}
	t.Service = service.NewService(ServiceName)
	t.output = nil
	t.Name = name
	t.ip = ip
	t.port = port
	t.Instance = t
	return t
}
