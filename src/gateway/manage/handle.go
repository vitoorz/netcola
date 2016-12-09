package manage

import (
	"library/logger"
	. "types"
)

//client req message, forward to server with userId
func HandleC2GMsg(clientMeta *ClientConnectionMeta, msg *NetMsg) bool {
	if clientMeta.ServerMeta == nil {
		clientLogin(clientMeta, msg)
	}

	serverMeta := clientMeta.ServerMeta
	if serverMeta == nil {
		return false
	}

	msg.UserId = clientMeta.UserId
	binary, err := msg.BinaryProto()
	if err != nil {
		logger.Error("encode net message to binary form error: %s", err.Error())
		return false
	}

	n, err := serverMeta.ServerConn.Write(binary)
	if err != nil {
		ServerLogout(serverMeta)
		logger.Error("gateway froward REQ <%16s> to server %s error: %s",
			msg.OpCode.ToString(), serverMeta.ServerId.ToIdString(), err.Error())
		return false
	} else {
		logger.Info("gateway forward REQ <%16s> to server %s success (%d bytes)",
			msg.OpCode.ToString(), serverMeta.ServerId.ToIdString(), n)
	}

	return true
}

//server ack/heartbeat message, forward to client without userId
func HandleG2CMsg(serverMeta *ServerConnectionMeta, msg *NetMsg) bool {
	if msg.OpCode == MT_HeatBeat {
		serverLogin(serverMeta, msg)
		return true
	}

	clientMeta, ok := getClientMeta(msg.UserId)
	if !ok || clientMeta.ClientConn == nil {
		logger.Error("client connection for player %d is nil", msg.UserId)
		return false
	}

	binary, err := msg.BinaryProtoNoId()
	if err != nil {
		logger.Error("encode net message to binary form error: %s", err.Error())
		return false
	}

	n, err := clientMeta.ClientConn.Write(binary)
	if err != nil {
		ClientLogout(clientMeta)
		logger.Error("gateway forward ACK <%16s> to player %s error: %s",
			msg.OpCode.ToString(), clientMeta.UserId.ToIdString(), err.Error())
		return false
	} else {
		logger.Info("gateway foward ACK <%16s> to player %s (%d bytes)",
			msg.OpCode.ToString(), clientMeta.UserId.ToIdString(), n)
	}

	return true
}
