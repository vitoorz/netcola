package servertcp

import (
	dm "library/core/datamsg"
	"library/logger"
	"net"
	. "types"
)

//send ack back
func (t *serverTCP) DataHandler(msg *dm.DataMsg) bool {
	if msg.MsgFlag != DataMsgFlagS2G {
		logger.Error("%s: recv invalid message type %d", t.Name, msg.MsgFlag)
		return false
	}

	meta, ok := msg.Meta(t.Name)
	if !ok {
		logger.Error("%s: wrong meta in datamsg:%+v", t.Name, msg)
		return false
	}

	connection, ok := meta.(net.Conn)
	if !ok {
		logger.Error("%s: wrong meta in datamsg(should be net.Conn):%+v", t.Name, msg)
		return false
	}

	//todo: need to verify if the data payload is []byte
	content, ok := msg.Payload.([]byte)
	if !ok {
		logger.Error("%s: payload is not in form []byte,", t.Name)
		return false
	}

	count, err := connection.Write(content)
	if err != nil {
		logger.Warn("%s:conn write err:%s", t.Name, err.Error())
		connection.Close()
		return false
	} else {
		logger.Info("%s:sent to network:%d byte", t.Name, count)
	}

	return true
}
