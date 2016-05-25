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

func (t *jobType) Init() bool {
	logger.Info("job init")
	return true
}

func (t *jobType) Start() bool {
	logger.Info("job start running")
	//go t.engine(t.BUS)
	return true
}

func (t *jobType) Pause() bool {
	return true
}

func (t *jobType) Exit() bool {
	return true
}

func (t *jobType) Kill() bool {
	return true
}
