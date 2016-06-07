package ticker

import (
	cm "library/core/controlmsg"
	dm "library/core/datamsg"
	"time"
)

func (t *tickerType) tick() {
	tickChan := time.NewTicker(time.Second).C
	for {
		select {
		case <-tickChan:
			t.WriteCmdNonblock(&cm.ControlMsg{MsgType: cm.ControlMsgTick, Payload: nil})
			break
		}
	}
	return
}

func (t *tickerType) DataHandler(msg *dm.DataMsg) bool {
	return true
}
