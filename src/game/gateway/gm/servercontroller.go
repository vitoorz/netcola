package gm

//server req send to gateway, should handle by gateway, and should send no more ack
//objectId is server_id
import . "types"

func Handle_ServerLoginReq(objectId IdString, opCode MsgType, req *ServerLoginReq) interface{} {
	ack := &ServerLoginAck{Status: OK}
	ack.Common = serverCommonAck(OK)
	return ack
}
