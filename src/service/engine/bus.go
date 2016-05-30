package engine

import (
	dm "library/core/datamsg"
	"service"
)

func (t *engineType) BusSchedule(msg *dm.DataMsg) (operate int, funCode int) {
	service.ServicePool.SendData(msg)
	return Continue, service.FunOK
}
