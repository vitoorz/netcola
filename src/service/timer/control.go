package timer

import (
	cm "library/core/controlmsg"
	dm "library/core/datamsg"
	"library/logger"
	"service"
)

func (t *timerType) Start(bus *dm.DataMsgPipe) bool {
	logger.Info("%s:start running", t.Name)
	t.output = bus
	go t.timer()
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

func (t *timerType) controlEntry(msg *cm.ControlMsg) (int, int) {
	switch msg.MsgType {
	case cm.ControlMsgExit:
		logger.Info("%s:ControlMsgPipe.Cmd Read %d", t.Name, msg.MsgType)
		t.Echo <- &cm.ControlMsg{MsgType: cm.ControlMsgExit}
		logger.Info("%s:exit", t.Name)
		return Return, service.FunOK
	case cm.ControlMsgPause:
		logger.Info("%s:paused", t.Name)
		t.Echo <- &cm.ControlMsg{MsgType: cm.ControlMsgPause}
		for {
			var resume bool = false
			select {
			case msg, ok := <-t.Cmd:
				if !ok {
					logger.Info("%s:Cmd Read error", t.Name)
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
		logger.Info("%s:resumed", t.Name)
	}
	return Continue, service.FunOK
}
