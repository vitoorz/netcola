package library

import (
	"fmt"
	"testing"
	"time"
)

func TestWatcher(t *testing.T) {
	handler := func(object string, nano int64) bool {
		fmt.Printf("user defined handler: watch object %s for %d ns\n", object, nano)
		return false
	}

	object := "test_object"
	w := NewWatcher(nil)
	w.WatchObjStart(object)
	time.Sleep(3 * time.Second)
	w.WatchObjOver(object)

	w.SetWatcherHandler(handler)
	w.WatchObjStart(object)
	time.Sleep(3 * time.Second)
	w.WatchObjOver(object)
}
