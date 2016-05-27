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

func (t *DataMsgPipe) ReadDownChan() chan *DataMsg {
	return t.Down
}

func (t *DataMsgPipe) WriteDownChanNB(msg *DataMsg) bool {
	select {
	case t.Down <- msg:
		break
	default:
		logger.Warn("down chan full")
		go func() {
			t.Down <- msg
		}()
		return false
	}
	return true
}

func (t *DataMsgPipe) Length() int {
	return len(t.Down)
}
