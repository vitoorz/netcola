package service

import (
	"errors"
	"fmt"
)

var ServicePool = &ServiceManager{
	services: make(map[string]ServiceI, 0),
}

func StartService(s ServiceI) {
	s.Init()
	s.Start()
	ServicePool.Register(s)
}

type ServiceManager struct {
	services map[string]ServiceI
}

func (t *ServiceManager) Service(name string) (ServiceI, bool) {
	s, ok := t.services[name]
	return s, ok
}

func (t *ServiceManager) Register(service ServiceI) error {
	c := service.CommonService()
	_, ok := t.services[c.Name]
	if ok {
		return errors.New(fmt.Sprint("Already exist service %s", c.Name))
	}
	t.services[c.Name] = service

	return nil
}
