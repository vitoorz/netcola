package engine

import (
	dm "library/core/datamsg"
	"service"
)

func (t *engineType) DataEntry(msg *dm.DataMsg) (int, bool) {
	service.ServicePool.SendDown(msg)
	return Continue, true
}
