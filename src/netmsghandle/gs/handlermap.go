package gs
//Auto generated, do not modify unless you know clearly what you are doing.
import . "types"

type NetMsgHandler func(objectId IdString, opcode MsgType, content []byte ) interface{};

type NetMsgCb struct{
	OpCode  MsgType
	RetCode MsgType
	Handler NetMsgHandler
	Desc    string
}

var NetMsgTypeHandler = map[MsgType]*NetMsgCb {
	MT_ServerLoginReq :&NetMsgCb{MT_ServerLoginReq, MT_ServerLoginAck, On_ServerLoginReq, ""},
	MT_ServerLoginAck :&NetMsgCb{MT_ServerLoginAck, MT_Blank        , On_ServerLoginAck, ""},
	MT_ServerLoginOutReq :&NetMsgCb{MT_ServerLoginOutReq, MT_Blank        , On_ServerLoginOutReq, ""},
	MT_BrdCastAddMemberReq :&NetMsgCb{MT_BrdCastAddMemberReq, MT_BrdCastAddMemberAck, On_BrdCastAddMemberReq, ""},
	MT_BrdCastAddMemberAck :&NetMsgCb{MT_BrdCastAddMemberAck, MT_Blank        , On_BrdCastAddMemberAck, ""},
	MT_BrdCastDelMemberReq :&NetMsgCb{MT_BrdCastDelMemberReq, MT_BrdCastDelMemberAck, On_BrdCastDelMemberReq, ""},
	MT_BrdCastDelMemberAck :&NetMsgCb{MT_BrdCastDelMemberAck, MT_Blank        , On_BrdCastDelMemberAck, ""},
	MT_BrdCastDestroyReq :&NetMsgCb{MT_BrdCastDestroyReq, MT_BrdCastDestroyAck, On_BrdCastDestroyReq, ""},
	MT_BrdCastDestroyAck :&NetMsgCb{MT_BrdCastDestroyAck, MT_Blank        , On_BrdCastDestroyAck, ""},
	MT_BrdCastSyncReq :&NetMsgCb{MT_BrdCastSyncReq, MT_BrdCastSyncAck, On_BrdCastSyncReq, ""},
	MT_BrdCastSyncAck :&NetMsgCb{MT_BrdCastSyncAck, MT_Blank        , On_BrdCastSyncAck, ""},
}
