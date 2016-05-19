package library

import (
	"os"
	"time"
	"runtime"
	"runtime/pprof"
	"library/logger"
)

type Watcher struct {
	Start   chan string
	End     chan string
	Objects map[string]int64 //string for object's unique identifier
	Handler func(object string, watchedNanoSecond int64) (stopWatch bool)
}

func NewWatcher(handler func(string, int64)(bool)) *Watcher {
	if handler == nil {
		handler = printPProfAfter2Second
	}
	w := &Watcher{
		Start:   make(chan string),
		End:     make(chan string),
		Objects: make(map[string]int64),
		Handler: handler,
	}
	go w.startWatchRoutine()

	return w
}


func (w *Watcher) WatchObjStart(obj string) {
	w.Start <- obj
}

func (w *Watcher) WatchObjOver(obj string) {
	w.End <- obj
}

func (w *Watcher) startWatchRoutine() {
	tickChan := time.NewTicker(time.Second).C
	for {
		select {
		case <-tickChan:
			curTime := time.Now().UTC().UnixNano()
			for obj, startTime := range w.Objects {
				if w.Handler(obj, curTime - startTime) {
					delete(w.Objects, obj)
				}
			}
		case obj, ok := <-w.Start:
			if ok {
				curTime := time.Now().UTC().UnixNano()
				w.Objects[obj] = curTime
			}

		case obj, ok := <-w.End:
			if ok {
				delete(w.Objects, obj)
			}
		}
	}
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