package timer

import (
	dm "library/core/datamsg"
	"library/logger"
	"service"
)

const ServiceName = "timer"

const (
	Break = iota
	Continue
	Return
)

type timerType struct {
	service.Service
	Output *dm.DataMsgPipe
}

func NewTimer() *timerType {
	t := &timerType{}
	t.Service = *service.NewService(ServiceName)
	t.Output = nil
	return t
}

func (t *timerType) job() {
	logger.Info("timer service running")

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
