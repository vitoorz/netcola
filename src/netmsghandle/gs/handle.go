package gs
//Auto generated, do not modify unless you know clearly what you are doing.

import (
	pb "github.com/golang/protobuf/proto"
	. "types"
    . "game/server/gs"
    . "game/gateway/gm"
)

func On_ServerLoginReq(objectId IdString, opCode MsgType, payLoad []byte) interface{} {
	req := &ServerLoginReq{}
	pb.Unmarshal(payLoad, req)
	return Handle_ServerLoginReq(objectId, opCode, req)
}

func On_ServerLoginAck(objectId IdString, opCode MsgType, payLoad []byte) interface{} {
	req := &ServerLoginAck{}
	pb.Unmarshal(payLoad, req)
	return Handle_ServerLoginAck(objectId, opCode, req)
}

func On_ServerLoginOutReq(objectId IdString, opCode MsgType, payLoad []byte) interface{} {
	req := &ServerLoginOutReq{}
	pb.Unmarshal(payLoad, req)
	return Handle_ServerLoginOutReq(objectId, opCode, req)
}

func On_BrdCastAddMemberReq(objectId IdString, opCode MsgType, payLoad []byte) interface{} {
	req := &BrdCastGroupManageReq{}
	pb.Unmarshal(payLoad, req)
	return Handle_BrdCastAddMemberReq(objectId, opCode, req)
}

func On_BrdCastAddMemberAck(objectId IdString, opCode MsgType, payLoad []byte) interface{} {
	req := &BrdCastGroupManageAck{}
	pb.Unmarshal(payLoad, req)
	return Handle_BrdCastAddMemberAck(objectId, opCode, req)
}

func On_BrdCastDelMemberReq(objectId IdString, opCode MsgType, payLoad []byte) interface{} {
	req := &BrdCastGroupManageReq{}
	pb.Unmarshal(payLoad, req)
	return Handle_BrdCastDelMemberReq(objectId, opCode, req)
}

func On_BrdCastDelMemberAck(objectId IdString, opCode MsgType, payLoad []byte) interface{} {
	req := &BrdCastGroupManageAck{}
	pb.Unmarshal(payLoad, req)
	return Handle_BrdCastDelMemberAck(objectId, opCode, req)
}

func On_BrdCastDestroyReq(objectId IdString, opCode MsgType, payLoad []byte) interface{} {
	req := &BrdCastGroupManageReq{}
	pb.Unmarshal(payLoad, req)
	return Handle_BrdCastDestroyReq(objectId, opCode, req)
}

func On_BrdCastDestroyAck(objectId IdString, opCode MsgType, payLoad []byte) interface{} {
	req := &BrdCastGroupManageAck{}
	pb.Unmarshal(payLoad, req)
	return Handle_BrdCastDestroyAck(objectId, opCode, req)
}

func On_BrdCastSyncReq(objectId IdString, opCode MsgType, payLoad []byte) interface{} {
	req := &BrdCastGroupManageReq{}
	pb.Unmarshal(payLoad, req)
	return Handle_BrdCastSyncReq(objectId, opCode, req)
}

func On_BrdCastSyncAck(objectId IdString, opCode MsgType, payLoad []byte) interface{} {
	req := &BrdCastGroupManageAck{}
	pb.Unmarshal(payLoad, req)
	return Handle_BrdCastSyncAck(objectId, opCode, req)
}
