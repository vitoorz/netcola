package logger

import (
	"log"
	"os"
	"runtime"
)

const (
	debug = " DBG  "
	info  = " INF  "
	warn  = " WAR  "
	error = " ERR  "
	fatal = " FTL  "
)

var l Logger = log.New(os.Stdout, "", log.LstdFlags)

type Logger interface {
	Printf(format string, args ...interface{})
}

func SetLogger(new Logger) {
	l = new
}

func Debug(format string, args ...interface{}) {
	l.Printf(debug+format, args...)
}

func Info(format string, args ...interface{}) {
	l.Printf(info+format, args...)
}

func Warn(format string, args ...interface{}) {
	l.Printf(warn+format, args...)
}

func Error(format string, args ...interface{}) {
	l.Printf(error+format, args...)
}

func Fatal(format string, args ...interface{}) {
	l.Printf(fatal+format, args...)
}

func Stack() {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			break
		}
		buf = make([]byte, 2*len(buf))
	}
	l.Printf(info + string(buf))
}
