package main

import (
	"library/logger"
	"service/syssignal"
)

func InterruptHandler(i interface{}) int {
	logger.Info("InterruptHandler called")
	return syssignal.SysSignalHandled
}

func SIGTERMHandler(i interface{}) int {
	logger.Info("SIGTERMHandler called")
	return syssignal.SysSignalHandled
}
