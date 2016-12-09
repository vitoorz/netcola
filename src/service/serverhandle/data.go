package serverhandle

import (
	pb "github.com/golang/protobuf/proto"
	dm "library/core/datamsg"
	"library/logger"
	"play"
	. "types"
)

func serverHandler(payload *NetMsg) bool {

	return true
}

func (t *serverHandle) playHandler(payload *NetMsg) bool {
	play.PlayFrame.SetFrameTime()

	handler, ok := play.NetMsgTypeHandler[payload.OpCode]
	if !ok {
		logger.Info("recv invalid payload type %d", payload.OpCode)
		return false
	}

	playerId := payload.UserId.ToIdString()

	ack := handler.Handler(playerId, payload.OpCode, payload.Content)
	if ack != nil {
		pbAck, ok := ack.(pb.Message)
		if !ok {
			logger.Error("Ack payload for req %d invaliid, can not marshal with protobuf", payload.OpCode)
			return false
		}
		payload.OpCode = handler.RetCode
		payload.SetPayLoad(pbAck)
		logger.Info("%s: ACK for REQ %16s: ID %s, ack code %16s",
			t.Name, handler.OpCode.ToString(), playerId, handler.RetCode.ToString())
	} else {
		logger.Info("%s: ACK for REQ %d is nil, ID %s", t.Name, handler.OpCode.ToString(), playerId)
	}

	return true
}

func (t *serverHandle) DataHandler(msg *dm.DataMsg) bool {

	if msg.MsgType != Inner_MsgTypeG2S {
		logger.Info("%s: recv invalid message type %d", t.Name, msg.MsgType)
		return false
	}

	payload, ok := msg.Payload.(*NetMsg)
	if !ok {
		logger.Info("%s: recv invalid message type %d, payload error", t.Name, msg.MsgType)
		return false
	}

	if payload.OpCode == MT_HeatBeat {
		return serverHandler(payload)
	}

	if ok = t.playHandler(payload); ok {
		msg.Payload, _ = payload.BinaryProto()
		msg.Receiver = msg.Sender
		msg.Sender = ServiceName
		msg.MsgType = Inner_MsgTypeS2G
		if msg.Receiver != dm.NoReceiver {
			t.Output.WritePipeNoBlock(msg)
		}
	}

	return true
}
