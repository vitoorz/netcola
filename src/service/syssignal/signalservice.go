package syssignal

import (
	"os"
	"os/signal"
	"sync"
)

import (
	"library/logger"
)

// signal handler route map
// map key: signal in std library
// map value: the handler chain
type SignalService struct {
	sync.Mutex
	listenCH chan os.Signal
	selector map[os.Signal]*signalManager
}

func NewSignalService() *SignalService {
	return &SignalService{
		listenCH: make(chan os.Signal),
		selector: make(map[os.Signal]*signalManager)}
}

func (s *SignalService) InitSignalService() {
	logger.Info("signal service routine bring up")

	go s.signalRoutine()
}

// register signal handler to service. when
func (s *SignalService) RegisterSignalCallback(sig os.Signal, f func(interface{}) int, d interface{}) chan int {

	// 注册要监听的信号
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

// sevice main routine
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
