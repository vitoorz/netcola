package ticker

import (
	cm "library/core/controlmsg"
	dm "library/core/datamsg"
	"service"
	"time"
)

func (t *tickerType) Start(bus *dm.DataMsgPipe) bool {
	t.output = bus
	go t.tick()
	return true
}

func (t *tickerType) Pause() bool {
	return true
}

func (t *tickerType) Resume() bool {
	return true
}

func (t *tickerType) Exit() bool {
	return true
}

func (t *tickerType) ControlHandler(msg *cm.ControlMsg) (int, int) {
	if msg.MsgType == cm.ControlMsgTick {
		now := time.Now().Unix()
		for id, ticker := range t.tickerBooker {
			if ticker.TickAt >= now {
				msg := &cm.ControlMsg{MsgType: cm.ControlMsgTick, Payload: nil}
				s, ok := service.ServicePool.ServiceById(id)
				if ok && s.Self().WriteCmdNonblock(msg) {
					ticker.TickAt = now + ticker.TickStep
				}
			}
		}
	}
	return cm.NextActionContinue, cm.ProcessStatOK
}
