package types

//Auto generated, do not modify unless you know clearly what you are doing.
type MsgType int32

const (
    MT_HeatBeat          =  MsgType(  123456)
    MT_Blank             =  MsgType(  0)

    MT_LoginReq          =  MsgType(  1)
    MT_LoginAck          =  MsgType(  2)
    MT_LogoutReq         =  MsgType(  3)
    MT_LogoutAck         =  MsgType(  4)
)

func (mt MsgType)ToString() string {
    return netMsgTypeName[mt]
}

var netMsgTypeName = map[MsgType]string {
    MT_HeatBeat         :  "MT_HeatBeat",
    MT_Blank            :  "MT_Blank",
    MT_LoginReq         :  "MT_LoginReq",
    MT_LoginAck         :  "MT_LoginAck",
    MT_LogoutReq        :  "MT_LogoutReq",
    MT_LogoutAck        :  "MT_LogoutAck",
}
