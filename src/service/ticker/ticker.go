package ticker

import (
	dm "library/core/datamsg"
	"library/idgen"
	"service"
	"time"
)

const ServiceName = "ticker"

type tickerType struct {
	service.Service
	output       *dm.DataMsgPipe
	tickerBooker map[idgen.ObjectID]*ticker //ID, should have a method to find service by unique id
}

type ticker struct {
	TickStep int64
	TickAt   int64
}

func NewTicker(name string) *tickerType {
	t := &tickerType{}
	t.Service = *service.NewService(ServiceName)
	t.output = nil
	t.Name = name
	t.tickerBooker = make(map[idgen.ObjectID]*ticker)
	t.Instance = t
	return t
}

func (t *tickerType) BookedBy(s service.IService, stepSecond int64) {
	t.tickerBooker[s.Self().ID] = &ticker{stepSecond, time.Now().Unix() + stepSecond}
}
