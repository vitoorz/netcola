package datamsg

import (
	"library/idgen"
)

type DataMsg struct {
	opcode  int
	owner   idgen.ObjectID
	payload interface{}
	next    *DataMsg
}

func NewDataMsg() *DataMsg {
	var msg = &DataMsg{}
	return msg
}

func (p *DataMsg) SetOpCode(opcode int) {
	p.opcode = opcode
}

func (p *DataMsg) OpCode() int {
	return p.opcode
}

func (p *DataMsg) SetOwner(id idgen.ObjectID) {
	p.owner = id
}

func (p *DataMsg) Owner() idgen.ObjectID {
	return p.owner
}

func (p *DataMsg) SetPayload(pkt interface{}) {
	p.payload = pkt
}

func (p *DataMsg) Payload() interface{} {
	return p.payload
}

func (p *DataMsg) PushBack(n *DataMsg) (h *DataMsg) {
	if p == nil {
		return n
	}
	m := p
	for ; m.next != nil; m = m.next {
	}
	m.next = n
	return p
}

func (p *DataMsg) PushFront(n *DataMsg) (h *DataMsg) {
	if n == nil {
		return p
	}
	n.next = p
	return n
}
