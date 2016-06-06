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

//consumer
func (t *BufferPool) reborn() ([]*dm.DataMsg, bool) {
	t.Cond.L.Lock()
	for len(t.Pool) <= 0 { //see usage of cond.wait(), this loop is needed
		t.Cond.Wait()
	}
	pool := t.Pool
	newPool := make([]*dm.DataMsg, 0)
	t.Pool = newPool
	t.Cond.L.Unlock()

	return pool, true
}

//provider
func (t *BufferPool) Append(msg *dm.DataMsg) {
	t.Pool = append(t.Pool, msg)
	t.Cond.Signal()
}

//should use lock to block for loop
func (t *BufferPool) Daemon() {
	for {
		todo, _ := t.reborn()
		for _, msg := range todo {
			ok := t.Host.Instance.DataHandler(msg)
			if !ok {
				logger.Error("%s:DataHandler failed", t.Host.Name)
			}
		}
	}
}
