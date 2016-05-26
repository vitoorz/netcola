package job

import (
	dm "library/core/datamsg"
	"library/logger"
	//"service"
)

func (t *jobType) DataEntry(msg *dm.DataMsg) (int, bool) {
	defer func() (int, bool) {
		if x := recover(); x != nil {
			logger.Error("job panic: %v", x)
			logger.Stack()
		}
		return Continue, false
	}()
	logger.Info("job: data msg:%+v,payload:%s", msg, msg.Payload.([]byte))
	msg.Receiver = "tcpserver"
	//service.ServicePool.SendDown(msg)
	t.BUS.Up <- msg
	return Continue, true
}
