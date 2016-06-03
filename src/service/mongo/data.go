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

	t.buffer.Append(msg)
	msg.Receiver = msg.Sender
	msg.Sender = t.Name
	t.output.WritePipeNB(msg)
	return Continue, service.FunOK
}
