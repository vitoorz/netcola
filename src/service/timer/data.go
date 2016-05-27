package timer

import (
	dm "library/core/datamsg"
	"library/logger"
	"service"
)

func (t *timerType) DataEntry(msg *dm.DataMsg) (operate int, funCode int) {
	defer func() {
		if x := recover(); x != nil {
			logger.Error("timer panic: %v", x)
			logger.Stack()
		}
		operate = Continue
		funCode = service.FunPanic
	}()

	logger.Info("timer: data msg:%+v,payload:%s", msg, msg.Payload.([]byte))
	msg.Receiver = "tcpserver"
	//service.ServicePool.SendDown(msg)
	ok := t.Output.WriteDownChanNB(msg)
	if !ok {
		// channel full
		return Continue, service.FunDownChanFull
	}
	return Continue, service.FunOK
}
