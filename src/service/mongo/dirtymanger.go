package mongo

import (
	"sync"
)

import (
	cm "library/core/controlmsg"
	dm "library/core/datamsg"
	"library/logger"
	"time"
	ts "types/service"
)

type DirtyPool struct {
	Lock sync.Mutex
	cm.ControlMsgPipe
	dm.DataMsgPipe
	Pool []*ts.Dirty
}

//var dirtyManager *DirtyPool = &DirtyPool{}

func NewDirtyPool() *DirtyPool {
	t := &DirtyPool{}
	t.ControlMsgPipe = *cm.NewControlMsgPipe()
	t.DataMsgPipe = *dm.NewDataMsgPipe(0)
	t.Pool = make([]*ts.Dirty, 0)
	return t
}

//
//func (t *DirtyPool) init() {
//	t.Lock.Lock()
//	defer t.Lock.Unlock()
//	t.Pool = make([]*ts.MongoDirty, 0)
//}

func (t *DirtyPool) PoolLen() int {
	t.Lock.Lock()
	defer t.Lock.Unlock()
	return len(t.Pool)
}

func (t *DirtyPool) reborn() ([]*ts.Dirty, bool) {
	if t.PoolLen() <= 0 {
		return nil, false
	}

	dirtyPool := t.Pool
	newDirtyPool := make([]*ts.Dirty, 0)

	t.Lock.Lock()
	defer t.Lock.Unlock()
	t.Pool = newDirtyPool
	return dirtyPool, true
}

func (t *DirtyPool) addDirty(trash *ts.Dirty) {
	t.Lock.Lock()
	defer t.Lock.Unlock()
	t.Pool = append(t.Pool, trash)
}

func (t *DirtyPool) recycleDaemon() {
	for {
		todo, hasDirty := t.reborn()
		if !hasDirty {
			time.Sleep(time.Millisecond)
			continue
		}
		var ok bool = false
		for _, dirty := range todo {
			ok = (*dirty).CRUD()
			for ok = (*dirty).CRUD(); !ok; {
				logger.Info("Read:%s", (*dirty).Inspect())
			}
			if !ok {
				logger.Error("operation failed")
			}
		}
	}
}
