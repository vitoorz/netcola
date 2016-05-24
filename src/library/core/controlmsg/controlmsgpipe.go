package controlmsg

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

func (p *ControlMsgPipe) ReadEchoNonblock() (*ControlMsg, bool) {
	var msg *ControlMsg = nil
	select {
	case msg = <-p.Echo:
		return msg, true
	default:
	}
	return nil, false
}

func (p *ControlMsgPipe) WriteEchoNonblock(msg *ControlMsg) bool {
	select {
	case p.Echo <- msg:
		return true
	default:
	}
	return false
}

func (p *ControlMsgPipe) ReadCmdNonblock() (*ControlMsg, bool) {
	var msg *ControlMsg = nil
	select {
	case msg = <-p.Cmd:
		return msg, true
	default:
	}
	return nil, false
}

func (p *ControlMsgPipe) WriteCmdNonblock(msg *ControlMsg) bool {
	select {
	case p.Cmd <- msg:
		return true
	default:
	}
	return false
}
