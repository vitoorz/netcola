package service

import (
	dm "library/core/datamsg"
	. "library/idgen"
	"library/logger"
)

//the only service manager instance in process, will maintain all service instance
var ServicePool = &ServiceManager{
	services: make(map[string]*serviceSeries, 0), //key:name,value:service objects
	byId:     make(map[ObjectID]IService, 0),     // index service by service id
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
		//todo: maybe need a load balancer here
		for _, s := range sList.Lists {
			return s, true
		}
	}

	return nil, ok
}

// get service object by id index
func (t *ServiceManager) ServiceById(id ObjectID) (IService, bool) {
	s, ok := t.byId[id]
	return s, ok
}

// return all the service name
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
	instance := s.Self()
	if s.Start(bus) {
		ServicePool.register(s)
		go instance.Background()
		go instance.Buffer.Daemon()
		return true
	} else {
		logger.Error("start service %s failed", instance.Name)
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
