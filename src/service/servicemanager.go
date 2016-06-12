package service

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
	logger.Info("start service:host:%v,instance:%+v", s.Self().Name, s.Self().Instance)
	parent := s.Self()
	if s.Start(bus) {
		ServicePool.register(s)
		go parent.Background()
		go parent.Buffer.Daemon()
		return true
	} else {
		logger.Error("start service %s failed", parent.Name)
		return false
	}
}

//send down (mostly request) data message to the right receiver
func (t *ServiceManager) SendData(msg *dm.DataMsg) bool {
	s, ok := t.Service(msg.Receiver)
	if !ok {
		logger.Error("Service %s to receive message not exist:%v", msg.Receiver)
		return ok
	}
	s.Self().DataMsgPipe.WritePipeNoBlock(msg)
	return true
}
