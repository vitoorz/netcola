package mongo

import (
	"gopkg.in/mgo.v2"
)

import (
	dm "library/core/datamsg"
	"service"
)

const ServiceName = "mongo"

type mongoType struct {
	*service.Service
	output  *dm.DataMsgPipe
	ip      string
	port    string
	session *mgo.Session
}

func NewMongo(name, ip, port string) *mongoType {
	t := &mongoType{}
	t.Service = service.NewService(ServiceName)
	t.Name = name
	t.output = nil
	t.ip = ip
	t.port = port
	t.session = nil
	t.Instance = t
	return t
}

type Dirty interface {
	CRUD(interface{}) bool
	Inspect() string
}
