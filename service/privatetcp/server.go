package privatetcp

import (
	"net"
)

import (
	dm "netcola/library/core/datamsg"
	"netcola/service"
)

const ServiceName = "privatetcpserver"

type privateTCPServer struct {
	*service.Service
	output *dm.DataMsgPipe

	listener *net.TCPListener
	ip       string
	port     string
}

func NewPrivateTCPServer(name, ip, port string) *privateTCPServer {
	t := &privateTCPServer{}
	t.Service = service.NewService(ServiceName)
	t.output = nil
	t.Name = name
	t.ip = ip
	t.port = port
	t.Instance = t
	return t
}
