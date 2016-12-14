package play

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
	MT_LoginReq:  {MT_LoginReq, MT_LoginAck, On_LoginReq, "player login request"},
	MT_LoginAck:  {MT_LoginAck, MT_Blank, On_LoginAck, "player login ack from client"},
	MT_LogoutReq: {MT_LogoutReq, MT_LogoutAck, On_LogoutReq, ""},
	MT_LogoutAck: {MT_LogoutAck, MT_Blank, On_LogoutAck, ""},
}
