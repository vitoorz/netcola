package service

import (
	"sync"
)

type Page interface {
	CRUD(interface{}) bool
	Inspect() string
}

type BufferPool struct {
	Lock sync.Mutex
	Pool []*Page
}

func NewBufferPool() *BufferPool {
	t := &BufferPool{}
	t.Pool = make([]*Page, 0)
	return t
}

func (t *BufferPool) Len() int {
	t.Lock.Lock()
	defer t.Lock.Unlock()
	return len(t.Pool)
}

func (t *BufferPool) Reborn() ([]*Page, bool) {
	if t.Len() <= 0 {
		return nil, false
	}

	pool := t.Pool
	newPool := make([]*Page, 0)

	t.Lock.Lock()
	defer t.Lock.Unlock()
	t.Pool = newPool
	return pool, true
}

func (t *BufferPool) Append(page *Page) {
	t.Lock.Lock()
	defer t.Lock.Unlock()
	t.Pool = append(t.Pool, page)
}
