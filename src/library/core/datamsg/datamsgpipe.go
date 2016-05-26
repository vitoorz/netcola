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

func (p *DataMsgPipe) WriteDownChanNB(msg *DataMsg) bool {
	select {
	case p.Down <- msg:
		break
	default:
		logger.Warn("down chan full")
		go func() {
			p.Down <- msg
		}()
		return false
	}
	return true
}

func (p *DataMsgPipe) Length() int {
	return len(p.Down)
}
