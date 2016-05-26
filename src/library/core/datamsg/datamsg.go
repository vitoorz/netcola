package datamsg

type DataMsg struct {
	Receiver string
	Meta     interface{}
	MsgType  int
	Next     *DataMsg
	Payload  interface{}
}

func NewDataMsg(recv string, meta interface{},
	msgtype int, payload interface{}) *DataMsg {
	var msg = &DataMsg{
		MsgType:  msgtype,
		Receiver: recv,
		Meta:     meta,
		Next:     nil,
		Payload:  payload,
	}
	return msg
}

func (p *DataMsg) PushBack(n *DataMsg) (h *DataMsg) {
	if p == nil {
		return n
	}
	m := p
	for ; m.Next != nil; m = m.Next {
	}
	m.Next = n
	return p
}

func (p *DataMsg) PushFront(n *DataMsg) (h *DataMsg) {
	if n == nil {
		return p
	}
	n.Next = p
	return n
}
