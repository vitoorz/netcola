package controlmsg

const (
	ControlMsgExit = iota
)

type ControlMsg struct {
	MsgType int
	Payload interface{}
}
