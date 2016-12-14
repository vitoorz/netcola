package gatewayinner

import (
	. "gateway/manage"
	dm "library/core/datamsg"
	"library/logger"
	. "types"
)

//messages receive from client
func (t *gatewayInner) DataHandler(msg *dm.DataMsg) bool {
	if msg.MsgFlag != DataMsgFlagC2G {
		logger.Error("%s: recv invalid message type %d", t.Name, msg.MsgFlag)
		return false
	}

	meta, ok := msg.Meta("gatewayoutter")
	if !ok {
		logger.Error("%s:wrong meta in datamsg:%+v", t.Name, msg)
		return false
	}

	clientMeta, ok := meta.(*ConnMeta)
	if !ok {
		logger.Error("%s:wrong meta in datamsg(should be *ClientConnectionMeta):%+v", t.Name, meta)
		return false
	}

	netMsg, ok := msg.Payload.(*NetMsg)
	if !ok {
		logger.Error("%s:wrong payload in datamsg(should be *NetMsg):%+v", t.Name, msg.Payload)
		return false
	}

	return HandleClientMessage(clientMeta, netMsg)
}
