package engine

import (
	cm "library/core/controlmsg"
	dm "library/core/datamsg"
	"library/logger"
	"service"
	"service/job"
)

const ServiceName = "engine"

type engineType struct {
	service.Service
}

func NewEngine(bus *dm.DataMsgPipe) *engineType {
	t := &engineType{}
	t.Service = *service.NewService(ServiceName)
	t.BUS = bus
	return t
}

func (t *engineType) ControlEntry() *cm.ControlMsgPipe {
	return &t.ControlMsgPipe
}

func (t *engineType) engine(datapipe *dm.DataMsgPipe) (err interface{}) {
	defer func() {
		if x := recover(); x != nil {
			logger.Error("Engine job panic: %v", x)
			logger.Stack()
		}
	}()

	logger.Info("engine cycle start")
	for {
		select {
		case msg, ok := <-datapipe.ReadDownChan():
			if !ok {
				logger.Info("ReadDownChan Read error")
				break
			}
			logger.Debug("engine recv data:%+v", msg)
			if msg.Receiver == job.ServiceName {
				service.ServicePool.SendDown(msg)
			}
		case msg, ok := <-t.Cmd:
			if !ok {
				logger.Info("ControlMsgPipe.Cmd Read error")
				break
			}
			if msg.MsgType == cm.ControlMsgExit {
				logger.Info("ControlMsgPipe.Cmd Read %d", msg.MsgType)
				t.Echo <- &cm.ControlMsg{MsgType: cm.ControlMsgExit}
				logger.Info("engine exit")
				return nil
			}
		}
	}
	return nil
}

func (t *engineType) Init() bool {
	logger.Info("engine init")
	return true
}

func (t *engineType) Start() bool {
	logger.Info("engine start running")
	go t.engine(t.BUS)
	return true
}

func (t *engineType) Pause() bool {
	return true
}

func (t *engineType) Exit() bool {
	return true
}

func (t *engineType) Kill() bool {
	return true
}
