package timer

import (
	dm "netcola/library/core/datamsg"
	"netcola/library/logger"
	"netcola/service/timer/task"
	ts "netcola/types/service"
	"time"
)

func (t *timerType) callBack(e ts.Event, msg *dm.DataMsg) {
	go func() {
		wakeAt := <-e.TimerObject.C
		logger.Info("event wake up:at:%s", wakeAt.String())
		task.DoLater(msg)
		msg.Receiver, msg.Sender = msg.Sender, msg.Receiver
		t.output.WritePipeNoBlock(msg)
	}()
}

func (t *timerType) DataHandler(msg *dm.DataMsg) bool {
	m, ok := msg.Meta(t.Name)
	if !ok {
		logger.Error("meta error for timer")
	} else {
		ev := m.(ts.Event)
		sleepTime := time.Second * time.Duration(ev.When-time.Now().Unix())
		if ev.When > time.Now().Unix() {
			ev.TimerObject = time.NewTimer(sleepTime)
			logger.Info("event will wake after:%d nano", sleepTime)
			t.callBack(ev, msg)
		}
	}
	return true
}
