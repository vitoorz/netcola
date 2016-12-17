package cs
//Auto generated, do not modify unless you know clearly what you are doing.

import (
	pb "github.com/golang/protobuf/proto"
	. "types"
    . "game/server/play"
)

func On_LoginReq(objectId IdString, opCode MsgType, payLoad []byte) interface{} {
	req := &LoginReq{}
	pb.Unmarshal(payLoad, req)
	return Handle_LoginReq(objectId, opCode, req)
}

func On_LoginAck(objectId IdString, opCode MsgType, payLoad []byte) interface{} {
	req := &LoginAck{}
	pb.Unmarshal(payLoad, req)
	return Handle_LoginAck(objectId, opCode, req)
}

func On_LogoutReq(objectId IdString, opCode MsgType, payLoad []byte) interface{} {
	req := &LogoutReq{}
	pb.Unmarshal(payLoad, req)
	return Handle_LogoutReq(objectId, opCode, req)
}

func On_LogoutAck(objectId IdString, opCode MsgType, payLoad []byte) interface{} {
	req := &LogoutAck{}
	pb.Unmarshal(payLoad, req)
	return Handle_LogoutAck(objectId, opCode, req)
}
