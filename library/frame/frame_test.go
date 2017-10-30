package frame

import (
	"testing"
	"time"
)

func TestNewFrame(t *testing.T) {
	frame := NewFrame()
	frame.SetFrameTime()
	a := frame.FrameTime()
	time.Sleep(1 * time.Second)
	b := frame.FrameTime()

	frame.SetFrameTime()
	c := frame.FrameTime()

	if a != b {
		t.Errorf("frame ts changed after 1 second before call SetFrameTime again")
	}

	if a == c {
		t.Errorf("frame ts not changed after 1 second and call SetFrameTime again")
	}
}
