package datamsg

type DataMsg struct {
	Receiver string
	meta     map[string]interface{}
	MsgType  int
	Next     *DataMsg
	Payload  interface{}
}

func NewDataMsg(recv string, msgtype int, payload interface{}) *DataMsg {
	var msg = &DataMsg{
		MsgType:  msgtype,
		Receiver: recv,
		meta:     make(map[string]interface{}),
		Next:     nil,
		Payload:  payload,
	}
	return msg
}

func (t *DataMsg) SetMeta(owner string, m interface{}) {
	t.meta[owner] = m
}

func (t *DataMsg) Meta(owner string) (interface{}, bool) {
	i, ok := t.meta[owner]
	return i, ok
}

func (t *DataMsg) PushBack(n *DataMsg) (d *DataMsg) {
	if t == nil {
		return n
	}
	m := t
	for ; m.Next != nil; m = m.Next {
	}
	m.Next = n
	return t
}

func (t *DataMsg) PushFront(n *DataMsg) (d *DataMsg) {
	if n == nil {
		return t
	}
	n.Next = t
	return n
}
