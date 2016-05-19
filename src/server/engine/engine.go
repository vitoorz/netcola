package engine

import (
	"runtime"
	"sync"
	"time"
)

import (
	cm "library/controlmsg"
	"library/logger"
	"library/netmsg"
	. "types"
)

type engineDefine struct {
	sync.Mutex
	cm.ControlMsgPipe
	State StateT
}

func (eg *engineDefine) Engine(pipe *netmsg.NetMsgPipe) (err interface{}) {
	// catch panic
	defer func() {
		if x := recover(); x != nil {
			logger.Error("Engine job panic: %v", x)
			logger.Stack()
		}
	}()

	// lock this go-routine to a system thread
	runtime.LockOSThread()

	tickChan := time.NewTicker(time.Millisecond * 100).C
	logger.Info("engine start running")
	for {
		select {

		case msg, ok := <-eg.ControlMsgPipe.Cmd:
			if !ok {
				logger.Info("WorkerCtrlMsg Read error %v", ok)
				break
			}
			if msg.MsgType == 1 {
				logger.Info("WorkerCtrlMsg Read %d", msg.MsgType)
				cm.ControlMsgPipe.Echo <- cm.ControlMsg{MsgType: 2}
				logger.Info("Worker exit")
				return nil
			}
		case <-tickChan:
			break
		}
	}
	return nil
}
