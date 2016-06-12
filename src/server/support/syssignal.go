package support

import (
	cm "library/core/controlmsg"
	sig "module/syssignal"
	"os"
)

var SignalService *sig.SignalService = nil

func init() {
	SignalService = sig.NewSignalService()
	SignalService.InitSignalService()
}

func RegisterSignalHandler(sig os.Signal, f func(interface{}) *cm.ControlMsg, d interface{}) {
	SignalService.RegisterSignalCallback(sig, f, d)
}
