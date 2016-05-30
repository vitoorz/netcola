package datamsg

import (
	"library/idgen"
)

type DataMsg struct {
	Receiver string
	Meta     map[idgen.ObjectID]interface{}
	MsgType  int
	Next     *DataMsg
	Payload  interface{}
}

func NewDataMsg(recv string, msgtype int, payload interface{}) *DataMsg {
	var msg = &DataMsg{
		MsgType:  msgtype,
		Receiver: recv,
		Meta:     make(map[idgen.ObjectID]interface{}),
		Next:     nil,
		Payload:  payload,
	}
	return msg
}

func (t *DataMsg) SetMeta(id idgen.ObjectID, meta interface{}) {
	t.Meta[id] = meta
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
