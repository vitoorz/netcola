package job

import (
	cm "library/core/controlmsg"
	dm "library/core/datamsg"
	"library/logger"
)

func (t *jobType) Start(bus *dm.DataMsgPipe) bool {
	logger.Info("%s:start running", t.Name)
	t.Output = bus
	return true
}

func (t *jobType) Pause() bool {
	return true
}

func (t *jobType) Resume() bool {
	return true
}

func (t *jobType) Exit() bool {
	return true
}

func (t *jobType) ControlHandler(msg *cm.ControlMsg) (int, int) {
	return cm.NextActionContinue, cm.ProcessStatOK
}
