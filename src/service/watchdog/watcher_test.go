package watchdog

import (
	"service"
	"testing"
	"time"
)

func TestWatcher(t *testing.T) {
	handler := func(object string, nano int64) bool {
		t.Logf("user defined handler: watch object %s for %d ns\n", object, nano)
		return false
	}

	object := "test_object"
	w := NewWatcher()
	w.Start()
	w.WatchObjStart(object)
	time.Sleep(3 * time.Second)
	w.WatchObjOver(object)

	w.SetWatcherHandler(handler)
	w.WatchObjStart(object)
	time.Sleep(2 * time.Second)
	w.WatchObjOver(object)
	if w.Status() != service.ServiceStatusRunning {
		t.Errorf("should be running: %s", w.Status())
	}
	w.Stop()
	time.Sleep(1 * time.Second)
	if w.Status() != service.ServiceStatusStopped {
		t.Errorf("should be running: %s", w.Status())
	}
	t.Logf("%v", w)
}
