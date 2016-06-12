package syssignal

import (
	cm "library/core/controlmsg"
	"os"
	"syscall"
	"testing"
)

func InterruptHandler1(i interface{}) *cm.ControlMsg {
	n := i.(int)
	return &cm.ControlMsg{MsgType: n}
}

func InterruptHandler3(i interface{}) *cm.ControlMsg {
	n := i.(int)
	return &cm.ControlMsg{MsgType: n}
}

func eq(s *SignalModule, exp int) bool {
	re := <-s.Echo
	if re.MsgType == exp {
		return true
	} else {
		return false
	}
}

// Test signal service.
// todo: do some parallel test
func TestSignal(t *testing.T) {

	var signalService *SignalModule = NewSignalService()
	signalService.InitSignalService()

	signalService.RegisterSignalCallback(os.Interrupt, InterruptHandler1, 11)
	signalService.RegisterSignalCallback(os.Interrupt, InterruptHandler1, 12)
	signalService.RegisterSignalCallback(os.Interrupt, InterruptHandler3, 13)

	t.Logf("test sigint...")
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	if !(eq(signalService, 11) &&
		eq(signalService, 12) &&
		eq(signalService, 13)) {
		t.Error("sigint test failed")
	}
}
