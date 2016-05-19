package support

import (
	"os"
	sig "service/syssignal"
)

var signalService *sig.SignalService = nil

func init() {
	signalService = sig.NewSignalService()
	signalService.InitSignalService()
}

func RegisterSignalHandler(sig os.Signal, f func(interface{}) int, d interface{}) chan int {
	return signalService.RegisterSignalCallback(sig, f, d)
}
