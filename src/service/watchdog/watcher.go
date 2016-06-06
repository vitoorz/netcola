package watchdog

import (
	cm "library/core/controlmsg"
	"service"
)

//should have no overlap with pre-defined control message type
const ServiceName = "watcher"

type watcherType struct {
	service.Service
	objects map[string]int64
}

func NewWatcher(name string) *watcherType {
	t := &watcherType{}
	t.Service = *service.NewService(ServiceName)
	t.objects = make(map[string]int64, 0)
	t.Name = name
	t.Buffer = service.NewBufferPool(&t.Service)
	t.Instance = t
	return t
}

func (t *watcherType) WatchObjStart(obj string) {
	t.WriteCmdNonblock(&cm.ControlMsg{watchCmdStartWatch, obj})
}

func (t *watcherType) WatchObjOver(obj string) {
	t.WriteCmdNonblock(&cm.ControlMsg{watchCmdEndWatch, obj})
}
