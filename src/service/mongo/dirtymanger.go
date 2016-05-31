package mongo

import (
	"sync"
)

import (
	cm "library/core/controlmsg"
	dm "library/core/datamsg"
	"library/logger"
	ts "types/service"
)

type DirtyPool struct {
	Lock sync.Mutex
	cm.ControlMsgPipe
	dm.DataMsgPipe
	Pool []*ts.MongoDirty
}

//var dirtyManager *DirtyPool = &DirtyPool{}

func NewDirtyPool() *DirtyPool {
	t := &DirtyPool{}
	t.ControlMsgPipe = *cm.NewControlMsgPipe()
	t.DataMsgPipe = *dm.NewDataMsgPipe(0)
	t.Pool = make([]*ts.MongoDirty, 0)
	return t
}

//
//func (t *DirtyPool) init() {
//	t.Lock.Lock()
//	defer t.Lock.Unlock()
//	t.Pool = make([]*ts.MongoDirty, 0)
//}

func (t *DirtyPool) reborn() []*ts.MongoDirty {
	dirtyPool := t.Pool
	newDirtyPool := make([]*ts.MongoDirty, 0)

	t.Lock.Lock()
	defer t.Lock.Unlock()
	t.Pool = newDirtyPool
	return dirtyPool
}

func (t *DirtyPool) addDirty(trash *ts.MongoDirty) {
	t.Lock.Lock()
	defer t.Lock.Unlock()
	t.Pool = append(t.Pool, trash)
}

func (t *DirtyPool) recycleDaemon() {
	todo := t.reborn()
	var ok bool = false
	for _, dirty := range todo {
		switch dirty.Action {
		case ts.MongoActionCreate:
			for ok = dirty.Create(); !ok; {
				// break session and re-dial?
				logger.Info("Create:%s", dirty.Inspect())
			}
		case ts.MongoActionRead:
			for ok = dirty.Read(); !ok; {
				logger.Info("Read:%s", dirty.Inspect())
			}
		case ts.MongoActionUpdate:
			for ok = dirty.Update(); !ok; {
				logger.Info("Update:%s", dirty.Inspect())
			}
		case ts.MongoActionDelete:
			for ok = dirty.Delete(); !ok; {
				logger.Info("Delete:%s", dirty.Inspect())
			}
		default:
			logger.Error("dirty has error", dirty.Inspect())
		}
		if !ok {
			logger.Error("operation failed")
		}
	}
}
