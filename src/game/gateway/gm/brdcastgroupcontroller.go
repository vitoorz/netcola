package gm

import (
	. "game/com"
	"library/logger"
	. "types"
)

func Handle_BrdCastAddMemberReq(objectId IdString, opCode MsgType, req *BrdCastGroupManageReq) interface{} {
	ack := &BrdCastGroupManageAck{GroupId: req.GroupId}

	grp, ok := FindBrdCastGroup(objectId, IdString(req.GroupId))
	if !ok {
		grp = NewBrdCastGroup4Server(objectId, IdString(req.GroupId))
	}

	ack.MemberNum = grp.AddMembers(req.MemberIds...)

	ack.Common = serverCommonAck(OK)

	return ack
}

func Handle_BrdCastDelMemberReq(objectId IdString, opCode MsgType, req *BrdCastGroupManageReq) interface{} {
	ack := &BrdCastGroupManageAck{GroupId: req.GroupId}

	grp, ok := FindBrdCastGroup(objectId, IdString(req.GroupId))
	if !ok {
		ack.Common = serverCommonAck(ERR_GROUP_NOT_FOUND)
		return ack
	}

	for _, memberId := range req.MemberIds {
		ack.MemberNum, ok = grp.DelMember(IdString(memberId))
		if !ok {
			logger.Warn("broadcsat group %s delete member %s error", req.GroupId, memberId)
		}
	}

	ack.Common = serverCommonAck(OK)

	return ack
}

func Handle_BrdCastDestroyReq(objectId IdString, opCode MsgType, req *BrdCastGroupManageReq) interface{} {
	ack := &BrdCastGroupManageAck{GroupId: req.GroupId}

	grp, ok := DestroyBrdCastGroup(objectId, IdString(req.GroupId))
	if !ok {
		ack.Common = serverCommonAck(ERR_GROUP_NOT_FOUND)
		return ack
	}

	ack.MemberNum = int32(len(grp.Members))
	ack.Common = serverCommonAck(OK)

	return ack
}

func Handle_BrdCastSyncReq(objectId IdString, opCode MsgType, req *BrdCastGroupManageReq) interface{} {
	ack := &BrdCastGroupManageAck{GroupId: req.GroupId}

	grp, ok := FindBrdCastGroup(objectId, IdString(req.GroupId))
	if !ok {
		grp = NewBrdCastGroup4Server(objectId, IdString(req.GroupId))
	}

	ack.MemberNum = grp.ResetMembers()
	ack.Common = serverCommonAck(OK)

	return ack
}
