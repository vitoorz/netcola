package timer

import (
	cm "library/core/controlmsg"
	dm "library/core/datamsg"
	"library/logger"
)

func (t *timerType) Start(bus *dm.DataMsgPipe) bool {
	logger.Info("%s:start running", t.Name)
	t.output = bus
	return true
}

func (t *timerType) Pause() bool {
	return true
}

func (t *timerType) Resume() bool {
	return true
}

func (t *timerType) Exit() bool {
	return true
}

func (t *timerType) ControlHandler(msg *cm.ControlMsg) (int, int) {
	return cm.NextActionContinue, cm.ProcessStatOK
}
