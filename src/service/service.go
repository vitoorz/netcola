package service

import (
	"sync"
)

import (
	cm "library/core/controlmsg"
	dm "library/core/datamsg"
	"library/idgen"
	"library/logger"
)

type Service struct {
	sync.RWMutex
	ID    idgen.ObjectID
	Name  string
	State StateT
	cm.ControlMsgPipe
	dm.DataMsgPipe
	Buffer    *BufferPool
	SelfDrive bool
	Instance  IService
}

func NewService(name string) *Service {
	t := &Service{}
	t.ID = idgen.NewObjectID()
	t.Name = name
	t.State = StateInit
	t.ControlMsgPipe = *cm.NewControlMsgPipe()
	t.DataMsgPipe = *dm.NewDataMsgPipe(0)
	t.Buffer = nil
	return t
}

func (t *Service) Self() *Service {
	return t
}

func (t *Service) Background() {
	var next, stat int = cm.NextActionContinue, cm.ProcessStatOK
	logger.Info("Sysbackgroud for service running: %s", t.Name)

	for {
		select {
		case msg, ok := <-t.Cmd:
			if !ok {
				logger.Info("%s:Cmd Read error", t.Name)
				break
			}
			next, stat = t.SysControlEntry(t.Name, msg)
			if stat == cm.ProcessStatIgnore {
				next, _ = t.Instance.ControlHandler(msg)
			}
			break
		case msg, ok := <-t.ReadPipe():
			if !ok {
				logger.Info("%s:Data Read error", t.Name)
				break
			}
			t.Buffer.Append(msg)
			break
		}

		switch next {
		case cm.NextActionBreak:
			break
		case cm.NextActionReturn:
			return
		case cm.NextActionContinue:
		}
	}
	return

}

type IService interface {
	Start(bus *dm.DataMsgPipe) bool
	Pause() bool
	Resume() bool
	Exit() bool
	Self() *Service
	ControlHandler(*cm.ControlMsg) (int, int)
	DataHandler(*dm.DataMsg) bool
}
