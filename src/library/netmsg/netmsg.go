package netmsg

import (
	"library/idgen"
)

type NetMsg struct {
	opcode  int
	owner   idgen.ObjectID
	payload interface{}
	next    *NetMsg
}

func NewNetMsg() *NetMsg {
	var msg = &NetMsg{}
	return msg
}

func (p *NetMsg) SetOpCode(opcode int) {
	p.opcode = opcode
}

func (p *NetMsg) OpCode() int {
	return p.opcode
}

func (p *NetMsg) SetOwner(id idgen.ObjectID) {
	p.owner = id
}

func (p *NetMsg) Owner() idgen.ObjectID {
	return p.owner
}

func (p *NetMsg) SetPayload(pkt interface{}) {
	p.payload = pkt
}

func (p *NetMsg) Payload() interface{} {
	return p.payload
}

func (p *NetMsg) PushBack(n *NetMsg) (h *NetMsg) {
	if p == nil {
		return n
	}
	m := p
	for ; m.next != nil; m = m.next {
	}
	m.next = n
	return p
}

func (p *NetMsg) PushFront(n *NetMsg) (h *NetMsg) {
	if n == nil {
		return p
	}
	n.next = p
	return n
}
