package engine

import (
	cm "library/core/controlmsg"
	"library/logger"
)

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
