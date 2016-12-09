package gatewayinner

import (
	. "gateway/manage"
	"io"
	cm "library/core/controlmsg"
	dm "library/core/datamsg"
	"library/logger"
	"net"
	"time"
	. "types"
)

func (t *gatewayInner) Start(bus *dm.DataMsgPipe) bool {
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
	go t.gatewayServer()
	return true
}

func (t *gatewayInner) Pause() bool {
	return true
}

func (t *gatewayInner) Resume() bool {
	return true
}

func (t *gatewayInner) Exit() bool {
	return true
}

func (t *gatewayInner) ControlHandler(msg *cm.ControlMsg) (int, int) {
	return cm.NextActionContinue, cm.ProcessStatOK
}

func (t *gatewayInner) gatewayServer() {

	for {
		connect, err := t.listener.AcceptTCP()
		if err != nil {
			logger.Error("%s:listener.AcceptTCP error,%s", t.Name, err.Error())
			time.Sleep(time.Second * 2)
			connect.Close()
			continue
		}
		go t.gatewayServerConn(connect)
	}
}

//get one message from server connection
func (t *gatewayInner) gatewayServerConn(connection *net.TCPConn) {
	serverMeta := &ServerConnectionMeta{ServerConn: connection}

	for serverMeta.ServerConn != nil {
		head := make([]byte, NetMsgHeadSize)
		_, err := io.ReadAtLeast(connection, head, NetMsgHeadSize)
		if err != nil {
			logger.Error("%s: read server %s ack head(+ID) error: %s",
				t.Name, serverMeta.ServerId.ToIdString(), err.Error())
			ServerLogout(serverMeta)
			break
		}

		msg, err := NewNetMsgFromHead(head)
		if err != nil {
			logger.Error("%s: decode server %s ack head(+ID) error: %s",
				t.Name, serverMeta.ServerId.ToIdString(), err.Error())
			ServerLogout(serverMeta)
			break
		}

		_, err = io.ReadAtLeast(connection, msg.Content, int(msg.Size))
		if err != nil {
			logger.Error("%s: read server %s ack payload error: %s",
				t.Name, serverMeta.ServerId.ToIdString(), err.Error())
			ServerLogout(serverMeta)
			break
		}

		logger.Info("%s: get server message <%16s>", t.Name, msg.OpCode.ToString())
		d := dm.NewDataMsg(ServiceName, "gatewayoutter", Inner_MsgTypeG2C, msg)
		d.SetMeta(t.Name, serverMeta)
		t.output.WritePipeNoBlock(d)
	}

	logger.Warn("%s: server %s exit", t.Name, serverMeta.ServerId.ToIdString())
}
