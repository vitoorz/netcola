package syssignal

import (
	"os"
	"os/signal"
)

import (
	"library/logger"
)

var listenCH chan os.Signal

func RegisterSignalCallback(s os.Signal, f func(interface{}), d interface{}) chan int {

	// 注册要监听的信号
	signal.Notify(listenCH, s)
	logger.Info("register listening on signal:%s", s.String())

	selector.Lock()
	defer selector.Unlock()

	_, ok := selector.S[s]
	if !ok {
		selector.S[s] = &SignalManager{
			Sig:         s,
			WorkerChain: make([]*SignalWorker, 0),
		}
	}

	sw := &SignalWorker{
		Data:     d,
		Handler:  f,
		EchoChan: make(chan int),
	}

	mgr := selector.S[s]
	mgr.Lock.Lock()
	mgr.WorkerChain = append(mgr.WorkerChain, sw)
	mgr.Lock.Unlock()

	return sw.EchoChan
}

func InitSignalService() {
	logger.Info("signal service routine bring up")

	listenCH = make(chan os.Signal)
	go signalRoutine()
}

// sevice main routine
func signalRoutine() {
	for {
		s := <-listenCH
		logger.Info("signal caught: %s", s.String())

		selector.Lock()
		manager, ok := selector.S[s]
		selector.Unlock()
		if !ok {
			logger.Error("Didn't register signal:%s", s.String())
			continue
		}
		signalWork(manager)
	}
}

func signalWork(sm *SignalManager) {
	sm.Lock.Lock()
	defer sm.Lock.Unlock()
	for i, w := range sm.WorkerChain {
		logger.Info("handle No.%d in signal chain", i)
		w.Handler(w.Data)
		w.EchoChan <- SysSignalHandled
	}
	logger.Info("finish handle chain for %s", sm.Sig.String())
}
