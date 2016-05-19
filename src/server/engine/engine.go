package engine

import (
	"runtime"
	"sync"
	//"time"
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

func NewEngineDefine() *engineDefine {
	e := &engineDefine{State: StateInit}
	e.ControlMsgPipe = *cm.NewControlMsgPipe()
	return e
}

func (eg *engineDefine) StartEngine(pipe *netmsg.NetMsgPipe) {
	logger.Info("engine start running")
	go eg.engine(pipe)
}

func (eg *engineDefine) ControlEntry() *cm.ControlMsgPipe {
	return &eg.ControlMsgPipe
}

func (eg *engineDefine) engine(pipe *netmsg.NetMsgPipe) (err interface{}) {
	// catch panic
	defer func() {
		if x := recover(); x != nil {
			logger.Error("Engine job panic: %v", x)
			logger.Stack()
		}
	}()

	// lock this go-routine to a system thread
	runtime.LockOSThread()

	//tickChan := time.NewTicker(time.Millisecond * 100).C
	logger.Info("engine cycle start")
	for {
		select {
		case msg, ok := <-eg.Cmd:
			if !ok {
				logger.Info("ControlMsgPipe.Cmd Read error")
				break
			}
			if msg.MsgType == cm.ControlMsgExit {
				logger.Info("ControlMsgPipe.Cmd Read %d", msg.MsgType)
				eg.Echo <- &cm.ControlMsg{MsgType: cm.ControlMsgExit}
				logger.Info("engine exit")
				return nil
			}
		}
	}
	return nil
}
