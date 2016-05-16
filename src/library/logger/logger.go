package logger

import (
	"log"
	"os"
)

var l Logger = log.New(os.Stdout, "", log.LstdFlags)

type Logger interface {
	Printf(format string, args ...interface{})
}

func SetLogger(new Logger) {
	l = new
}

func Debug(format string, args ...interface{}) {
	l.Printf(format, args...)
}

func Info(format string, args ...interface{}) {
	l.Printf(format, args...)
}

func Warn(format string, args ...interface{}) {
	l.Printf(format, args...)
}

func Error(format string, args ...interface{}) {
	l.Printf(format, args...)
}

func Fatal(format string, args ...interface{}) {
	l.Printf(format, args...)
}
