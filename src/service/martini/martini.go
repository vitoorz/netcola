package martini

import (
	dm "library/core/datamsg"
	"service"
)

const ServiceName = "martini"

type martiniType struct {
	*service.Service
	Output *dm.DataMsgPipe
}

func NewMartini(name string) *martiniType {
	t := &martiniType{}
	t.Service = service.NewService(ServiceName)
	t.Output = nil
	t.Name = name
	t.Instance = t
	return t
}
