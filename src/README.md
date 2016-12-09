type NetMsgHeadNoId struct {
	OpCode   MsgType
	Size     int32
}

type NetMsgHead struct{
	UserId   PlayerId
	NetMsgHeadNoId
}

#--------------------Client<-----------------------> GateWay <------------------>Server<------------------------
                       . <clientManage>                 : <serverManage>            .
                       . userId = connection->$UserId   :                           .
                       . server = user->$ServerId       : connection = server->$Conn.
                       \.........................................................../
@req                             |                                     |                                       |
protoBinary = req->$protoMarshal |                                     |                                       |
................................................................................................................
                                 |     +    UserId      int64     ===>>|                                       |
OpCode  = req->OpCode            |===>>     OpCode      int32     ===>>| req = un-marshal(protoBinary, OpCode) |
Size    = protoBinary->$Size     |===>>     Size        int32     ===>>|                                       |
content = protoBinary            |===>>     content     []byte    ===>>|                                       |
---------------------------------|-------------------------------------|---------------------------------------|


                                                                         ack = handle(UserId, OpCode, req)
---------------------------------------------------------------------------------------------------------------|
                                 |                                     |  @ack                                 |
                                 |                                     |  protoBinary = ack->$protoMarshal     |
................................................................................................................
                                 |                                     |                                       |
                                 |     -    UserId      int64     <<===|  UserId  = UserId                     |
OpCode      int32                |<<===     OpCode      int32     <<===|  OpCode  = req->$AckOpCode            |
Size        int32                |<<===     Size        int32     <<===|  Size    = protoBinary->$Size         |
content     []byte               |<<===     content     []byte    <<===|  content = protoBinary                |
