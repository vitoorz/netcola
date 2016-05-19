package frame

import (
	"library/logger"
	"time"
	. "types"
)

const (
	SecondPerDay     = UnixTS(24 * 60 * 60)
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

func BeginningOfToday() UnixTS {
	ts := UnixTS(time.Now().UTC().Unix())
	return ts - (ts % SecondPerDay)
}

func EndingOfToday() UnixTS {
	ts := UnixTS(time.Now().UTC().Unix())
	return ts - (ts % SecondPerDay) + SecondPerDay
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
