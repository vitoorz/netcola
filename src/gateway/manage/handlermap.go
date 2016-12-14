package manage

//Auto generated, do not modify unless you know clearly what you are doing.
import . "types"

type NetMsgHandler func(objectId IdString, opcode MsgType, content []byte) interface{}

type NetMsgCb struct {
	OpCode  MsgType
	RetCode MsgType
	Handler NetMsgHandler
	Desc    string
}

var NetMsgTypeHandler = map[MsgType]*NetMsgCb{
	MT_ServerLoginReq: {MT_ServerLoginReq, MT_ServerLoginAck, On_ServerLoginReq, ""},
	MT_ServerLoginAck: {MT_ServerLoginAck, MT_Blank, On_ServerLoginAck, ""},
}
