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

func (t *BufferPool) Len() int {
	t.Cond.L.Lock()
	defer t.Cond.L.Unlock()
	return len(t.Pool)
}

//consumer
func (t *BufferPool) reborn() []*dm.DataMsg {
	t.Cond.L.Lock()
	for len(t.Pool) <= 0 { //see usage of cond.wait(), this loop is needed
		t.Cond.Wait()
	}
	//logger.Info("%s:buffer reborn wake up", t.Host.Name)
	pool := t.Pool

	t.Pool = make([]*dm.DataMsg, 0)
	t.Cond.L.Unlock()
	return pool
}

//provider
func (t *BufferPool) Append(msg *dm.DataMsg) {
	t.Cond.L.Lock()
	t.Pool = append(t.Pool, msg)
	t.Cond.Signal()
	t.Cond.L.Unlock()
}

//should use lock to block for loop
func (t *BufferPool) Daemon() {

	for {
		todo := t.reborn()
		for _, msg := range todo {
			//logger.Debug("todo elem:%+v, host:%v", msg, t.Host.Name)
			ok := t.Host.Instance.DataHandler(msg)
			if !ok {
				logger.Error("%s:DataHandler failed", t.Host.Name)
			}
		}
	}
}
