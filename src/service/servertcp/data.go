package servertcp

import (
	"game/server/play"
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

	sendLists := play.AsyncSender.GetAsyncNetMessages()
	nm, ok := msg.Payload.(*NetMsg)
	if !ok {
		logger.Error("%s: payload %s not type of NetMsg", t.Name)
	} else {
		sendLists = append(sendLists, nm)
	}

	for _, nm = range sendLists {
		content, err := nm.BinaryProto()
		if err != nil {
			logger.Warn("%s:netmsg poylad %s marshal error: %s", nm.TypeString(), err.Error())
			continue
		}

		count, err := connection.Write(content)
		if err != nil {
			logger.Warn("%s:conn write err:%s", t.Name, err.Error())
			connection.Close()
			return false
		} else {
			logger.Info("%s:sent to network: %s (%d byte)", t.Name, nm.TypeString(), count)
		}
	}

	return true
}
