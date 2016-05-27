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

func NewEngine() *engineType {
	t := &engineType{}
	t.Service = *service.NewService(ServiceName)
	t.BUS = dm.NewDataMsgPipe(0)
	return t
}

func (t *engineType) engine() {
	logger.Info("engine service running")

	var next, fun int = Continue, service.FunUnknown
	for {
		select {
		case msg, ok := <-t.Cmd:
			if !ok {
				logger.Info("Cmd Read error")
				break
			}
			next, ok = t.ControlEntry(msg)
			break
		case msg, ok := <-t.ReadDownChan():
			if !ok {
				logger.Info("DownChan Read error")
				break
			}
			next, fun = t.DataEntry(msg)
			break
		case msg, ok := <-t.BUS.ReadDownChan():
			if !ok {
				logger.Info("DownChan Read error")
				break
			}
			logger.Info("read bus chan msg:%+v", msg)
			next, fun = t.BusSchedule(msg)
			if fun == service.FunDownChanFull {
				logger.Warn("need do something when full")
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
