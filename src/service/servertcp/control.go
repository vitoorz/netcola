package servertcp

import (
	"io"
	cm "library/core/controlmsg"
	dm "library/core/datamsg"
	//"library/idgen"
	"library/logger"
	"net"
	"service/serverhandle"
	"time"
	. "types"
)

func (t *serverTCP) Start(bus *dm.DataMsgPipe) bool {
	t.output = bus

	gatewayConnection := t.loginGateway()
	if gatewayConnection == nil {
		return false
	}

	logger.Info("%s:connect to gateway success: port %s", t.Name, t.port)
	go t.readConn(gatewayConnection)

	return true
}

func (t *serverTCP) Pause() bool {
	return true
}

func (t *serverTCP) Resume() bool {
	return true
}

func (t *serverTCP) Exit() bool {
	return true
}

func (t *serverTCP) ControlHandler(msg *cm.ControlMsg) (int, int) {
	return cm.NextActionContinue, cm.ProcessStatOK
}

func (t *serverTCP) readConn(connection net.Conn) {

	for {
		head := make([]byte, NetMsgHeadSize)
		n, err := io.ReadAtLeast(connection, head, NetMsgHeadSize)
		if err != nil {
			logger.Warn("%s:read byte:%d,error:%s", t.Name, n, err.Error())
			break
		}
		msg, err := NewNetMsgFromHead(head)
		if err != nil {
			logger.Warn("%s:decode msg head error:%s", t.Name, err.Error())
			break
		}
		_, err = io.ReadAtLeast(connection, msg.Content, int(msg.Size))
		if err != nil {
			logger.Warn("%s:read byte:%d,error:%s", t.Name, n, err.Error())
			break
		}

		d := dm.NewDataMsg(ServiceName, serverhandle.ServiceName, DataMsgFlagG2S, msg)
		d.SetMeta(t.Name, connection)
		t.output.WritePipeNoBlock(d)
	}

	connection.Close()
	logger.Error("lose server connection with gateway, try restart....")
	t.Start(t.output)
}

func (t *serverTCP) loginGateway() net.Conn {
	var (
		gatewayConnection net.Conn
		err               error
		retryTime         time.Duration
	)

	gatewayAddress := t.ip + ":" + t.port
	for {
		retryTime += 1
		logger.Info("begin to dial gateway server (%d times)....", retryTime)
		gatewayConnection, err = net.DialTimeout("tcp", gatewayAddress, time.Second*2)
		if err != nil {
			logger.Error("%s:net.ListenTCP error,%s", t.Name, err.Error())
		} else {
			break
		}
		time.Sleep(time.Second * retryTime)
	}

	req := &ServerLoginReq{ServerId: "0x84a5d1600002c001", ServerName: "xixihaha"}

	msg := &NetMsg{}
	msg.ObjectID = IdString(req.ServerId).ToObjectID()
	msg.SetPayLoad(MT_ServerLoginReq, req, NetMsgIdFlagServer)

	bin, err := msg.BinaryProto()
	if err != nil {
		logger.Error("encode headbeat message error %s", err.Error())
	}

	gatewayConnection.Write(bin)

	return gatewayConnection
}
