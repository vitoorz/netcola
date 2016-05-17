package syssignal

import (
	"os"
	"sync"
)

const (
	SysSignalHandled = 1
)

type SignalWorker struct {
	Data     interface{}
	Handler  func(interface{})
	EchoChan chan int
}

type SignalManager struct {
	Lock        sync.Mutex
	Sig         os.Signal
	WorkerChain []*SignalWorker
}

// signal handler route map
// map key: signal in std library
// map value: the handler chain
type signalSelector struct {
	sync.Mutex
	S map[os.Signal]*SignalManager
}

var selector signalSelector = signalSelector{S: make(map[os.Signal]*SignalManager)}
