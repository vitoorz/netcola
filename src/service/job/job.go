package job

import (
	"library/logger"
	"service"
)

const ServiceName = "job"

const (
	Break = iota
	Continue
	Return
)

type jobType struct {
	service.Service
}

func NewJob() *jobType {
	t := &jobType{}
	t.Service = *service.NewService(ServiceName)
	t.BUS = nil
	return t
}

func (t *jobType) job() {
	for {
		select {
		case msg, _ := <-t.DataMsgPipe.Down:
			logger.Debug("job msg:%+v", msg)
		}
	}
}
