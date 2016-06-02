package mongo

import (
	"sync"
)

import (
	ts "types/service"
)

type DirtyPool struct {
	Lock sync.Mutex
	Pool []*ts.Dirty
}

//var dirtyManager *DirtyPool = &DirtyPool{}

func NewDirtyPool() *DirtyPool {
	t := &DirtyPool{}
	t.Pool = make([]*ts.Dirty, 0)
	return t
}

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
