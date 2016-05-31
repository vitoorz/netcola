package service

import (
	"time"
)

type Event struct {
	TimerObject *time.Timer
	When int64
}
