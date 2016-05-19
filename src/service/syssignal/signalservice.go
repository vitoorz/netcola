package syssignal

import (
	"os"
	"os/signal"
	"sync"
)

import (
	"library/logger"
)

// signal service
type SignalService struct {
	sync.Mutex
	listenCH chan os.Signal               //chan for golang std library use
	selector map[os.Signal]*signalManager //key: signal in std library; value: the handler chain
}

// Alloc memory for signal service
func NewSignalService() *SignalService {
	return &SignalService{
		listenCH: make(chan os.Signal),
		selector: make(map[os.Signal]*signalManager)}
}

// Start signal daemon routine
func (s *SignalService) InitSignalService() {
	logger.Info("signal service routine bring up")
	go s.signalRoutine()
}

// Register signal handler to service
// Register to the same signal would be executed according the register sequence
func (s *SignalService) RegisterSignalCallback(sig os.Signal, f func(interface{}) int, d interface{}) chan int {

	// register the signal
	signal.Notify(s.listenCH, sig)
	logger.Info("register listening on signal:%s", sig.String())

	s.Lock()
	defer s.Unlock()

	_, ok := s.selector[sig]
	if !ok {
		s.selector[sig] = &signalManager{
			sig:         sig,
			workerChain: make([]*signalWorker, 0),
		}
	}

	worker := &signalWorker{
		data:     d,
		handler:  f,
		echoChan: make(chan int),
	}

	mgr := s.selector[sig]
	mgr.Lock()
	mgr.workerChain = append(mgr.workerChain, worker)
	mgr.Unlock()

	return worker.echoChan
}

// service daemon routine
func (s *SignalService) signalRoutine() {
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
