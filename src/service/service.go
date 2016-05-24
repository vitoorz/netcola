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
	BUS   *dm.DataMsgPipe
	ID    idgen.ObjectID
	State StateT
	Name  string
}

func NewService(name string) *Service {
	t := &Service{}
	t.ControlMsgPipe = *cm.NewControlMsgPipe()
	t.BUS = nil
	t.ID = idgen.NewObjectID()
	t.Name = name
	return t
}

type ServiceI interface {
	Init() bool
	Start() bool
	Pause() bool
	Exit() bool
	Kill() bool
}

func StartService(s ServiceI) {
	s.Init()
	s.Start()
}
