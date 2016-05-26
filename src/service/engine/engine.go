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
}

func NewEngine() *engineType {
	t := &engineType{}
	t.Service = *service.NewService(ServiceName)
	t.BUS = nil
	return t
}

func (t *engineType) engine(datapipe *dm.DataMsgPipe) {
	logger.Info("engine service running")

	var next int
	for {
		select {
		case msg, ok := <-t.Cmd:
			if !ok {
				logger.Info("Cmd Read error")
				break
			}
			next, ok = t.ControlEntry(msg)
			break
		case msg, ok := <-datapipe.ReadDownChan():
			if !ok {
				logger.Info("DownChan Read error")
				break
			}
			logger.Debug("engine recv data:%+v", msg)
			next, ok = t.DataEntry(msg)
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
