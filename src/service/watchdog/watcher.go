package library

import (
	"library/logger"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

type Watcher struct {
	start   chan string
	end     chan string
	objects map[string]int64 //string for object's unique identifier
	handler func(object string, watchedNanoSecond int64) (stopWatch bool)
}

func NewWatcher(handler func(string, int64) bool) *Watcher {
	if handler == nil {
		handler = printPProfAfter2Second
	}
	w := &Watcher{
		start:   make(chan string),
		end:     make(chan string),
		objects: make(map[string]int64),
		handler: handler,
	}
	go w.startWatchRoutine()

	return w
}

func (w *Watcher) WatchObjStart(obj string) {
	w.start <- obj
}

func (w *Watcher) WatchObjOver(obj string) {
	w.end <- obj
}

func (w *Watcher) SetWatcherHandler(handler func(string, int64) bool) {
	w.handler = handler
}

func (w *Watcher) startWatchRoutine() {
	tickChan := time.NewTicker(time.Second).C
	for {
		select {
		case <-tickChan:
			curTime := time.Now().UTC().UnixNano()
			for obj, startTime := range w.objects {
				if w.handler(obj, curTime-startTime) {
					delete(w.objects, obj)
				}
			}
		case obj, ok := <-w.start:
			if ok {
				curTime := time.Now().UTC().UnixNano()
				w.objects[obj] = curTime
			}

		case obj, ok := <-w.end:
			if ok {
				delete(w.objects, obj)
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
