package watchdog

import (
	cm "library/controlmsg"
	"library/logger"
	"os"
	"runtime"
	"runtime/pprof"
	. "service"
	"time"
)

//should have no overlap with pre-defined control message type
const (
	watchCmdStartWatch = cm.ControlMsgMax + iota
	watchCmdEndWatch
	watchCmdSetPassTickHandler
)

type PassTickHandler func(object string, watchedNanoSecond int64) (stopWatch bool)

type Watcher struct {
	service         *ServiceCommon
	objects         map[string]int64
	passTickHandler PassTickHandler
}

func NewWatcher() *Watcher {
	w := &Watcher{
		objects:         make(map[string]int64),
		passTickHandler: printPProfAfter2Second,
	}

	s := NewServiceCommon(1e9, nil)
	s.SetRealInstance(w)

	w.service = s

	w.service.SysHandler(cm.ControlMsgTick, handleOnTick)
	w.service.SysHandler(watchCmdSetPassTickHandler, handleWatchSetPassTickHandler)
	w.service.SysHandler(watchCmdStartWatch, handleWatchStart)
	w.service.SysHandler(watchCmdEndWatch, handleWatchEnd)

	return w
}

func (w *Watcher) WatchObjStart(obj string) {
	w.service.HandleSysCmd(&cm.ControlMsg{watchCmdStartWatch, obj})
}

func (w *Watcher) WatchObjOver(obj string) {
	w.service.HandleSysCmd(&cm.ControlMsg{watchCmdEndWatch, obj})
}

func (w *Watcher) SetWatcherHandler(obj PassTickHandler) {
	w.service.HandleSysCmd(&cm.ControlMsg{watchCmdSetPassTickHandler, obj})
}

func (w *Watcher) Start() {
	w.service.Start(nil)
}

func (w *Watcher) Stop() {
	w.service.Stop()
}

func (w *Watcher) Status() int {
	return w.service.Status()
}

func handleWatchSetPassTickHandler(ws interface{}, code int, object interface{}) (int, interface{}) {
	w := ws.(*Watcher)
	w.passTickHandler = object.(PassTickHandler)
	return 0, nil
}

func handleWatchStart(ws interface{}, code int, object interface{}) (int, interface{}) {
	w := ws.(*Watcher)
	w.objects[object.(string)] = time.Now().UnixNano()
	return 0, nil
}

func handleWatchEnd(ws interface{}, code int, object interface{}) (int, interface{}) {
	w := ws.(*Watcher)
	delete(w.objects, object.(string))
	return 0, nil
}

func handleOnTick(ws interface{}, code int, data interface{}) (int, interface{}) {
	w := ws.(*Watcher)

	curTime := time.Now().UTC().UnixNano()
	for obj, startTime := range w.objects {
		if w.passTickHandler(obj, curTime-startTime) {
			delete(w.objects, obj)
		}
	}

	return 0, nil
}

func printPProfAfter2Second(object string, nanoSecond int64) bool {
	logger.Warn("watched object %s for %ds %dns", object, nanoSecond/1e9, nanoSecond%1e9)
	if nanoSecond >= 2e9 {
		profile := pprof.Lookup("goroutine")

		logger.Warn("%d goroutine(s) are currently in proccess, detail...", runtime.NumGoroutine())
		profile.WriteTo(os.Stdout, 3)
		logger.Warn("-----------over----------")

		return true
	}
	return false
}
