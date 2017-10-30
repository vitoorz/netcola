package logger

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

const (
	debug = " DBG  "
	info  = " INF  "
	warn  = " WAR  "
	error = " ERR  "
	fatal = " FTL  "
)

var l Logger = log.New(os.Stdout, "", log.LstdFlags)
var debugSwitch bool = false

type Logger interface {
	Printf(format string, args ...interface{})
}

func SetLogger(new Logger) {
	l = new
}

func Debug(format string, args ...interface{}) {
	if debugSwitch {
		l.Printf(debug+format, args...)
	}
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

func EnableDebugLog(s bool) {
	debugSwitch = s
}

func PProf() string {
	profile := pprof.Lookup("goroutine")
	msg := fmt.Sprintf("%d goroutine(s) are currently in proccess, detail...", runtime.NumGoroutine())
	buffer := bytes.NewBuffer([]byte(msg))
	profile.WriteTo(buffer, 3)
	buffer.Write([]byte("-----------over----------"))

	Warn(msg)
	return buffer.String()
}
