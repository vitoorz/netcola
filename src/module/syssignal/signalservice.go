package syssignal

import (
	"os"
	"os/signal"
	"sync"
)

import (
	cm "library/core/controlmsg"
	"library/logger"
)

// signal module
type SignalModule struct {
	sync.Mutex
	cm.ControlMsgPipe
	listenCH chan os.Signal               //chan for golang std library use
	selector map[os.Signal]*signalManager //key: signal in std library; value: the handler chain
}

// Alloc memory for signal module
func NewSignalModule() *SignalModule {
	s := &SignalModule{
		listenCH: make(chan os.Signal),
		selector: make(map[os.Signal]*signalManager)}
	s.ControlMsgPipe = *cm.NewControlMsgPipe()
	return s
}

// Start signal daemon routine
func (s *SignalModule) InitSignalModule() {
	logger.Info("signal module routine bring up")
	go s.signalRoutine()
}

// Start signal daemon routine
func (s *SignalModule) Exit() {
	//todo:
	// 1. unregister all the signals
	// 2. recycle the memory use in module
}

// Register signal handler to module
// Register to the same signal would be executed according the register sequence
func (s *SignalModule) RegisterSignalCallback(
	sig os.Signal, f func(interface{}) *cm.ControlMsg, d interface{}) {

	// register the signal
	signal.Notify(s.listenCH, sig)
	logger.Info("register listening on signal:%s", sig.String())

	s.Lock()
	defer s.Unlock()

	_, ok := s.selector[sig]
	if !ok {
		s.selector[sig] = &signalManager{
			c:           &s.ControlMsgPipe,
			sig:         sig,
			workerChain: make([]*signalWorker, 0),
		}
	}

	worker := &signalWorker{
		data:    d,
		handler: f,
	}

	mgr := s.selector[sig]
	mgr.Lock()
	mgr.workerChain = append(mgr.workerChain, worker)
	mgr.Unlock()
	return
}

// module daemon routine
func (s *SignalModule) signalRoutine() {
	for {
		sig := <-s.listenCH
		logger.Info("signal caught: %s", sig.String())

		s.Lock()
		manager, ok := s.selector[sig]
		s.Unlock()
		if !ok {
			logger.Error("Didn't register signal:%s", sig.String())
			continue
		}
		manager.signalWork()
	}
}
