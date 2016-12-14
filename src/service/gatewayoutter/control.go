package gatewayoutter

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

func (t *gatewayOutter) Start(bus *dm.DataMsgPipe) bool {
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
	go t.gatewayClient()
	return true
}

func (t *gatewayOutter) Pause() bool {
	return true
}

func (t *gatewayOutter) Resume() bool {
	return true
}

func (t *gatewayOutter) Exit() bool {
	return true
}

func (t *gatewayOutter) ControlHandler(msg *cm.ControlMsg) (int, int) {
	return cm.NextActionContinue, cm.ProcessStatOK
}

func (t *gatewayOutter) gatewayClient() {

	for {
		connect, err := t.listener.AcceptTCP()
		if err != nil {
			logger.Error("%s:listener.AcceptTCP error,%s", t.Name, err.Error())
			time.Sleep(time.Second * 2)
			connect.Close()
			continue
		}
		go t.gatewayClientConn(connect)
	}
}

func (t *gatewayOutter) gatewayClientConn(connection *net.TCPConn) {
	clientMeta := NewClientMeta(connection)

	for clientMeta.Conn != nil {
		head := make([]byte, NetMsgHeadNoIdSize)
		_, err := io.ReadAtLeast(connection, head, NetMsgHeadNoIdSize)
		if err != nil {
			logger.Error("%s: read player %s req head[NOID] error: %s",
				t.Name, clientMeta.ID, err.Error())
			Clients.Logout(clientMeta)
			break
		}

		msg, err := NewNetMsgFromHeadNoId(head)
		if err != nil {
			logger.Error("%s:decode player %s req head(NOID) error: %s",
				t.Name, clientMeta.ID, err.Error())
			Clients.Logout(clientMeta)
			break
		}

		msg.MsgType = AddMsgFlag(msg.MsgType, NetMsgIdFlagClient)

		_, err = io.ReadAtLeast(connection, msg.Content, int(msg.Size))
		if err != nil {
			logger.Error("%s: read player %s req [%s] payload error",
				t.Name, clientMeta.ID, msg.TypeString())
			Clients.Logout(clientMeta)
			break
		}

		d := dm.NewDataMsg(ServiceName, "gatewayinner", DataMsgFlagC2G, msg)
		d.SetMeta(t.Name, clientMeta)
		t.output.WritePipeNoBlock(d)
	}

	logger.Warn("%s: player %s connection exit", t.Name, clientMeta.ID)
}
