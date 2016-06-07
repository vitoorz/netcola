package watchdog

import (
	cm "library/core/controlmsg"
	dm "library/core/datamsg"
	"library/logger"
	"service"
	"time"
)

const (
	watchCmdStartWatch = iota + cm.ControlMsgMax + 1
	watchCmdEndWatch
)

func (t *watcherType) Start(bus *dm.DataMsgPipe) bool {
	logger.Info("watcher start running")
	return true
}

func (t *watcherType) Pause() bool {
	return true
}

func (t *watcherType) Resume() bool {
	return true
}

func (t *watcherType) Exit() bool {
	return true
}

func (t *watcherType) ControlHandler(msg *cm.ControlMsg) (int, int) {
	switch msg.MsgType {
	case cm.ControlMsgExit:
		logger.Info("ControlMsgPipe.Cmd Read %d", msg.MsgType)
		t.Echo <- &cm.ControlMsg{MsgType: cm.ControlMsgExit}
		logger.Info("watcher exit")
		return Return, service.FunOK
	case cm.ControlMsgPause:
		logger.Info("watcher paused")
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
		logger.Info("watcher resumed")
	default:
		object, ok := msg.Payload.(string)
		if !ok {
			logger.Error("invalid watcher cmd: type: %d, payload: %v", msg.MsgType, msg.Payload)
			break
		}
		t.objectProcess(msg.MsgType, object)
	}
	return cm.NextActionContinue, cm.ProcessStatOK
}

func (t *watcherType) objectProcess(cmd int, object string) {
	switch cmd {
	case watchCmdStartWatch:
		t.objects[object] = time.Now().Unix()
	case watchCmdEndWatch:
		delete(t.objects, object)
	default:
		logger.Error("watcher received invalid message type %v", cmd)
	}
}
