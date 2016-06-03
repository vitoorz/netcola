package engine

import (
	dm "library/core/datamsg"
	"service"
)

func (t *engineType) BackDoorEntry(msg *dm.DataMsg) (operate int, funCode int) {
	//service.ServicePool.SendDown(msg)
	return Continue, service.FunOK
}

func (t *engineType) DataHandler(msg *dm.DataMsg) bool {
	return true
}
