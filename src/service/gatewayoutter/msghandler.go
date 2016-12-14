package gatewayoutter

import (
	. "gateway/manage"
	pb "github.com/golang/protobuf/proto"
	"library/logger"
	. "types"
)

func forwardServerMessageToClient(serverMeta *ConnMeta, msg *NetMsg) bool {
	clientMeta, ok := Clients.GetMeta(msg.ToIdString())
	if !ok || clientMeta.Conn == nil {
		logger.Error("client connection for player %s is nil or serv", msg.ToIdString())
		return false
	}

	binary, err := msg.BinaryProtoToClient()
	if err != nil {
		logger.Error("encode net message to binary form error: %s", err.Error())
		return false
	}

	n, err := clientMeta.Conn.Write(binary)
	if err != nil {
		Clients.Logout(clientMeta)
		logger.Error("gateway: forward MSG <%16s> to player %s error: %s",
			msg.TypeString(), clientMeta.ID, err.Error())
		return false
	} else {
		logger.Info("gateway: foward MSG <%16s> to player %s success(%d/%d bytes)",
			msg.TypeString(), clientMeta.ID, msg.Size, n)
	}

	return true
}

//gateway receives message from server (protocol between server & gateway)
func handleServer2GatewayMessage(serverMeta *ConnMeta, msg *NetMsg) bool {
	handler, ok := NetMsgTypeHandler[msg.Code()]
	if !ok {
		logger.Warn("message from server: code %s do not handler", msg.TypeString())
		return false
	}

	if msg.Code() == MT_ServerLoginReq {
		Servers.Login(serverMeta, msg.Content)
	}

	opName := msg.TypeString()
	ack := handler.Handler(msg.ToIdString(), msg.Code(), msg.Content)
	if ack == nil {
		logger.Info("ACK for REQ %s is nil, ID %s", opName, msg.ToIdString())
		return true
	}

	pbAck, ok := ack.(pb.Message)
	if !ok {
		logger.Error("Ack payload for req %s invaliid, can not marshal with protobuf", msg.TypeString())
		return false
	}
	msg.SetPayLoad(handler.RetCode, pbAck, NetMsgIdFlagServer)

	logger.Info("ACK for REQ %16s: ID %s, ack code %16s",
		opName, msg.ToIdString(), msg.TypeString())

	bin, err := msg.BinaryProto()
	if err != nil {
		logger.Error("Ack payload for req %s protobuf encode error %s", opName, err.Error())
		return false
	}

	n, err := serverMeta.Send(bin)
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

func distributeBroadCastMessageToClients(serverMeta *ConnMeta, msg *NetMsg) bool {

	return true
}

//server ack/heartbeat message, forward to client without userId
func HandleMessageFromServer(serverMeta *ConnMeta, msg *NetMsg) bool {
	//server message for
	switch {
	case msg.HasFlag(NetMsgIdFlagClient):
		return forwardServerMessageToClient(serverMeta, msg)
	case msg.HasFlag(NetMsgIdFlagServer):
		return handleServer2GatewayMessage(serverMeta, msg)
	case msg.HasFlag(NetMsgIdFlagBroadCast):
		return distributeBroadCastMessageToClients(serverMeta, msg)
	}

	logger.Info("gateway: server message with invalid flag, code %d",
		msg.TypeString())

	return false
}
