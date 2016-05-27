package timer

import (
	cm "library/core/controlmsg"
	dm "library/core/datamsg"
	"library/logger"
	"service"
)

func (t *timerType) Start(name string, bus *dm.DataMsgPipe) bool {
	logger.Info("timer start running")
	t.Name = name
	t.Output = bus
	go t.job()
	return true
}

func (t *timerType) Pause() bool {
	return true
}

func (t *timerType) Resume() bool {
	return true
}

func (t *timerType) Exit() bool {
	return true
}

func (t *timerType) ControlEntry(msg *cm.ControlMsg) (int, int) {
	switch msg.MsgType {
	case cm.ControlMsgExit:
		logger.Info("ControlMsgPipe.Cmd Read %d", msg.MsgType)
		t.Echo <- &cm.ControlMsg{MsgType: cm.ControlMsgExit}
		logger.Info("timer exit")
		return Return, service.FunOK
	case cm.ControlMsgPause:
		logger.Info("timer paused")
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
		logger.Info("timer resumed")
	}
	return Continue, service.FunOK
}
