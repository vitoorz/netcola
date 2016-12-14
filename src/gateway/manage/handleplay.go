package manage

//Auto generated, do not modify unless you know clearly what you are doing.

import (
	. "types"
)

func Handle_InvalidReq(objectId IdString, opCode MsgType, req interface{}) interface{} {
	return nil
}

func Handle_ServerLoginReq(objectId IdString, opCode MsgType, req *ServerLoginReq) interface{} {
	ack := &ServerLoginAck{ServerId: "0x84a5d16000028001"}
	return ack
}
