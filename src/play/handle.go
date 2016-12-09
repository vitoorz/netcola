package play
//Auto generated, do not modify unless you know clearly what you are doing.

import (
	pb "github.com/golang/protobuf/proto"
	. "types"
)


func On_LoginReq(playerId IdString, opCode MsgType, payLoad []byte) interface{} {
	req :=&LoginReq{}
	pb.Unmarshal(payLoad, req)
	return Play_LoginReq(playerId, opCode, req)
}

func On_LoginAck(playerId IdString, opCode MsgType, payLoad []byte) interface{} {
	return Play_InvalidReq(playerId, opCode, payLoad)
}

func On_LogoutReq(playerId IdString, opCode MsgType, payLoad []byte) interface{} {
	req :=&LogoutReq{}
	pb.Unmarshal(payLoad, req)
	return Play_LogoutReq(playerId, opCode, req)
}

func On_LogoutAck(playerId IdString, opCode MsgType, payLoad []byte) interface{} {
	return Play_InvalidReq(playerId, opCode, payLoad)
}
