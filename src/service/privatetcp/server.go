package privatetcp

import (
	"io"
	"net"
)

import (
	"library/logger"
)

import (
	cm "library/core/controlmsg"
	dm "library/core/datamsg"
	svc "service/core"
)

type PrivateTCPServer struct {
	svc.Service
	Listener *net.TCPListener
	Conn     *net.TCPConn
	IP       string
	Port     string
}

func NewPrivateTCPServer(bus *dm.DataMsgPipe) *PrivateTCPServer {
	t := &PrivateTCPServer{}
	t.Service = *svc.NewService("")
	t.State = svc.StateInit
	t.ControlMsgPipe = *cm.NewControlMsgPipe()
	t.BUS = bus
	return t
}

func (t *PrivateTCPServer) OnInit(bus *dm.DataMsgPipe) bool {
	logger.Info("Start PrivateTCPServer")
	t.BUS = bus
	tcpAddr, err := net.ResolveTCPAddr("tcp", t.IP+":"+t.Port)
	if err != nil {
		logger.Error("net.ResolveTCPAddr error,%s", err.Error())
		return false
	}

	t.Listener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		logger.Error("net.ListenTCP error,%s", err.Error())
		return false
	}

	logger.Info("listening port:%s", t.Port)

	for {
		t.Conn, err = t.Listener.AcceptTCP()
		if err != nil {
			logger.Error("listener.AcceptTCP error,%s", err.Error())
			continue
		}
		go t.serve()
	}

	return true
}

func (t *PrivateTCPServer) serve() {
	for {
		data := make([]byte, 2)
		n, err := io.ReadAtLeast(t.Conn, data, 2)
		if err != nil {
			logger.Warn("read byte:%d,error:%s", n, err.Error())
		}
		logger.Debug("read byte:%+v", data)
	}
}

func (t *PrivateTCPServer) OnExit() {

}
