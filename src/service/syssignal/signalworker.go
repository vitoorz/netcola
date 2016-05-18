package syssignal

import (
	"os"
	"sync"
)

import (
	"library/logger"
)

type signalWorker struct {
	data     interface{}
	handler  func(interface{}) int
	echoChan chan int
}

type signalManager struct {
	sync.Mutex
	sig         os.Signal
	workerChain []*signalWorker
}

func (sm *signalManager) signalWork() {
	sm.Lock()
	defer sm.Unlock()
	for i, w := range sm.workerChain {
		logger.Info("handle No.%d in signal chain", i)
		w.echoChan <- w.handler(w.data)
	}
	logger.Info("finish handle chain for %s", sm.sig.String())
}
