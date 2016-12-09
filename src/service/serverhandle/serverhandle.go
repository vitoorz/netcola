package serverhandle

import (
	dm "library/core/datamsg"
	"service"
)

const ServiceName = "serverhandle"

type serverHandle struct {
	*service.Service
	Output *dm.DataMsgPipe
}

func NewServerHandle(name string) *serverHandle {
	t := &serverHandle{}
	t.Service = service.NewService(ServiceName)
	t.Output = nil
	t.Name = name
	t.Instance = t
	return t
}
