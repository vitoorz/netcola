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
	object, ok := msg.Payload.(string)
	if !ok {
		logger.Error("invalid watcher cmd: type: %d, payload: %v", msg.MsgType, msg.Payload)
		return cm.NextActionContinue, cm.ProcessStatIgnore
	}

	switch msg.MsgType {
	case watchCmdStartWatch:
		t.objects[object] = time.Now().Unix()
	case watchCmdEndWatch:
		delete(t.objects, object)
	default:
		logger.Error("watcher received invalid message type %v", msg.MsgType)
	}

	return cm.NextActionContinue, cm.ProcessStatOK
}
