package sim
//Auto generated, do not modify unless you know clearly what you are doing.

import . "types"

func Play_InvalidReq(playerId IdString, opCode MsgType, req interface{}) interface{} {

	return getCommonAck(ERR_INVALID_REQ)
}

func Play_LoginReq(playerId IdString, opCode MsgType, req *LoginReq) interface{} {
	return nil
}

func Play_LogoutReq(playerId IdString, opCode MsgType, req *LogoutReq) interface{} {
	return nil
}
