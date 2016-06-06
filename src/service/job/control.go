package job

import (
	cm "library/core/controlmsg"
	dm "library/core/datamsg"
)

func (t *jobType) Start(bus *dm.DataMsgPipe) bool {
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
