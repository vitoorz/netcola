package mongo

import (
	"gopkg.in/mgo.v2"
	dm "library/core/datamsg"
	"library/logger"
	"service"
)

const ServiceName = "mongo"

const (
	Break = iota
	Continue
	Return
)

type mongoType struct {
	service.Service
	output  *dm.DataMsgPipe
	ip      string
	port    string
	session *mgo.Session
	buffer  *service.BufferPool
}

func NewMongo(name, ip, port string) *mongoType {
	t := &mongoType{}
	t.Service = *service.NewService(ServiceName)
	t.Name = name
	t.output = nil
	t.ip = ip
	t.port = port
	t.session = nil
	t.buffer = service.NewBufferPool(t)
	return t
}

func (t *mongoType) mongo() {
	logger.Info("%s:service running", t.Name)
	var next, fun int = Continue, service.FunUnknown

	// todo: how to control this daemon?
	//go t.buffer.Daemon()
	for {
		select {
		case msg, ok := <-t.Cmd:
			if !ok {
				logger.Info("%s:Cmd Read error", t.Name)
				break
			}
			next, fun = t.controlEntry(msg)
			break
		case msg, ok := <-t.ReadPipe():
			if !ok {
				logger.Info("%s:Data Read error", t.Name)
				break
			}

			next, fun = t.dataEntry(msg)
			if fun == service.FunDataPipeFull {
				logger.Warn("%s:need do something when full", t.Name)
			}
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

func (t *mongoType) Handle(msg *dm.DataMsg) bool {
	logger.Info("this is handle:%+v", msg)
	i, ok := msg.Meta(t.Name)
	if !ok {
		logger.Error("get Dirty interface error")
		return false
	}

	d, ok := i.(Dirty)
	if !ok {
		logger.Error("msg is not Dirty interface")
		return false
	}
	for ok = d.CRUD(t.session); !ok; {
		logger.Info("CRUD failed:%s", d.Inspect())
		ok = d.CRUD(t.session)
	}
	return true
}
