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

const (
	Break = iota
	Continue
	Return
)

type Service struct {
	sync.RWMutex
	ID    idgen.ObjectID
	Name  string
	State StateT
	cm.ControlMsgPipe
	dm.DataMsgPipe
	Buffer   *BufferPool
	Instance IService
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
	logger.Info("%s:service running", t.Name)
	var next int = Continue
	for {
		select {
		case msg, ok := <-t.Cmd:
			if !ok {
				logger.Info("%s:Cmd Read error", t.Name)
				break
			}
			next, _ = t.Instance.ControlHandler(msg)
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
		case Break:
			break
		case Return:
			return
		case Continue:
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
