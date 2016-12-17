package gatewayinner

import (
	"game/gateway/gm"
	"library/logger"
	. "types"
)

//client req message, forward to server with userId
func GatewayOnClientMessage(clientMeta *gm.ConnMeta, msg *NetMsg) bool {

	if clientMeta.ForwardMeta == nil {
		gm.Clients.Login(clientMeta, msg.Content)
	}

	logger.Info("gateway: MSG <%16s> from player <%s> received", clientMeta.ID, msg.TypeString())
	serverMeta := clientMeta.ForwardMeta
	if serverMeta == nil {
		logger.Error("gateway: player %s, server not distributed", clientMeta.ID)
	}

	msg.ObjectID = clientMeta.ObjID

	binary, err := msg.BinaryProto()
	if err != nil {
		logger.Error("gateway: player %s, payload marshal error: %s", clientMeta.ID, err.Error())
		return false
	}

	return serverMeta.CsToServer(msg.TypeString(), binary)
}
