package mongo

import (
//"time"
)

import (
	dm "library/core/datamsg"
	"library/logger"
)

func (t *mongoType) DataHandler(msg *dm.DataMsg) bool {
	//logger.Info("this is handle:%+v", msg)
	i, ok := msg.Meta(t.Name)
	if !ok {
		logger.Error("get Dirty interface error")
		return false
	}

	d, ok := i.(Dirty)
	if !ok {
		logger.Error("msg is not Dirty interface")
		return false
	}
	for ok = d.CRUD(t.session); !ok; {
		logger.Info("CRUD failed:%s", d.Inspect())
		ok = d.CRUD(t.session)
	}

	msg.Receiver = msg.Sender
	msg.Sender = t.Name
	t.output.WritePipeNoBlock(msg)
	return true
}
