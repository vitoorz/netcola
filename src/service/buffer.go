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
	Host *Service
	cm.ControlMsgPipe
	Pool []*dm.DataMsg
}

func NewBufferPool(h *Service) *BufferPool {
	t := &BufferPool{}
	t.Host = h
	t.Pool = make([]*dm.DataMsg, 0)
	t.ControlMsgPipe = *cm.NewControlMsgPipe()
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
	for {
		todo, hasPage := t.Reborn()
		if !hasPage {
			time.Sleep(time.Millisecond)
			continue
		}
		for _, msg := range todo {
			ok := t.Host.Instance.DataHandler(msg)
			if !ok {
				logger.Error("%s:DataHandler failed", t.Host.Name)
			}
		}
	}
}

type BufferHandler interface {
}
