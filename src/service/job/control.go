package job

import (
	dm "library/core/datamsg"
	"library/logger"
)

func (t *jobType) Start(name string, bus *dm.DataMsgPipe) bool {
	logger.Info("job start running")
	t.Name = name
	t.BUS = bus
	go t.job()
	return true
}

func (t *jobType) Pause() bool {
	return true
}

func (t *jobType) Resume() bool {
	return true
}

func (t *jobType) Exit() bool {
	return true
}
