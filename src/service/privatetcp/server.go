package privatetcp

import (
	"net"
)

import (
	dm "library/core/datamsg"
	"service"
)

const ServiceName = "privatetcpserver"

type privateTCPServer struct {
	service.Service
	output *dm.DataMsgPipe

	listener *net.TCPListener
	ip       string
	port     string
}

func NewPrivateTCPServer(name, ip, port string) *privateTCPServer {
	t := &privateTCPServer{}
	t.Service = *service.NewService(ServiceName)
	t.output = nil
	t.Name = name
	t.ip = ip
	t.port = port
	t.Buffer = service.NewBufferPool(&t.Service)
	t.Instance = t
	return t
}
