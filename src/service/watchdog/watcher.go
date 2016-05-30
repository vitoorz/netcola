package watchdog

import (
	cm "library/core/controlmsg"
	"library/logger"
	"os"
	"runtime"
	"runtime/pprof"
	"service"
	"time"
)

//should have no overlap with pre-defined control message type
const ServiceName = "watcher"

const (
	Break = iota
	Continue
	Return
)

type watcherType struct {
	service.Service
	objects map[string]int64
}

func NewWatcher() *watcherType {
	t := &watcherType{}
	t.Service = *service.NewService(ServiceName)
	t.objects = make(map[string]int64)
	return t
}

func (t *watcherType) WatchObjStart(obj string) {
	t.WriteCmdNonblock(&cm.ControlMsg{watchCmdStartWatch, obj})
}

func (t *watcherType) WatchObjOver(obj string) {
	t.WriteCmdNonblock(&cm.ControlMsg{watchCmdEndWatch, obj})
}

func (t *watcherType) watch() {
	logger.Info("watcher service running")

	tickChan := time.NewTicker(time.Second).C
	var next, fun int = Continue, service.FunUnknown
	for {
		select {
		case <-tickChan:
			t.onTick()
			break
		case msg, ok := <-t.Cmd:
			if !ok {
				logger.Info("Cmd Read error")
				break
			}
			next, fun = t.ControlEntry(msg)
			if fun != service.FunOK {
				logger.Info("watcher control chan full")
			}
			break
		}

		switch next {
		case Break:
			break
		case Return:
			return
		case Continue:
		}
	}
	return
}

func (t *watcherType) onTick() {
	curTime := time.Now().Unix()
	for obj, startTime := range t.objects {
		if t.printPProfAfter2Second(obj, curTime-startTime) {
			delete(t.objects, obj)
		}
	}
}

func (t *watcherType) printPProfAfter2Second(object string, time int64) bool {
	logger.Warn("watched object %s for %d seconds", object, time)
	if time >= 2 {
		profile := pprof.Lookup("goroutine")

		logger.Warn("%d goroutine(s) are currently in proccess, detail...", runtime.NumGoroutine())
		profile.WriteTo(os.Stdout, 3)
		logger.Warn("-----------over----------")

		return true
	}
	return false
}
