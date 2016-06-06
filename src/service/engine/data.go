package engine

import (
	cm "library/core/controlmsg"
	dm "library/core/datamsg"
)

func (t *engineType) BackDoorEntry(msg *dm.DataMsg) (operate int, funCode int) {
	return cm.NextActionContinue, cm.ProcessStatOK
}

func (t *engineType) DataHandler(msg *dm.DataMsg) bool {
	return true
}
