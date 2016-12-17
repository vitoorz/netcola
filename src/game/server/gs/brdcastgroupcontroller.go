package gs

import (
	"game/com"
	"game/server/play"
	. "types"
)

func chkAndSyncBrdCastGroup(serverId IdString, ack *BrdCastGroupManageAck) (*com.BrdCastGroup, bool) {
	grp, ok := com.FindBrdCastGroup(serverId, IdString(ack.GroupId))
	if !ok {
		return nil, false
	}

	if ack.MemberNum == grp.MemberNum() {
		return grp, true
	}

	play.AsyncSender.SendServerNotify(MT_BrdCastSyncReq, grp.GroupDetail())

	return grp, true
}

func Handle_BrdCastAddMemberAck(objectId IdString, opCode MsgType, ack *BrdCastGroupManageAck) interface{} {
	chkAndSyncBrdCastGroup(objectId, ack)
	return nil
}

func Handle_BrdCastDelMemberAck(objectId IdString, opCode MsgType, ack *BrdCastGroupManageAck) interface{} {
	chkAndSyncBrdCastGroup(objectId, ack)
	return nil
}

func Handle_BrdCastDestroyAck(objectId IdString, opCode MsgType, ack *BrdCastGroupManageAck) interface{} {
	return nil
}

func Handle_BrdCastSyncAck(objectId IdString, opCode MsgType, ack *BrdCastGroupManageAck) interface{} {
	chkAndSyncBrdCastGroup(objectId, ack)
	return nil
}
