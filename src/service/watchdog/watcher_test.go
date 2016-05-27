package watchdog

import (
	"testing"
	"time"
)

func TestWatcher(t *testing.T) {
	object := "test_object"
	w := NewWatcher()
	w.Start("watcher", nil)
	w.WatchObjStart(object)
	time.Sleep(3 * time.Second)
	w.WatchObjOver(object)

	w.WatchObjStart(object)
	time.Sleep(2 * time.Second)
	w.WatchObjOver(object)

	time.Sleep(1 * time.Second)
	t.Logf("%v", w)
}
