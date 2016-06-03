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

	t.Buffer.Append(msg)
	msg.Receiver = msg.Sender
	msg.Sender = t.Name
	t.output.WritePipeNB(msg)
	return Continue, service.FunOK
}

func (t *mongoType) Handle(msg *dm.DataMsg) bool {
	logger.Info("this is handle:%+v", msg)
	i, ok := msg.Meta(t.Name)
	if !ok {
		logger.Error("get Dirty interface error")
		return false
	}

	d, ok := i.(Dirty)
	if !ok {
		logger.Error("msg is not Dirty interface")
		return false
	}
	for ok = d.CRUD(t.session); !ok; {
		logger.Info("CRUD failed:%s", d.Inspect())
		ok = d.CRUD(t.session)
	}
	return true
}
