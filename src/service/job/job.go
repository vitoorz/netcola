package job

import (
	dm "library/core/datamsg"
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
	Output *dm.DataMsgPipe
}

func NewJob(name string) *jobType {
	t := &jobType{}
	t.Service = *service.NewService(ServiceName)
	t.Output = nil
	t.Name = name
	return t
}

func (t *jobType) job() {
	logger.Info("%s:service running", t.Name)
	var next, fun int = Continue, service.FunUnknown
	for {
		select {
		case msg, ok := <-t.Cmd:
			if !ok {
				logger.Info("%s:Cmd Read error", t.Name)
				break
			}
			next, fun = t.controlEntry(msg)
			break
		case msg, ok := <-t.ReadPipe():
			if !ok {
				logger.Info("%s:Data Read error", t.Name)
				break
			}
			next, fun = t.dataEntry(msg)
			if fun == service.FunDataPipeFull {
				logger.Warn("%s:need do something when full", t.Name)
			}
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
