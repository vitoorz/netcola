package netmsg

import (
	"library/logger"
)

type NetMsgPipe struct {
	send chan *NetMsg
	recv chan *NetMsg
}

func NewNetMsgPipe(sendSize, recvSize int) *NetMsgPipe {
	var pipe = &NetMsgPipe{}
	sendPipeLen := 1024
	recvPipeLen := 1024
	if sendSize >= 1 {
		sendPipeLen = sendSize
	}
	if recvSize >= 1 {
		recvPipeLen = recvSize
	}
	pipe.send = make(chan *NetMsg, sendPipeLen)
	pipe.recv = make(chan *NetMsg, recvPipeLen)
	return pipe
}

func (p *NetMsgPipe) ReadRecvChan() *chan *NetMsg {
	return &p.recv
}

func (p *NetMsgPipe) WriteRecvChan(msg *NetMsg) {
	select {
	case p.recv <- msg:
	default:
		logger.Warn("Recv chan overflow, need write to disk queue")
	}

}

func (p *NetMsgPipe) ReadSendChan() *chan *NetMsg {
	return &p.send
}

func (p *NetMsgPipe) WriteSendChan(msg *NetMsg) {
	select {
	case p.send <- msg:
	default:
		logger.Warn("WriteSendChan overflow, need write to disk queue")
	}

}

func (p *NetMsgPipe) Length() (int, int) {
	return len(p.send), len(p.recv)
}
