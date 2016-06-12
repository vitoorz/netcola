package syssignal

import (
	"os"
	"sync"
)

import (
	cm "library/core/controlmsg"
	"library/logger"
)

type signalWorker struct {
	data    interface{}
	handler func(interface{}) *cm.ControlMsg
}

type signalManager struct {
	sync.Mutex
	c           *cm.ControlMsgPipe
	sig         os.Signal
	workerChain []*signalWorker
}

func (sm *signalManager) signalWork() {
	sm.Lock()
	defer sm.Unlock()
	for i, w := range sm.workerChain {
		logger.Info("handle No.%d in signal chain", i)
		sm.c.Echo <- w.handler(w.data)
	}
	logger.Info("finish handle chain for %s", sm.sig.String())
}
