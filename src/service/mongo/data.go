package mongo

import (
//"time"
)

import (
	dm "library/core/datamsg"
	"library/logger"
	"service"
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
	d, ok := m.(service.Page)
	if !ok {
		logger.Error("msg is not Page")
		return Continue, service.FunUnknown
	}
	t.buffer.Append(&d)
	msg.Receiver = msg.Sender
	msg.Sender = t.Name
	t.output.WritePipeNB(msg)
	return Continue, service.FunOK
}
