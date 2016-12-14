package gatewayinner

import (
	. "gateway/manage"
	"library/logger"
	. "types"
)

//client req message, forward to server with userId
func HandleClientMessage(clientMeta *ConnMeta, msg *NetMsg) bool {

	if clientMeta.ForwardMeta == nil {
		Clients.Login(clientMeta, msg.Content)
	}

	logger.Info("gateway: recv from client ID <%s>, MSG <%s>", clientMeta.ID, msg.TypeString())
	serverMeta := clientMeta.ForwardMeta
	if serverMeta == nil {
		logger.Error("gateway: server meta for player %s is nil", clientMeta.ID)
	}

	msg.ObjectID = clientMeta.ObjID
	binary, err := msg.BinaryProto()
	if err != nil {
		logger.Error("encode net message to binary form error: %s", err.Error())
		return false
	}

	n, err := serverMeta.Send(binary)
	if err != nil {
		Servers.Logout(serverMeta)
		logger.Error("gateway: froward MSG <%16s> to server %s error: %s",
			msg.TypeString(), serverMeta.ID, err.Error())
		return false
	} else {
		logger.Info("gateway: forward MSG <%16s> to server %s success (%d/%d bytes)",
			msg.TypeString(), serverMeta.ID, msg.Size, n)
	}

	return true
}
