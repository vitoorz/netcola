package engine

import (
	cm "library/core/controlmsg"
	dm "library/core/datamsg"
	"service"
)

func (t *engineType) BusSchedule(msg *dm.DataMsg) (operate int, funCode int) {
	service.ServicePool.SendData(msg)
	return cm.NextActionContinue, cm.ProcessStatOK
}
