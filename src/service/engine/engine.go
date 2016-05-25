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
	//defer func() {
	//	if x := recover(); x != nil {
	//		logger.Error("Engine job panic: %v", x)
	//		logger.Stack()
	//	}
	//}()

	logger.Info("engine cycle start")

	for {

		select {
		case msg, ok := <-t.Cmd:
			if !ok {
				logger.Info("Cmd Read error")
				break
			}
			switch msg.MsgType {
			case cm.ControlMsgExit:
				logger.Info("ControlMsgPipe.Cmd Read %d", msg.MsgType)
				t.Echo <- &cm.ControlMsg{MsgType: cm.ControlMsgExit}
				logger.Info("engine exit")
				return nil
			case cm.ControlMsgPause:
				logger.Info("engine paused")
				t.Echo <- &cm.ControlMsg{MsgType: cm.ControlMsgPause}
				for {
					var resume bool = false
					select {
					case msg, ok := <-t.Cmd:
						if !ok {
							logger.Info("Cmd Read error")
							break
						}
						switch msg.MsgType {
						case cm.ControlMsgResume:
							t.Echo <- &cm.ControlMsg{MsgType: cm.ControlMsgResume}
							resume = true
							break
						}
					}
					if resume {
						break
					}
				}
				logger.Info("engine resumed")
			}
		case msg, ok := <-datapipe.ReadDownChan():
			if !ok {
				logger.Info("DownChan Read error")
				break
			}
			logger.Debug("engine recv data:%+v", msg)
			service.ServicePool.SendDown(msg)
		}
	}
	return nil
}
