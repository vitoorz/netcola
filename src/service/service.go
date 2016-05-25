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
	cm.ControlMsgPipe
	dm.DataMsgPipe
	BUS   *dm.DataMsgPipe
	ID    idgen.ObjectID
	State StateT
	Name  string
}

func NewService(name string) *Service {
	t := &Service{}
	t.ControlMsgPipe = *cm.NewControlMsgPipe()
	t.ID = idgen.NewObjectID()
	t.Name = name
	return t
}

func (t *Service) Self() *Service {
	return t
}

type IService interface {
	Init() bool
	Start() bool
	Pause() bool
	Exit() bool
	Kill() bool
	Self() *Service
}
