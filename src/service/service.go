package service

import (
	"sync"
)

import (
	cm "library/core/controlmsg"
	dm "library/core/datamsg"
	"library/idgen"
)

type Service struct {
	sync.RWMutex
	ID    idgen.ObjectID
	Name  string
	State StateT
	cm.ControlMsgPipe
	dm.DataMsgPipe
	BUS *dm.DataMsgPipe
}

func NewService(name string) *Service {
	t := &Service{}
	t.ID = idgen.NewObjectID()
	t.Name = name
	t.State = StateInit
	t.ControlMsgPipe = *cm.NewControlMsgPipe()
	t.DataMsgPipe = *dm.NewDataMsgPipe(0)
	t.BUS = nil
	return t
}

func (t *Service) Self() *Service {
	return t
}

type IService interface {
	Start(name string, bus *dm.DataMsgPipe) bool
	Pause() bool
	Resume() bool
	Exit() bool
	Self() *Service
}
