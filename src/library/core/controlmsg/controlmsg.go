package controlmsg

const (
	ControlMsgUnknown = iota
	ControlMsgPause
	ControlMsgResume
	ControlMsgExit
	ControlMsgTick
	ControlMsgMax
)

type ControlMsg struct {
	MsgType int
	Payload interface{}
}
