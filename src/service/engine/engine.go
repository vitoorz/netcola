package engine

import (
	"runtime"
)

import (
	cm "library/core/controlmsg"
	dm "library/core/datamsg"
	"library/logger"
	svc "service/core"
)

type engineT struct {
	svc.Service
}

func NewEngineDefine(bus *dm.DataMsgPipe) *engineT {
	t := &engineT{}
	t.Service = *svc.NewService("")
	t.State = svc.StateInit
	t.ControlMsgPipe = *cm.NewControlMsgPipe()
	t.BUS = bus
	return t
}

func (t *engineT) StartEngine(pipe *dm.DataMsgPipe) {
	logger.Info("engine start running")
	go t.engine(pipe)
}

func (t *engineT) ControlEntry() *cm.ControlMsgPipe {
	return &t.ControlMsgPipe
}

func (t *engineT) engine(datapipe *dm.DataMsgPipe) (err interface{}) {
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
		case msg, ok := <-datapipe.ReadRecvChan():
			if !ok {
				logger.Info("ReadRecvChan Read error")
				break
			}
			logger.Info("recv data:%+v", msg)
		case msg, ok := <-t.Cmd:
			if !ok {
				logger.Info("ControlMsgPipe.Cmd Read error")
				break
			}
			if msg.MsgType == cm.ControlMsgExit {
				logger.Info("ControlMsgPipe.Cmd Read %d", msg.MsgType)
				t.Echo <- &cm.ControlMsg{MsgType: cm.ControlMsgExit}
				logger.Info("engine exit")
				return nil
			}
		}
	}
	return nil
}
