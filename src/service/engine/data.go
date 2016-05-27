package engine

import (
	dm "library/core/datamsg"
	"service"
)

func (t *engineType) DataEntry(msg *dm.DataMsg) (operate int, funCode int) {
	//service.ServicePool.SendDown(msg)
	return Continue, service.FunOK
}
