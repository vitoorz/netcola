package engine

import (
	"runtime"
)

import (
	cm "library/core/controlmsg"
	dm "library/core/datamsg"
	"library/logger"
	"service"
)

type engineType struct {
	service.Service
}

func NewEngine(bus *dm.DataMsgPipe) *engineType {
	t := &engineType{}
	t.Service = *service.NewService("")
	t.State = service.StateInit
	t.ControlMsgPipe = *cm.NewControlMsgPipe()
	t.BUS = bus
	return t
}

func (t *engineType) ControlEntry() *cm.ControlMsgPipe {
	return &t.ControlMsgPipe
}

func (t *engineType) engine(datapipe *dm.DataMsgPipe) (err interface{}) {
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
		case msg, ok := <-datapipe.ReadDownChan():
			if !ok {
				logger.Info("ReadDownChan Read error")
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

func (t *engineType) Init() bool {
	logger.Info("engine init")
	return true
}

func (t *engineType) Start() bool {
	logger.Info("engine start running")
	go t.engine(t.BUS)
	return true
}

func (t *engineType) Pause() bool {
	return true
}

func (t *engineType) Exit() bool {
	return true
}

func (t *engineType) Kill() bool {
	return true
}

func (t *engineType) CommonService() *service.Service {
	return &t.Service
}
