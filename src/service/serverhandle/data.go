package serverhandle

import (
	dm "library/core/datamsg"
	"library/logger"
	. "types"
)

func (t *serverHandle) DataHandler(msg *dm.DataMsg) bool {

	if msg.MsgFlag != DataMsgFlagG2S {
		logger.Info("%s: recv invalid message type %d", t.Name, msg.MsgFlag)
		return false
	}

	payload, ok := msg.Payload.(*NetMsg)
	if !ok {
		logger.Info("%s: recv invalid message type %d, payload error", t.Name, msg.MsgFlag)
		return false
	}

	ack, _ := HandleMessageFromGateway(payload)
	msg.Receiver = msg.Sender
	msg.Sender = ServiceName
	if ack != nil && msg.Receiver != dm.NoReceiver {
		msg.MsgFlag = DataMsgFlagS2G
		t.Output.WritePipeNoBlock(msg)
	}

	return true
}
