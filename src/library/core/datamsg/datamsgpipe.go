package datamsg

import (
	"library/logger"
)

type DataMsgPipe struct {
	Pipe chan *DataMsg
}

func NewDataMsgPipe(size int) *DataMsgPipe {
	var pipe = &DataMsgPipe{}
	downPipeLen := 1024
	if size >= 1 {
		downPipeLen = size
	}
	pipe.Pipe = make(chan *DataMsg, downPipeLen)
	return pipe
}

func (t *DataMsgPipe) ReadPipe() chan *DataMsg {
	return t.Pipe
}

func (t *DataMsgPipe) WritePipeNoBlock(msg *DataMsg) bool {
	select {
	case t.Pipe <- msg:
		break
	default:
		logger.Info("Data chan full, so fork a routine wait here")
		go func() {
			t.Pipe <- msg
		}()
		return false
	}
	return true
}

func (t *DataMsgPipe) Length() int {
	return len(t.Pipe)
}
