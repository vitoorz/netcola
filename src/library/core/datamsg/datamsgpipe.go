package datamsg

import (
	"library/logger"
)

type DataMsgPipe struct {
	Up   chan *DataMsg
	Down chan *DataMsg
}

func NewDataMsgPipe(upSize, downSize int) *DataMsgPipe {
	var pipe = &DataMsgPipe{}
	upPipeLen := 1024
	downPipeLen := 1024
	if upSize >= 1 {
		upPipeLen = upSize
	}
	if downSize >= 1 {
		downPipeLen = downSize
	}
	pipe.Up = make(chan *DataMsg, upPipeLen)
	pipe.Down = make(chan *DataMsg, downPipeLen)
	return pipe
}

func (p *DataMsgPipe) ReadDownChan() chan *DataMsg {
	return p.Down
}

func (p *DataMsgPipe) WriteDownChanNB(msg *DataMsg) {
	select {
	case p.Down <- msg:
	default:
		logger.Warn("down chan overflow, may need to write to buffer")
	}

}

func (p *DataMsgPipe) ReadUpChan() chan *DataMsg {
	return p.Up
}

func (p *DataMsgPipe) WriteUpChanNB(msg *DataMsg) {
	select {
	case p.Up <- msg:
	default:
		logger.Warn("up chan overflow, may need to write to buffer")
	}

}

func (p *DataMsgPipe) Length() (int, int) {
	return len(p.Up), len(p.Down)
}
