package controlmsg

const (
	ControlMsgExit = iota
	ControlMsgTick
	ControlMsgMax
)

type ControlMsg struct {
	MsgType int
	Payload interface{}
}
