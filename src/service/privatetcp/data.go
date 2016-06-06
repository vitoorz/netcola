package privatetcp

import (
	dm "library/core/datamsg"
	"library/logger"
	"net"
)

func (t *privateTCPServer) DataHandler(msg *dm.DataMsg) bool {
	meta, ok := msg.Meta(t.Name)
	if !ok {
		logger.Error("%s:wrong meta in datamsg:%+v", t.Name, msg)
		return false
	}
	connection := meta.(*net.TCPConn)
	//todo: need to verify if the data payload is []byte
	count, err := connection.Write(msg.Payload.([]byte))
	if err != nil {
		logger.Warn("%s:conn write err:%s", t.Name, err.Error())
		connection.Close()
		return false
	} else {
		logger.Info("%s:sent to network:%d byte", t.Name, count)
	}
	return true
}
