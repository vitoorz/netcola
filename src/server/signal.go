package main

import (
	cm "library/core/controlmsg"
	"library/logger"
	"service/syssignal"
)

func InterruptHandler(i interface{}) *cm.ControlMsg {
	logger.Info("InterruptHandler called")
	return &cm.ControlMsg{MsgType: syssignal.SysSignalHandled}
}

func SIGTERMHandler(i interface{}) *cm.ControlMsg {
	logger.Info("SIGTERMHandler called")
	return &cm.ControlMsg{MsgType: syssignal.SysSignalHandled}
}
