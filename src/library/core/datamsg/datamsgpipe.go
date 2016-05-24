package datamsg

import (
	"library/logger"
)

type DataMsgPipe struct {
	send chan *DataMsg
	recv chan *DataMsg
}

func NewNetMsgPipe(sendSize, recvSize int) *DataMsgPipe {
	var pipe = &DataMsgPipe{}
	sendPipeLen := 1024
	recvPipeLen := 1024
	if sendSize >= 1 {
		sendPipeLen = sendSize
	}
	if recvSize >= 1 {
		recvPipeLen = recvSize
	}
	pipe.send = make(chan *DataMsg, sendPipeLen)
	pipe.recv = make(chan *DataMsg, recvPipeLen)
	return pipe
}

func (p *DataMsgPipe) ReadRecvChan() *chan *DataMsg {
	return &p.recv
}

func (p *DataMsgPipe) WriteRecvChan(msg *DataMsg) {
	select {
	case p.recv <- msg:
	default:
		logger.Warn("Recv chan overflow, need write to disk queue")
	}

}

func (p *DataMsgPipe) ReadSendChan() *chan *DataMsg {
	return &p.send
}

func (p *DataMsgPipe) WriteSendChan(msg *DataMsg) {
	select {
	case p.send <- msg:
	default:
		logger.Warn("WriteSendChan overflow, need write to disk queue")
	}

}

func (p *DataMsgPipe) Length() (int, int) {
	return len(p.send), len(p.recv)
}
