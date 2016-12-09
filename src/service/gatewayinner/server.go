package gatewayinner

import (
	dm "library/core/datamsg"
	"net"
	"service"
)

const ServiceName = "gatewayinner"

type gatewayInner struct {
	*service.Service
	output *dm.DataMsgPipe

	listener *net.TCPListener
	ip       string
	port     string
}

func NewGatewayInner(name, ip, port string) *gatewayInner {
	t := &gatewayInner{}
	t.Service = service.NewService(ServiceName)
	t.output = nil
	t.Name = name
	t.ip = ip
	t.port = port
	t.Instance = t
	return t
}
