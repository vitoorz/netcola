package play

//Auto generated, do not modify unless you know clearly what you are doing.

import . "types"

func Handle_InvalidReq(objectId IdString, opCode MsgType, req interface{}) interface{} {
	return getCommonAck(ERR_INVALID_REQ)
}

func Handle_LoginReq(objectId IdString, opCode MsgType, req *LoginReq) interface{} {
	return &LoginAck{Common: getCommonAck(OK)}
}

func Handle_LogoutReq(objectId IdString, opCode MsgType, req *LogoutReq) interface{} {
	return nil
}
