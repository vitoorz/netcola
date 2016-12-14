package manage

//Auto generated, do not modify unless you know clearly what you are doing.

import (
	pb "github.com/golang/protobuf/proto"
	. "types"
)

func On_ServerLoginReq(objectId IdString, opCode MsgType, payLoad []byte) interface{} {
	req := &ServerLoginReq{}
	pb.Unmarshal(payLoad, req)
	return Handle_ServerLoginReq(objectId, opCode, req)
}

func On_ServerLoginAck(objectId IdString, opCode MsgType, payLoad []byte) interface{} {
	return Handle_InvalidReq(objectId, opCode, payLoad)
}
