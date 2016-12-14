package manage

import (
	"library/frame"
	. "types"
)

var GateWayFrame = frame.NewFrame()

func getCommonAck(code int32) *CommonAck {
	return &CommonAck{
		Status:    code,
		Message:   ErrDesc[code],
		TimeStamp: uint32(GateWayFrame.FrameTime()),
	}
}
