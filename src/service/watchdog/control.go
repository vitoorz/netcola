package watchdog

import (
	cm "library/core/controlmsg"
	dm "library/core/datamsg"
	"library/logger"
	"time"
)

const (
	watchCmdStartWatch = iota + cm.ControlMsgMax + 1
	watchCmdEndWatch
)

func (t *watcherType) Start(bus *dm.DataMsgPipe) bool {
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
	case watchCmdStartWatch:
		object := msg.Payload.(string)
		t.objects[object] = time.Now().Unix()
	case watchCmdEndWatch:
		object := msg.Payload.(string)
		delete(t.objects, object)
	case cm.ControlMsgTick:
		t.onTick()
	default:
		logger.Error("watcher received invalid message type %v", msg.MsgType)
	}

	return cm.NextActionContinue, cm.ProcessStatOK
}

func (t *watcherType) onTick() {
	curTime := time.Now().Unix()
	for obj, startTime := range t.objects {
		past := curTime - startTime
		logger.Warn("watched object %s for %d seconds", obj, past)
		if past >= 2 {
			logger.PProf()
			delete(t.objects, obj)
		}
	}
}
