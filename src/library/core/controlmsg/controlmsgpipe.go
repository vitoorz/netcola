package controlmsg

import (
	"library/logger"
)

type ControlMsgPipe struct {
	Cmd  chan *ControlMsg
	Echo chan *ControlMsg
}

func NewControlMsgPipe() *ControlMsgPipe {
	new := &ControlMsgPipe{}
	new.Cmd = make(chan *ControlMsg)
	new.Echo = make(chan *ControlMsg)
	return new
}

func (p *ControlMsgPipe) CloseControlMsgPipe() {
	close(p.Cmd)
	close(p.Echo)
}

func (p *ControlMsgPipe) ReadEchoNoBlock() (*ControlMsg, bool) {
	var msg *ControlMsg = nil
	select {
	case msg = <-p.Echo:
		return msg, true
	default:
	}
	return nil, false
}

func (p *ControlMsgPipe) WriteEchoNoBlock(msg *ControlMsg) bool {
	select {
	case p.Echo <- msg:
		return true
	default:
	}
	return false
}

func (p *ControlMsgPipe) ReadCmdNoBlock() (*ControlMsg, bool) {
	var msg *ControlMsg = nil
	select {
	case msg = <-p.Cmd:
		return msg, true
	default:
	}
	return nil, false
}

func (p *ControlMsgPipe) WriteCmdNoBlock(msg *ControlMsg) bool {
	select {
	case p.Cmd <- msg:
		return true
	default:
	}
	return false
}

func (t *ControlMsgPipe) SysControlEntry(serviceName string, msg *ControlMsg) (int, int) {
	switch msg.MsgType {
	case ControlMsgExit:
		logger.Info("%s:ControlMsgPipe.Cmd Read %d", serviceName, msg.MsgType)
		t.Echo <- &ControlMsg{MsgType: ControlMsgExit}
		logger.Info("%s:exit", serviceName)
		return NextActionReturn, ProcessStatOK
	case ControlMsgPause:
		logger.Info("%s:paused", serviceName)
		t.Echo <- &ControlMsg{MsgType: ControlMsgPause}
		for {
			var resume bool = false
			select {
			case msg, ok := <-t.Cmd:
				if !ok {
					logger.Info("%s:Cmd Read error", serviceName)
					break
				}
				switch msg.MsgType {
				case ControlMsgResume:
					t.Echo <- &ControlMsg{MsgType: ControlMsgResume}
					resume = true
					break
				}
			}
			if resume {
				break
			}
		}
		logger.Info("%s:resumed", serviceName)
	default:
		return NextActionContinue, ProcessStatIgnore
	}
	return NextActionContinue, ProcessStatOK
}
