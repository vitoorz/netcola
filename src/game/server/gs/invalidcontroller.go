package gs

//gateway ack send to server, should handle by server, and should send no more ack
//objectId is server_id
import . "types"

func Handle_ServerLoginAck(objectId IdString, opCode MsgType, ack *ServerLoginAck) interface{} {
	return nil
}

func Handle_ServerLoginOutReq(objectId IdString, opCode MsgType, ack *ServerLoginOutReq) interface{} {
	return nil
}
