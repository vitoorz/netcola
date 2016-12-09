package gatewayoutter

import (
	dm "library/core/datamsg"
	"net"
	"service"
)

const ServiceName = "gatewayoutter"

type gatewayOutter struct {
	*service.Service
	output *dm.DataMsgPipe

	listener *net.TCPListener
	ip       string
	port     string
}

func NewGatewayOutter(name, ip, port string) *gatewayOutter {
	t := &gatewayOutter{}
	t.Service = service.NewService(ServiceName)
	t.output = nil
	t.Name = name
	t.ip = ip
	t.port = port
	t.Instance = t
	return t
}
