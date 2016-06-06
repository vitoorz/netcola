package timer

import (
	cm "library/core/controlmsg"
	dm "library/core/datamsg"
)

func (t *timerType) Start(bus *dm.DataMsgPipe) bool {
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
