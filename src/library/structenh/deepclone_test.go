package structenh

import (
	"testing"
)

type tCeil struct {
	Ca int
	Cb string
}

type tNest struct {
	AInt     int
	BPointer *tCeil
	CSlice   []tCeil
	DSliceP  []*tCeil
	EArray   [2]*tCeil
	FMap     map[string]tCeil
	GMapP    map[string]*tCeil
}

var testData *tNest = &tNest{
	AInt:     9,
	BPointer: &tCeil{1, "ceil_1"},
	CSlice:   []tCeil{tCeil{2, "ceil_2"}, tCeil{3, "ceil3"}},
	DSliceP:  []*tCeil{&tCeil{4, "ceil_4"}, &tCeil{5, "ceil5"}},
	EArray:   [2]*tCeil{&tCeil{6, "ceil_6"}, &tCeil{7, "ceil7"}},
	FMap:     map[string]tCeil{"ceil_8": tCeil{6, "ceil_8"}, "ceil_9": tCeil{9, "ceil9"}},
	GMapP:    map[string]*tCeil{"ceil_10": &tCeil{6, "ceil_10"}, "ceil_11": &tCeil{9, "ceil11"}},
}

func TestDeepClone(t *testing.T) {
	cpy := DeepClone(testData)

	if !ValueEqual(cpy, testData) {
		t.Errorf("cloned value not equal to origin!")
	}
}
