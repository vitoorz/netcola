package controlmsg

const (
	ControlMsgUnknown = iota
	ControlMsgPause
	ControlMsgResume
	ControlMsgExit
	ControlMsgTick
	ControlMsgDegrade
	ControlMsgMax
)

type ControlMsg struct {
	MsgType int
	Payload interface{}
}
