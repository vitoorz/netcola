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

const (
	ProcessStatOK = iota
	ProcessStatIgnore
	ProcessStatPanic
	ProcessPipeFull
	ProcessPipeReceiverLost
	ProcessStatUnknown
)

const (
	NextActionBreak = iota
	NextActionContinue
	NextActionReturn
)

type ControlMsg struct {
	MsgType int
	Payload interface{}
}
