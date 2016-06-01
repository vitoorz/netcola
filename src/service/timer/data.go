package timer

import (
	dm "library/core/datamsg"
	"library/logger"
	"service"
	"service/timer/task"
	"time"
	ts "types/service"
)

func (t *timerType) dataEntry(msg *dm.DataMsg) (operate int, funCode int) {
	defer func() {
		if x := recover(); x != nil {
			logger.Error("timer panic: %v", x)
			logger.Stack()
		}
		operate = Continue
		funCode = service.FunPanic
	}()

	logger.Info("timer: data msg:%+v,payload:%s", msg, msg.Payload.([]byte))

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

	return Continue, service.FunOK
}

func (t *timerType) callBack(e ts.Event, msg *dm.DataMsg) {
	go func() {
		wakeAt := <-e.TimerObject.C
		logger.Info("event wake up:at:%s", wakeAt.String())
		task.DoLater(msg)
		msg.Sender = t.Name
		ok := t.Output.WritePipeNB(msg)
		if !ok {
			// channel full
			//return Continue, service.FunDataPipeFull
		}
	}()
}
