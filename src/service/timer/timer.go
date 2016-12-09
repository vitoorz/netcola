package timer

import (
	dm "library/core/datamsg"
	"service"
	. "types"
)

const ServiceName = "timer"

type timerType struct {
	*service.Service
	output    *dm.DataMsgPipe
	timerPool map[ObjectID]*dm.DataMsg
}

func NewTimer(name string) *timerType {
	t := &timerType{}
	t.Service = service.NewService(ServiceName)
	t.output = nil
	t.Name = name
	t.timerPool = make(map[ObjectID]*dm.DataMsg)
	t.Instance = t
	return t
}
