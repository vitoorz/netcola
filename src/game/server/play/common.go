package play

import (
	"game/server/support"
	"library/frame"
	. "types"
)

var AsyncSender = new(support.AsyncSender)

var PlayFrame = frame.NewFrame()

func getCommonAck(code int32) *CommonAck {
	return &CommonAck{
		Status:    code,
		Message:   ErrDesc[code],
		TimeStamp: uint32(PlayFrame.FrameTime()),
	}
}

func getMyServerId() IdString {
	return support.MyServerId
}
