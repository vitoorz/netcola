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

func NewEngine(bus *dm.DataMsgPipe) *engineType {
	t := &engineType{}
	t.Service = *service.NewService(ServiceName)
	t.BUS = bus
	return t
}

func (t *engineType) engine(datapipe *dm.DataMsgPipe) (err interface{}) {
	//defer func() {
	//	if x := recover(); x != nil {
	//		logger.Error("Engine job panic: %v", x)
	//		logger.Stack()
	//	}
	//}()
	err = nil
	logger.Info("engine cycle start")
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
			return err
		case Continue:
		}
	}
	return err
}
