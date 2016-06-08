package support

import (
	"service"
	"service/syswatcher"
)

var SysWatcher = syswatcher.NewWatcher()

func init() {
	SysWatcher.Watching()
}

func BookTicker(service service.IService, stepSecond int64) {
	SysWatcher.BookedBy(service, stepSecond)
}

func WatchObjectStart(object string) {
	SysWatcher.WatchObjectStart(object)
}

func WatchObjectEnd(object string) {
	SysWatcher.WatchObjectEnd(object)
}
