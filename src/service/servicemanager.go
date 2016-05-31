package service

import (
	"errors"
	"fmt"
	dm "library/core/datamsg"
	. "library/idgen"
)

//the only service manager instance in process, will maintain all service instance
var ServicePool = &ServiceManager{
	services: make(map[string]*serviceSeries, 0),
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

func (t *ServiceManager) Register(service IService) error {
	s := service.Self()
	sList, ok := t.services[s.Name]
	if !ok {
		sList = newServiceSeries()
		t.services[s.Name] = sList
	}

	sList.Lists[s.ID] = service
	return nil
}

func StartService(s IService, bus *dm.DataMsgPipe) bool {
	if s.Start(bus) {
		ServicePool.Register(s)
		return true
	} else {
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
