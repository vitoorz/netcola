package service

import (
	"time"
)

type event struct {
	TimerObject *time.Timer
	Callback    func()
	WakeTime    uint64
}
