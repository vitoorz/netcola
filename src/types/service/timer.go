package service

import (
	"time"
)

type Event struct {
	TimerObject *time.Timer
	//Callback    func(interface{}) interface{}
	When int64
}
