package datamsg

import (
	"library/logger"
)

type DataMsgPipe struct {
	Down chan *DataMsg
}

func NewDataMsgPipe(size int) *DataMsgPipe {
	var pipe = &DataMsgPipe{}
	downPipeLen := 1024
	if size >= 1 {
		downPipeLen = size
	}
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

func (p *DataMsgPipe) Length() int {
	return len(p.Down)
}
