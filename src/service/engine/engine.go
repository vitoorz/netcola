package engine

import (
	cm "library/core/controlmsg"
	dm "library/core/datamsg"
	"library/logger"
	"service"
)

const ServiceName = "engine"

type engineType struct {
	service.Service
	BUS *dm.DataMsgPipe
}

func NewEngine(name string) *engineType {
	t := &engineType{}
	t.Service = *service.NewService(ServiceName)
	t.BUS = dm.NewDataMsgPipe(0)
	t.Name = name
	t.SelfDrive = true
	t.Instance = t
	return t
}

func (t *engineType) engine() {
	var next, fun int = cm.NextActionContinue, cm.ProcessStatUnknown
	for {
		select {
		case msg, ok := <-t.Cmd:
			if !ok {
				logger.Info("%s:Cmd Read error", t.Name)
				break
			}
			next, fun = t.SysControlEntry(t.Name, msg)
			break
		case msg, ok := <-t.ReadPipe():
			if !ok {
				logger.Info("%s:DataPipe Read error", t.Name)
				break
			}
			next, fun = t.BackDoorEntry(msg)
			break
		case msg, ok := <-t.BUS.ReadPipe():
			if !ok {
				logger.Info("%s:DataPipe Read error", t.Name)
				break
			}
			logger.Debug("%s:read bus chan msg:%+v", t.Name, msg)
			next, fun = t.BusSchedule(msg)
			if fun == cm.ProcessPipeFull {
				logger.Warn("%s:need do something when full", t.Name)
			}
			break
		}

		switch next {
		case cm.NextActionBreak:
			break
		case cm.NextActionReturn:
			return
		case cm.NextActionContinue:
		}
	}
	return
}
