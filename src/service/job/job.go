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
	logger.Info("job service running")

	var next, fun int = Continue, service.FunUnknown
	for {
		select {
		case msg, ok := <-t.DataMsgPipe.Down:
			if !ok {
				logger.Info("Data Read error")
				break
			}
			next, fun = t.DataEntry(msg)
			if fun == service.FunDownChanFull {
				logger.Warn("need do something when full")
			}
			break
		case msg, ok := <-t.Cmd:
			if !ok {
				logger.Info("Cmd Read error")
				break
			}
			next, fun = t.ControlEntry(msg)
			break
		}

		switch next {
		case Break:
			break
		case Return:
			return
		case Continue:
		}
	}
	return

}
