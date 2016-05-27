package engine

import (
	cm "library/core/controlmsg"
	dm "library/core/datamsg"
	"library/logger"
)

func (t *engineType) Start(name string, bus *dm.DataMsgPipe) bool {
	logger.Info("engine start running")
	t.Name = name
	go t.engine()
	return true
}

func (t *engineType) Pause() bool {
	logger.Info("engine will pause")
	t.Cmd <- &cm.ControlMsg{MsgType: cm.ControlMsgPause}
	echo := <-t.Echo
	if echo.MsgType != cm.ControlMsgPause {
		return false
	}
	return true
}

func (t *engineType) Resume() bool {
	logger.Info("engine will resume")
	t.Cmd <- &cm.ControlMsg{MsgType: cm.ControlMsgResume}
	echo := <-t.Echo
	if echo.MsgType != cm.ControlMsgResume {
		return false
	}
	return true
}

func (t *engineType) Exit() bool {
	logger.Info("engine will exit")
	t.Cmd <- &cm.ControlMsg{MsgType: cm.ControlMsgExit}
	echo := <-t.Echo
	if echo.MsgType != cm.ControlMsgExit {
		return false
	}
	return true
}

func (t *engineType) ControlEntry(msg *cm.ControlMsg) (int, bool) {
	switch msg.MsgType {
	case cm.ControlMsgExit:
		logger.Info("ControlMsgPipe.Cmd Read %d", msg.MsgType)
		t.Echo <- &cm.ControlMsg{MsgType: cm.ControlMsgExit}
		logger.Info("engine exit")
		return Return, true
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
	return Continue, true
}
