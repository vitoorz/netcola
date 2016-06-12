package syswatcher

import (
	cm "library/core/controlmsg"
	"library/idgen"
	"library/logger"
	"service"
	"time"
)

//singleton service routine to provide both tick service, and provide watch time consumption of procedure object
//a system service control message of type cm.ControlMsgTick is send to the booker on each tick point

const (
	watchCmdStart = 1
	watchCmdEnd   = 2
)

type tickBooker struct {
	TickStep int64
	TickAt   int64
	Service  service.IService
}

type watcherType struct {
	tickerBooker map[idgen.ObjectID]*tickBooker //ID, should have a method to find service by unique id
	watchObject  map[string]int64               //key: object unique identifier, value: start watch timestamp
	cmd          chan *watchCmd
}

type watchCmd struct {
	cmdType uint8
	object  string
}

func NewWatcher() *watcherType {
	return &watcherType{
		tickerBooker: make(map[idgen.ObjectID]*tickBooker),
		watchObject:  make(map[string]int64, 0),
		cmd:          make(chan *watchCmd, 0),
	}
}

func (t *watcherType) Watching() {
	go t.watcher()
}

func (t *watcherType) BookedBy(s service.IService, stepSecond int64) {
	t.tickerBooker[s.Self().ID] = &tickBooker{stepSecond, time.Now().Unix() + stepSecond, s}
}

func (t *watcherType) WatchObjectStart(object string) {
	t.cmd <- &watchCmd{watchCmdStart, object}
}

func (t *watcherType) WatchObjectEnd(object string) {
	t.cmd <- &watchCmd{watchCmdEnd, object}
}

func (t *watcherType) watcher() {
	tickChan := time.NewTicker(time.Second).C
	for {
		select {
		case <-tickChan:
			t.onTick(time.Now().Unix())
		case cmd := <-t.cmd:
			t.onCmd(cmd)
		}
	}
}

func (t *watcherType) doTickNotify(curTime int64) {
	for _, ticker := range t.tickerBooker {
		if ticker.TickAt >= curTime {
			msg := &cm.ControlMsg{MsgType: cm.ControlMsgTick, Payload: nil}
			if ticker.Service.Self().WriteCmdNoBlock(msg) {
				ticker.TickAt = curTime + ticker.TickStep
			}
		}
	}
}

func (t *watcherType) onTick(curTime int64) {
	t.doTickNotify(curTime)

	for obj, startTime := range t.watchObject {
		past := curTime - startTime
		logger.Warn("watched object %s for %d seconds", obj, past)
		if past >= 2 {
			logger.PProf()
			delete(t.watchObject, obj)
		}
	}
}

func (t *watcherType) onCmd(cmd *watchCmd) {
	switch cmd.cmdType {
	case watchCmdStart:
		t.watchObject[cmd.object] = time.Now().Unix()
	case watchCmdEnd:
		delete(t.watchObject, cmd.object)
	}
}
