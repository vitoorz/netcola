package job

import (
	dm "library/core/datamsg"
	"library/logger"
	"service"
)

func (t *jobType) DataEntry(msg *dm.DataMsg) (operate int, funCode int) {
	defer func() {
		if x := recover(); x != nil {
			logger.Error("job panic: %v", x)
			logger.Stack()
		}
		operate = Continue
		funCode = service.FunPanic
	}()

	logger.Info("job: data msg:%+v,payload:%s", msg, msg.Payload.([]byte))

	switch msg.MsgType {
	case 1:
		msg.Receiver = "tcpserver"
	case 2:
		msg.Receiver = "timer"
	}

	//service.ServicePool.SendDown(msg)
	ok := t.Output.WritePipeNB(msg)
	if !ok {
		// channel full
		return Continue, service.FunDownChanFull
	}
	return Continue, service.FunOK

}
