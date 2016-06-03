package service

import (
	"sync"
)

import (
	cm "library/core/controlmsg"
	dm "library/core/datamsg"
	"library/logger"
	"time"
)

type BufferPool struct {
	Lock sync.Mutex
	Host IService
	cm.ControlMsgPipe
	Pool []*dm.DataMsg
}

func NewBufferPool(h IService) *BufferPool {
	t := &BufferPool{}
	t.Host = h
	t.Pool = make([]*dm.DataMsg, 0)
	t.ControlMsgPipe = *cm.NewControlMsgPipe()
	go t.Daemon()
	return t
}

func (t *BufferPool) Len() int {
	t.Lock.Lock()
	defer t.Lock.Unlock()
	return len(t.Pool)
}

func (t *BufferPool) Reborn() ([]*dm.DataMsg, bool) {
	if t.Len() <= 0 {
		return nil, false
	}

	pool := t.Pool
	newPool := make([]*dm.DataMsg, 0)

	t.Lock.Lock()
	defer t.Lock.Unlock()
	t.Pool = newPool
	return pool, true
}

func (t *BufferPool) Append(msg *dm.DataMsg) {
	t.Lock.Lock()
	defer t.Lock.Unlock()
	t.Pool = append(t.Pool, msg)
}

func (t *BufferPool) Daemon() {
	serviceName := t.Host.Self().Name
	h, ok := t.Host.(BufferHandler)
	if !ok {
		logger.Error("%s:msg is not BufferHandler interface", serviceName)
		return
	}

	for {
		todo, hasPage := t.Reborn()
		if !hasPage {
			time.Sleep(time.Millisecond)
			continue
		}
		for _, msg := range todo {
			h.Handle(msg)
			if !ok {
				logger.Error("%s:Operation failed", serviceName)
			}
		}
	}
}

type BufferHandler interface {
	Handle(*dm.DataMsg) bool
}
