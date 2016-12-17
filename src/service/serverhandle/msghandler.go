package serverhandle

import (
	"game/gateway/gm"
	"game/server/play"
	pb "github.com/golang/protobuf/proto"
	"library/logger"
	"netmsghandle/cs"
	"netmsghandle/gs"
	. "types"
)

func HandleMessageFromGateway(payload *NetMsg) (interface{}, bool) {
	var (
		ack     interface{} = nil
		ackCode MsgType
	)

	opName := payload.TypeString()
	msgFlag := uint32(NetMsgIdFlagClient)

	switch {
	case payload.HasFlag(NetMsgIdFlagServer):
		ack, ackCode = serverHandler(payload)
		msgFlag = NetMsgIdFlagServer
	case payload.HasFlag(NetMsgIdFlagClient):
		ack, ackCode = playHandler(payload)
		msgFlag = NetMsgIdFlagClient
	case payload.HasFlag(NetMsgIdFlagBroadCast):
		ack, ackCode = serverHandler(payload)
		msgFlag = NetMsgIdFlagBroadCast
	}

	if ack == nil {
		logger.Info("ACK for MSG %s is nil, ID %s", opName, payload.ToIdString())
		return nil, true
	}

	pbAck, ok := ack.(pb.Message)
	if !ok {
		logger.Error("Ack for MSG %s: %s invaliid, can not marshal with protobuf", opName, ackCode.TypeString())
		return nil, false
	}

	payload.SetPayLoad(ackCode, pbAck, msgFlag)

	logger.Info("ACK for MSG %16s: ID %s, ack code %16s", opName, payload.ToIdString(), payload.TypeString())
	return ack, true
}

//Message from client
func playHandler(payload *NetMsg) (interface{}, MsgType) {
	play.PlayFrame.SetFrameTime()

	handler, ok := cs.NetMsgTypeHandler[payload.Code()]
	if !ok {
		logger.Info("PlayHandler get invalid payload type %s", payload.TypeString())
		return nil, MT_Blank
	}

	ack := handler.Handler(payload.ToIdString(), payload.Code(), payload.Content)

	return ack, handler.RetCode
}

//message from gateway
func serverHandler(payload *NetMsg) (interface{}, MsgType) {
	gm.GatewayFrame.SetFrameTime()

	handler, ok := gs.NetMsgTypeHandler[payload.Code()]
	if !ok {
		logger.Info("ServerHandler get invalid payload type %d", payload.TypeString())
		return nil, MT_Blank
	}

	ack := handler.Handler(payload.ToIdString(), payload.Code(), payload.Content)

	return ack, handler.RetCode
}
