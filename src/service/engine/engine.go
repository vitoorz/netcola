package engine

import (
	dm "library/core/datamsg"
	"library/logger"
	"service"
)

const ServiceName = "engine"

const (
	Break = iota
	Continue
	Return
)

type engineType struct {
	service.Service
	BUS *dm.DataMsgPipe
}

func NewEngine(name string) *engineType {
	t := &engineType{}
	t.Service = *service.NewService(ServiceName)
	t.BUS = dm.NewDataMsgPipe(0)
	t.Name = name
	return t
}

func (t *engineType) engine() {
	logger.Info("%s:service running", t.Name)

	var next, fun int = Continue, service.FunUnknown
	for {
		select {
		case msg, ok := <-t.Cmd:
			if !ok {
				logger.Info("%s:Cmd Read error", t.Name)
				break
			}
			next, ok = t.controlEntry(msg)
			break
		case msg, ok := <-t.ReadPipe():
			if !ok {
				logger.Info("%s:DataPipe Read error", t.Name)
				break
			}
			next, fun = t.dataEntry(msg)
			break
		case msg, ok := <-t.BUS.ReadPipe():
			if !ok {
				logger.Info("%s:DataPipe Read error", t.Name)
				break
			}
			logger.Debug("%s:read bus chan msg:%+v", t.Name, msg)
			next, fun = t.BusSchedule(msg)
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
