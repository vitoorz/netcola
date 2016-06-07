package mongo

import (
	"gopkg.in/mgo.v2"
	"time"
)

import (
	cm "library/core/controlmsg"
	dm "library/core/datamsg"
	"library/logger"
)

const (
	MAX_SLEEP_SECONDS = 10
)

func (t *mongoType) Start(bus *dm.DataMsgPipe) bool {
	t.output = bus
	return t.dial()
}

func (t *mongoType) Pause() bool {
	return true
}

func (t *mongoType) Resume() bool {
	return true
}

func (t *mongoType) Exit() bool {
	return true
}

func (t *mongoType) ControlHandler(msg *cm.ControlMsg) (int, int) {
	//switch msg.MsgType {
	//case cm.ControlMsgExit:
	//	logger.Info("%s:ControlMsgPipe.Cmd Read %d", t.Name, msg.MsgType)
	//	t.Echo <- &cm.ControlMsg{MsgType: cm.ControlMsgExit}
	//	logger.Info("%s:exit", t.Name)
	//	return cm.NextActionReturn, cm.ProcessStatOK
	//case cm.ControlMsgPause:
	//	logger.Info("%s:paused", t.Name)
	//	t.Echo <- &cm.ControlMsg{MsgType: cm.ControlMsgPause}
	//	for {
	//		var resume bool = false
	//		select {
	//		case msg, ok := <-t.Cmd:
	//			if !ok {
	//				logger.Info("%s:Cmd Read error", t.Name)
	//				break
	//			}
	//			switch msg.MsgType {
	//			case cm.ControlMsgResume:
	//				t.Echo <- &cm.ControlMsg{MsgType: cm.ControlMsgResume}
	//				resume = true
	//				break
	//			}
	//		}
	//		if resume {
	//			break
	//		}
	//	}
	//	logger.Info("%s:resumed", t.Name)
	//}
	return cm.NextActionContinue, cm.ProcessStatOK
}

func (t *mongoType) dial() bool {
	var keepTrying bool = true
	logger.Info("%s:connecting:%s:%s", t.Name, t.ip, t.port)
	if t.session != nil {
		t.session.Close()
	}

	for stepping := 1; keepTrying; stepping += 1 {
		sess, err := mgo.Dial(t.ip + ":" + t.port)
		if err != nil {
			logger.Error("%s:connect fail,err:%s", t.Name, err.Error())
			if stepping > 5 {
				return false
			}
			time.Sleep(time.Second * time.Duration(stepping*2%MAX_SLEEP_SECONDS))
		} else {
			t.session = sess
			break
		}
	}
	logger.Info("%s:connection established,session:%p", t.Name, t.session)
	t.session.SetMode(mgo.Strong, true)
	return true
}
