package service

import (
	"errors"
	"fmt"
)

var ServicePool = &ServiceManager{
	services: make(map[string]IService, 0),
}

func StartService(s IService) {
	s.Init()
	s.Start()
	ServicePool.Register(s)
}

type ServiceManager struct {
	services map[string]IService
}

func (t *ServiceManager) Service(name string) (IService, bool) {
	s, ok := t.services[name]
	return s, ok
}

func (t *ServiceManager) Register(service IService) error {
	c := service.Self()
	_, ok := t.services[c.Name]
	if ok {
		return errors.New(fmt.Sprint("Already exist service %s", c.Name))
	}
	t.services[c.Name] = service

	return nil
}
