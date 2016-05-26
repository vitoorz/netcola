package job

import (
	dm "library/core/datamsg"
	"service"
)

func (t *jobType) DataEntry(msg *dm.DataMsg) (int, bool) {
	msg.Receiver = ""
	service.ServicePool.SendDown(msg)
	return Continue, true
}
