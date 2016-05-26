package job

import (
	dm "library/core/datamsg"
	"library/logger"
	"service"
)

const ServiceName = "job"

type jobType struct {
	service.Service
}

func NewJob(bus *dm.DataMsgPipe) *jobType {
	t := &jobType{}
	t.Service = *service.NewService(ServiceName)
	t.BUS = bus
	return t
}

func (t *jobType) Start() bool {
	logger.Info("job start running")
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

func (t *jobType) job() {
	for {
		select {
		case msg, _ := <-t.DataMsgPipe.Down:
			logger.Debug("job msg:%+v", msg)
		}
	}
}
