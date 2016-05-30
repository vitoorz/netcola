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

	"types"
)

const ServiceName = "privatetcpserver"

type PrivateTCPServer struct {
	service.Service
	Listener *net.TCPListener
	Output   *dm.DataMsgPipe
	IP       string
	Port     string
}

func NewPrivateTCPServer() *PrivateTCPServer {
	t := &PrivateTCPServer{}
	t.Service = *service.NewService(ServiceName)
	t.IP = "0.0.0.0"
	t.Port = "7171"
	return t
}

func (t *PrivateTCPServer) Start(name string, bus *dm.DataMsgPipe) bool {
	logger.Info("Start PrivateTCPServer")
	t.Name = name
	t.Output = bus
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
	go t.writeConn()
	for {
		connect, err := t.Listener.AcceptTCP()
		if err != nil {
			logger.Error("listener.AcceptTCP error,%s", err.Error())
			time.Sleep(time.Second * 2)
			connect.Close()
			continue
		}
		go t.readConn(connect)
	}
}

func (t *PrivateTCPServer) readConn(connection *net.TCPConn) {

	for {
		var stream []byte
		for {
			data := make([]byte, 1)
			n, err := io.ReadAtLeast(connection, data, 1)
			if err != nil {
				logger.Warn("read byte:%d,error:%s", n, err.Error())
				connection.Close()
				return
			}
			if data[0] == 10 {
				break
			} else {
				stream = append(stream, data...)
			}
		}
		var d *dm.DataMsg
		d = dm.NewDataMsg("job", types.MsgTypeTelnet, stream)
		d.SetMeta(t.Name, connection)
		t.Output.WritePipeNB(d)
	}
}

func (t *PrivateTCPServer) writeConn() {

	for {
		select {
		case data, ok := <-t.ReadPipe():
			if !ok {
				logger.Info("%s:Data Read error", t.Name)
				break
			}
			logger.Debug("%s:get msg from chan:%+v", t.Name, data)
			meta, ok := data.Meta(t.Name)
			if !ok {
				logger.Error("%s:wrong meta in datamsg:%+v", t.Name, data)
				break
			}
			connection := meta.(*net.TCPConn)
			//todo: need to verify if the data payload is []byte
			count, err := connection.Write(data.Payload.([]byte))
			if err != nil {
				logger.Warn("%s:conn write err:%s", t.Name, err.Error())
				connection.Close()
				return
			} else {
				logger.Info("%s:sent to network:%d byte", t.Name, count)
			}
		}
	}
}
