package privatetcp

import (
	"io"
	cm "library/core/controlmsg"
	dm "library/core/datamsg"
	"library/logger"
	"net"
	"time"
	"types"
)

func (t *privateTCPServer) Start(bus *dm.DataMsgPipe) bool {
	t.output = bus
	tcpAddr, err := net.ResolveTCPAddr("tcp", t.ip+":"+t.port)
	if err != nil {
		logger.Error("%s:net.ResolveTCPAddr error,%s", t.Name, err.Error())
		return false
	}

	t.listener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		logger.Error("%s:net.ListenTCP error,%s", t.Name, err.Error())
		return false
	}

	logger.Info("%s:listening port:%s", t.Name, t.port)
	go t.serve()
	return true
}

func (t *privateTCPServer) Pause() bool {
	return true
}

func (t *privateTCPServer) Resume() bool {
	return true
}

func (t *privateTCPServer) Exit() bool {
	return true
}

func (t *privateTCPServer) ControlHandler(msg *cm.ControlMsg) (int, int) {
	return cm.NextActionContinue, cm.ProcessStatOK
}

func (t *privateTCPServer) serve() {

	for {
		connect, err := t.listener.AcceptTCP()
		if err != nil {
			logger.Error("%s:listener.AcceptTCP error,%s", t.Name, err.Error())
			time.Sleep(time.Second * 2)
			connect.Close()
			continue
		}
		go t.readConn(connect)
	}
}

func (t *privateTCPServer) readConn(connection *net.TCPConn) {

	for {
		var stream []byte
		for {
			data := make([]byte, 1)
			n, err := io.ReadAtLeast(connection, data, 1)
			if err != nil {
				logger.Warn("%s:read byte:%d,error:%s", t.Name, n, err.Error())
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
		d = dm.NewDataMsg("tcpserver", "job", types.MsgTypeTelnet, stream)
		d.SetMeta(t.Name, connection)
		t.output.WritePipeNB(d)
	}
}
