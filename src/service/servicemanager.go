package service

import (
	"errors"
	"fmt"
)

import (
	dm "library/core/datamsg"
	. "library/idgen"
	"library/logger"
)

//the only service manager instance in process, will maintain all service instance
var ServicePool = &ServiceManager{
	services: make(map[string]*serviceSeries, 0),
	byId:     make(map[ObjectID]IService, 0),
}

//Service with the same name, will have the same function
type serviceSeries struct {
	Lists map[ObjectID]IService
}

func newServiceSeries() *serviceSeries {
	return &serviceSeries{
		Lists: make(map[ObjectID]IService, 0),
	}
}

type ServiceManager struct {
	services map[string]*serviceSeries
	byId     map[ObjectID]IService
}

func (t *ServiceManager) Service(name string) (IService, bool) {
	//should have strategy to sample from multiple instance
	sList, ok := t.services[name]
	if ok {
		for _, s := range sList.Lists {
			return s, true
		}
	}

	return nil, ok
}

func (t *ServiceManager) ServiceById(id ObjectID) (IService, bool) {
	s, ok := t.byId[id]
	return s, ok
}

func (t *ServiceManager) NameList() (list []string) {
	for name := range t.services {
		list = append(list, name)
	}
	return list
}

func (t *ServiceManager) register(service IService) error {
	s := service.Self()
	sList, ok := t.services[s.Name]
	if !ok {
		sList = newServiceSeries()
		t.services[s.Name] = sList
	}

	sList.Lists[s.ID] = service
	t.byId[s.ID] = service

	return nil
}

func StartService(s IService, bus *dm.DataMsgPipe) bool {
	if s.Start(bus) {
		instance := s.Self()
		if instance.Buffer != nil {
			go instance.Buffer.Daemon()
		}
		go instance.Background()

		ServicePool.register(s)
		return true
	} else {
		logger.Error("start %s failed", s.Self().Name)
		return false
	}
}

//send down (mostly request) data message to the right receiver
func (t *ServiceManager) SendData(msg *dm.DataMsg) error {
	s, ok := t.Service(msg.Receiver)
	if !ok {
		return errors.New(fmt.Sprint("Service %s to receive message not exist", msg.Receiver))
	}
	s.Self().DataMsgPipe.WritePipeNB(msg)

	return nil
}
