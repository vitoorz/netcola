package syssignal

import (
	"os"
	"syscall"
	"testing"
)

func InterruptHandler1(i interface{}) int {
	n := i.(int)
	return n
}

func InterruptHandler3(i interface{}) int {
	n := i.(int)
	return n
}

func eq(re chan int, exp int) bool {
	renum := <-re
	if renum == exp {
		return true
	} else {
		return false
	}
}

// Test signal service.
func TestSignal(t *testing.T) {

	var signalService *SignalService = NewSignalService()
	signalService.InitSignalService()

	echo1 := signalService.RegisterSignalCallback(os.Interrupt, InterruptHandler1, 11)
	echo2 := signalService.RegisterSignalCallback(os.Interrupt, InterruptHandler1, 12)
	echo3 := signalService.RegisterSignalCallback(os.Interrupt, InterruptHandler3, 13)

	t.Logf("test sigint...")
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	if !(eq(echo1, 11) && eq(echo2, 12) && eq(echo3, 13)) {
		t.Error("sigint test failed")
	}
}
