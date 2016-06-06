package engine

import (
	cm "library/core/controlmsg"
	dm "library/core/datamsg"
	"library/logger"
)

func (t *engineType) Start(bus *dm.DataMsgPipe) bool {
	logger.Info("%s:start running", t.Name)
	go t.engine()
	return true
}

func (t *engineType) Pause() bool {
	logger.Info("%s:will pause", t.Name)
	t.Cmd <- &cm.ControlMsg{MsgType: cm.ControlMsgPause}
	echo := <-t.Echo
	if echo.MsgType != cm.ControlMsgPause {
		return false
	}
	return true
}

func (t *engineType) Resume() bool {
	logger.Info("%s:will resume", t.Name)
	t.Cmd <- &cm.ControlMsg{MsgType: cm.ControlMsgResume}
	echo := <-t.Echo
	if echo.MsgType != cm.ControlMsgResume {
		return false
	}
	return true
}

func (t *engineType) Exit() bool {
	logger.Info("%s:will exit", t.Name)
	t.Cmd <- &cm.ControlMsg{MsgType: cm.ControlMsgExit}
	echo := <-t.Echo
	if echo.MsgType != cm.ControlMsgExit {
		return false
	}
	return true
}

func (t *engineType) ControlHandler(msg *cm.ControlMsg) (int, int) {
	switch msg.MsgType {
	case cm.ControlMsgExit:
		logger.Info("%s:ControlMsgPipe.Cmd Read %d", t.Name, msg.MsgType)
		t.Echo <- &cm.ControlMsg{MsgType: cm.ControlMsgExit}
		logger.Info("%s:exit", t.Name)
		return cm.NextActionReturn, cm.ProcessStatOK
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
	return cm.NextActionContinue, cm.ProcessStatOK
}
