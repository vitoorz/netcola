package controlmsg

const (
	ControlMsgDummy = iota
	ControlMsgExit
	ControlMsgTick
	ControlMsgMax
)

type ControlMsg struct {
	MsgType int
	Payload interface{}
}
