package play

import (
    . "types"
    "library/frame"
)

var PlayFrame = frame.NewFrame()

func getCommonAck(code int32) *CommonAck {
    return  &CommonAck{
        Status: code,
        Message: ErrDesc[code],
        TimeStamp: uint32(PlayFrame.FrameTime()),
    }
}