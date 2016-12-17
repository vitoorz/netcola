package play

import (
	"game/com"
	. "types"
)

func Handle_LoginReq(objectId IdString, opCode MsgType, req *LoginReq) interface{} {
	grp, ok := com.FindBrdCastGroup(getMyServerId(), objectId)
	if !ok {
		grp = com.NewBrdCastGroup4Server(getMyServerId(), objectId)
		grp.AddMember(objectId)
	}

	notify := grp.GroupDetail()

	AsyncSender.SendServerNotify(MT_BrdCastAddMemberReq, notify)

	ack := &LoginAck{Common: getCommonAck(OK)}

	AsyncSender.SendInstant(MT_LoginAck, NetMsgIdFlagBroadCast, grp.Id, ack)

	return nil
}

func Handle_LogoutReq(objectId IdString, opCode MsgType, req *LogoutReq) interface{} {
	return nil
}
