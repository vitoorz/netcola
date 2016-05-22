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

func (n *ControlMsgPipe) CloseControlMsgPipe() {
	close(n.Cmd)
	close(n.Echo)
}
