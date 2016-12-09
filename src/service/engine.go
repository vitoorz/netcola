package service

import (
	cm "library/core/controlmsg"
	dm "library/core/datamsg"
	"library/logger"
)

type engineType struct {
	*cm.ControlMsgPipe
	*dm.DataMsgPipe
	BUS  *dm.DataMsgPipe
	Name string
}

func NewEngine(name string) *engineType {
	t := &engineType{}
	t.ControlMsgPipe = cm.NewControlMsgPipe()
	t.DataMsgPipe = dm.NewDataMsgPipe(0)
	t.Name = name
	t.BUS = dm.NewDataMsgPipe(0)
	return t
}

func (t *engineType) engine() {
	logger.Info("%s:running", t.Name)
	var next, fun int = cm.NextActionContinue, cm.ProcessStatUnknown
	for {
		select {
		case msg, ok := <-t.Cmd:
			if !ok {
				logger.Info("%s:Cmd Read error", t.Name)
				break
			}
			next, fun = t.ControlHandler(msg)
			break
		case msg, ok := <-t.ReadPipe():
			if !ok {
				logger.Info("%s:DataPipe Read error", t.Name)
				break
			}
			next, fun = t.DataEntry(msg)
			break
		case msg, ok := <-t.BUS.ReadPipe():
			if !ok {
				logger.Info("%s:DataPipe Read error", t.Name)
				break
			}
			//logger.Debug("%s:read bus chan msg:%+v", t.Name, msg)
			next, fun = t.BusSchedule(msg)
			if fun == cm.ProcessPipeReceiverLost {
				logger.Warn("%s:reciver lost:%v", t.Name, msg.Receiver)
			}
			break
		}

		switch next {
		case cm.NextActionBreak:
			break
		case cm.NextActionReturn:
			return
		case cm.NextActionContinue:
		}
	}
	return
}

func (t *engineType) BusSchedule(msg *dm.DataMsg) (operate int, funCode int) {
	funCode = cm.ProcessStatOK
	ok := ServicePool.SendData(msg)
	if !ok {
		funCode = cm.ProcessPipeReceiverLost
	}
	return cm.NextActionContinue, funCode
}

func (t *engineType) Start(bus *dm.DataMsgPipe) bool {
	go t.engine()
	return true
}

func (t *engineType) Resume() bool {
	logger.Info("%s:will resume", t.Name)
	t.Cmd <- &cm.ControlMsg{MsgType: cm.ControlMsgResume}
	echo := <-t.Echo
	if echo.MsgType != cm.ControlMsgResume {
		return false
	}
	return true
}

func (t *engineType) ControlHandler(msg *cm.ControlMsg) (int, int) {
	switch msg.MsgType {
	case cm.ControlMsgExit:
		logger.Info("%s:ControlMsgPipe.Cmd Read %d", t.Name, msg.MsgType)
		t.Echo <- &cm.ControlMsg{MsgType: cm.ControlMsgExit}
		logger.Info("%s:exit", t.Name)
		return cm.NextActionReturn, cm.ProcessStatOK
	case cm.ControlMsgPause:
		logger.Info("%s:paused", t.Name)
		t.Echo <- &cm.ControlMsg{MsgType: cm.ControlMsgPause}
		for {
			var resume bool = false
			select {
			case msg, ok := <-t.Cmd:
				if !ok {
					logger.Info("%s:Cmd Read error", t.Name)
					break
				}
				switch msg.MsgType {
				case cm.ControlMsgResume:
					t.Echo <- &cm.ControlMsg{MsgType: cm.ControlMsgResume}
					resume = true
					break
				}
			}
			if resume {
				break
			}
		}
		logger.Info("%s:resumed", t.Name)
	}
	return cm.NextActionContinue, cm.ProcessStatOK
}

func (t *engineType) DataEntry(msg *dm.DataMsg) (operate int, funCode int) {
	return cm.NextActionContinue, cm.ProcessStatOK
}
