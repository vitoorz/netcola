package main

import (
	"os"
	"syscall"
)

import (
	cm "library/core/controlmsg"
	"library/logger"
	sig "module/syssignal"
)

func initSystemSignalHandler() *sig.SignalModule {
	singer := sig.NewSignalModule()
	singer.InitSignalModule()
	singer.RegisterSignalCallback(os.Interrupt, InterruptHandler, nil)
	singer.RegisterSignalCallback(syscall.SIGTERM, SIGTERMHandler, nil)
	return singer
}

func InterruptHandler(i interface{}) *cm.ControlMsg {
	logger.Info("InterruptHandler called")
	return &cm.ControlMsg{MsgType: sig.SysSignalHandled}
}

func SIGTERMHandler(i interface{}) *cm.ControlMsg {
	logger.Info("SIGTERMHandler called")
	return &cm.ControlMsg{MsgType: sig.SysSignalHandled}
}
