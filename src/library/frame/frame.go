package frame

import (
	"library/logger"
	"time"
	. "types"
)

const (
	SecondPerDay     = UnixTS(24 * 60 * 60)
	SecondPerHour    = UnixTS(3600)
	SecondPerMinute  = UnixTS(60)
	debugWarnOnDelay = true
)

//reason for frame: Using same timestamp in single request process
//set frame time when we start new request process,
//and (get) frame time during processing, to keep timestamp the same

type Frame struct {
	TimeStamp UnixTS
}

func NewFrame() *Frame {
	return &Frame{
		TimeStamp: UnixTS(time.Now().UTC().Unix()),
	}
}

func (f *Frame) SetFrameTime() {
	f.TimeStamp = UnixTS(time.Now().UTC().Unix())
}

func (f *Frame) FrameTime() UnixTS {
	if debugWarnOnDelay && UnixTS(time.Now().UTC().Unix())-f.TimeStamp > 2 {
		logger.Warn("frame time delay, should SetFrameTime or last procedure last too long")
	}
	return f.TimeStamp
}

func (f *Frame) BeginningOfDay() UnixTS {
	return f.TimeStamp - (f.TimeStamp % SecondPerDay)
}

func (f *Frame) EndingOfDay() UnixTS {
	return f.TimeStamp - (f.TimeStamp % SecondPerDay) + SecondPerDay
}

func (f *Frame) ThisDayAt(h, m, s uint8) UnixTS {
	at := f.BeginningOfDay()
	if h >= 24 || m >= 60 || s >= 60 || h < 0 || m < 0 || s < 0 {
		return at
	}
	return at + UnixTS(h)*SecondPerHour + UnixTS(m)*SecondPerMinute + UnixTS(s)
}
