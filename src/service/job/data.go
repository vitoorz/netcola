package job

import (
	dm "library/core/datamsg"
	"library/logger"
	//"service"
)

func (t *jobType) DataEntry(msg *dm.DataMsg) (int, bool) {
	logger.Info("job: data msg:%+v,payload:%s", msg, msg.Payload.([]byte))
	msg.Receiver = "tcpserver"
	//service.ServicePool.SendDown(msg)
	t.BUS.Up <- msg
	return Continue, true
}
