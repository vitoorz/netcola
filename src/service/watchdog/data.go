package watchdog

import (
	dm "library/core/datamsg"
)

func (t *watcherType) DataHandler(msg *dm.DataMsg) bool {
	return true
}
