package idgen

import (
	"testing"
)

func TestInitIDGen(t *testing.T) {
	var ok bool
	expectCase := []string{"16383", "4"}
	badCase := []string{"0x4FFF", "20479", "0", "-65535"}

	for idx, seed := range expectCase {
		ok = InitIDGen(seed)
		if !ok {
			t.Errorf("expectCase:%v,input:%+v failed", idx, seed)
		}
	}

	for idx, seed := range badCase {
		ok = InitIDGen(seed)
		if ok {
			t.Errorf("badCase:%v,input:%+v failed", idx, seed)
		}
	}

}

// test the NewObjectID speed
func BenchmarkGetNewObject(b *testing.B) {
	InitIDGen("126")
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			NewObjectID()
		}
	})
}
