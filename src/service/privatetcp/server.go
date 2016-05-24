package privatetcp

import (
	"io"
	"net"
	"sync"
)

import (
	"library/logger"
)

import (
	cm "library/core/controlmsg"
	dm "library/core/datamsg"
)

type PrivateTCPServer struct {
	Lock sync.RWMutex
	cm.ControlMsgPipe
	BUS *dm.DataMsgPipe

	ID       uint64
	Listener *net.TCPListener
	Conn     *net.TCPConn
	IP       string
	Port     string
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
	}
}
