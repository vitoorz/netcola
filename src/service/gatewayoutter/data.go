package gatewayoutter

import (
	. "gateway/manage"
	dm "library/core/datamsg"
	"library/logger"
	. "types"
)

//message from server
func (t *gatewayOutter) DataHandler(msg *dm.DataMsg) bool {
	if msg.MsgType != Inner_MsgTypeG2C {
		logger.Error("%s: recv invalid message type %d", t.Name, msg.MsgType)
		return false
	}

	meta, ok := msg.Meta("gatewayinner")
	if !ok {
		logger.Error("%s:wrong meta in datamsg:%+v", t.Name, msg)
		return false
	}

	serverMeta, ok := meta.(*ServerConnectionMeta)
	if !ok {
		logger.Error("%s:wrong meta in datamsg(should be *ClientConnectionMeta):%+v", t.Name, meta)
		return false
	}

	netMsg, ok := msg.Payload.(*NetMsg)
	if !ok {
		logger.Error("%s:wrong payload in datamsg(should be *NetMsg):%+v", t.Name, msg.Payload)
		return false
	}

	return HandleG2CMsg(serverMeta, netMsg)
}
