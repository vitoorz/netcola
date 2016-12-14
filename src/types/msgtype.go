package types

import "fmt"

const(
    MT_Blank             =  MsgType(  0)
)

type MsgType uint32


const(
    msgCodeMask      = 0x0000FFFF
    msgFlagMask      = 0xFFFF0000
)

//to mark a special netmessage
const(
    NetMsgIdFlagClient       = 1 << 16
    NetMsgIdFlagServer       = 1 << 17
    NetMsgIdFlagBroadCast    = 1 << 18//BroadCastId ?

    NetMsgFlagLargeSize      = 1 << 30
    NetMsgFlagMax            = 1 << 31
)

func (t MsgType) Code() MsgType {
    return t & msgCodeMask
}

func (t MsgType) HasFlag(flag uint32) bool {
    return ((uint32(t) & msgFlagMask) & flag) != 0
}

func (mt MsgType) TypeString() string {
    switch {
    case mt.HasFlag(NetMsgIdFlagClient) || mt == 0:
        return "client: " + gnetmsgtypesNames[mt.Code()]
    case mt.HasFlag(NetMsgIdFlagServer):
        return "server: " + innermsgtypesNames[mt.Code()]
    case mt.HasFlag(NetMsgIdFlagBroadCast):
        return "broadcast: " + innermsgtypesNames[mt.Code()]
    }

    return fmt.Sprintf("message with invalid flag: value 0x%x", mt)
}

func AddMsgFlag(t MsgType, flag uint32) MsgType {
    return MsgType(uint32(t) | flag)
}

