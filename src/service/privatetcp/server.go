package privatetcp

import (
	"io"
	"net"
)

import (
	dm "library/core/datamsg"
	"library/logger"
	"service"
	"time"
)

const ServiceName = "privatetcpserver"

type PrivateTCPServer struct {
	service.Service
	Listener *net.TCPListener
	ConnList []*net.TCPConn
	IP       string
	Port     string
}

func NewPrivateTCPServer(bus *dm.DataMsgPipe) *PrivateTCPServer {
	t := &PrivateTCPServer{}
	t.Service = *service.NewService(ServiceName)
	t.BUS = nil
	t.IP = "0.0.0.0"
	t.Port = "7171"
	return t
}

func (t *PrivateTCPServer) Start(name string, bus *dm.DataMsgPipe) bool {
	logger.Info("Start PrivateTCPServer")
	t.Name = name
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
	go t.serve()
	return true
}

func (t *PrivateTCPServer) Pause() bool {
	return true
}

func (t *PrivateTCPServer) Resume() bool {
	return true
}

func (t *PrivateTCPServer) Exit() bool {
	return true
}

func (t *PrivateTCPServer) serve() {
	for {
		connect, err := t.Listener.AcceptTCP()
		if err != nil {
			logger.Error("listener.AcceptTCP error,%s", err.Error())
			time.Sleep(time.Second * 2)
			continue
		}
		go t.readConn(connect)
		go t.writeConn(connect)

		//t.Conn.Close()
	}
}

func (t *PrivateTCPServer) readConn(connection *net.TCPConn) {
	defer func() {
		connection.Close()
	}()
	for {
		data := make([]byte, 1)
		n, err := io.ReadAtLeast(connection, data, 1)
		if err != nil {
			logger.Warn("read byte:%d,error:%s", n, err.Error())
			return
		}
		logger.Debug("read %d byte:%+v", n, data)
		t.BUS.Down <- dm.NewDataMsg("job", 0, data)
	}
}

func (t *PrivateTCPServer) writeConn(connection *net.TCPConn) {
	defer func() {
		connection.Close()
	}()
	for {
		select {
		case data, ok := <-t.BUS.Up:
			if !ok {
				logger.Info("Data Read error")
				break
			}
			count, err := connection.Write(data.Payload.([]byte))
			if err != nil {
				logger.Warn("conn write err:%s", err.Error())
				return
			} else {
				logger.Info("send bus.up:%d byte", count)
			}
		}
	}
}
