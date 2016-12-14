package serverhandle

import (
	gateway "gateway/manage"
	pb "github.com/golang/protobuf/proto"
	"library/logger"
	"server/play"
	. "types"
)

func HandleMessageFromGateway(payload *NetMsg) interface{} {
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
	}

	if ack == nil {
		logger.Info("ACK for MSG %s is nil, ID %s", opName, payload.ToIdString())
		return nil
	}

	pbAck, ok := ack.(pb.Message)
	if !ok {
		logger.Error("Ack for MSG %s: %s invaliid, can not marshal with protobuf", opName, ackCode.TypeString())
		return false
	}

	payload.SetPayLoad(ackCode, pbAck, msgFlag)
	logger.Info("ACK for MSG %16s: ID %s, ack code %16s", opName, payload.ToIdString(), payload.TypeString())
	return ack
}

//Message from client
func playHandler(payload *NetMsg) (interface{}, MsgType) {
	play.PlayFrame.SetFrameTime()

	handler, ok := play.NetMsgTypeHandler[payload.Code()]
	if !ok {
		logger.Info("PlayHandler get invalid payload type %d", payload.TypeString())
		return nil, MT_Blank
	}

	ack := handler.Handler(payload.ToIdString(), payload.Code(), payload.Content)

	return ack, handler.RetCode
}

//message from gateway
func serverHandler(payload *NetMsg) (interface{}, MsgType) {
	gateway.GateWayFrame.SetFrameTime()

	handler, ok := gateway.NetMsgTypeHandler[payload.Code()]
	if !ok {
		logger.Info("ServerHandler get invalid payload type %d", payload.TypeString())
		return nil, MT_Blank
	}

	ack := handler.Handler(payload.ToIdString(), payload.Code(), payload.Content)

	return ack, handler.RetCode
}
