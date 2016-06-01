package mongo

import (
//"time"
)

import (
	dm "library/core/datamsg"
	"library/logger"
	"service"
	ts "types/service"
)

func (t *mongoType) dataEntry(msg *dm.DataMsg) (operate int, funCode int) {
	defer func() {
		if x := recover(); x != nil {
			logger.Error("%s:mongo panic: %v", t.Name, x)
			logger.Stack()
		}
		operate = Continue
		funCode = service.FunPanic
	}()

	logger.Info("%s:get data msg:%d,payload:%v", t.Name, msg.MsgType, msg.Payload.([]byte))
	m, ok := msg.Meta(t.Name)
	if !ok {
		//todo: do more
		return Continue, service.FunOK
	}
	d, ok := m.(ts.Dirty)
	if !ok {
		logger.Error("msg is not Dirty")
		return Continue, service.FunUnknown
	}
	t.dirtyPool.addDirty(&d)
	msg.Receiver = msg.Sender
	msg.Sender = t.Name
	ok = t.output.WritePipeNB(msg)
	if !ok {
		// channel full
		return Continue, service.FunDataPipeFull
	}

	return Continue, service.FunOK
}
