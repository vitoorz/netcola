package library

import (
	"testing"
	"time"
)

func TestWatcher(t *testing.T) {
	object := "3second"
	w := NewWatcher(nil)
	w.WatchObjStart(object)
	time.Sleep(3 * time.Second)
	w.WatchObjOver(object)
}
