package service

import (
	"sync"
)

import (
	cm "library/core/controlmsg"
	dm "library/core/datamsg"
	"library/logger"
)

type BufferPool struct {
	sync.Mutex
	Cond *sync.Cond //as a conditional variable
	Host *Service
	cm.ControlMsgPipe
	Pool []*dm.DataMsg
}

func NewBufferPool(h *Service) *BufferPool {
	t := &BufferPool{}
	t.Host = h
	t.Pool = make([]*dm.DataMsg, 0)
	t.ControlMsgPipe = *cm.NewControlMsgPipe()
	t.Cond = sync.NewCond(&sync.Mutex{})
	return t
}

func (t *BufferPool) Len() int {
	t.Lock()
	defer t.Unlock()
	return len(t.Pool)
}

func (t *BufferPool) Reborn() []*dm.DataMsg {
	t.Cond.L.Lock()
	for len(t.Pool) <= 0 { //see usage of cond.wait(), this loop is needed
		t.Cond.Wait()
	}
	t.Cond.L.Unlock()
	logger.Info("%s:buffer reborn wake up", t.Host.Name)

	pool := t.Pool
	newPool := make([]*dm.DataMsg, 0)

	t.Lock()
	defer t.Unlock()
	t.Pool = newPool
	return pool
}

func (t *BufferPool) Append(msg *dm.DataMsg) {
	t.Lock()
	defer t.Unlock()
	t.Pool = append(t.Pool, msg)
	t.Cond.Signal()
}

func (t *BufferPool) Daemon() {
	for {
		todo := t.Reborn() //, hasPage
		for _, msg := range todo {
			ok := t.Host.Instance.DataHandler(msg)
			if !ok {
				logger.Error("%s:DataHandler failed", t.Host.Name)
			}
		}
	}
}
