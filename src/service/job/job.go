package job

import (
	dm "library/core/datamsg"
	"library/logger"
	"service"
)

const ServiceName = "job"

type jobType struct {
	*service.Service
	Output *dm.DataMsgPipe
}

func NewJob(name string) *jobType {
	t := &jobType{}
	t.Service = service.NewService(ServiceName)
	t.Output = nil
	t.Name = name
	t.Instance = t
	return t
}
