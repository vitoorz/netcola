package gm

import (
	"library/frame"
	. "types"
)

var GatewayFrame = frame.NewFrame()

func serverCommonAck(code int32) *ServerCommonAck {
	return &ServerCommonAck{ErrCode: code}
}
