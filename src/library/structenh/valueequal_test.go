package structenh

import (
	"testing"
)

func TestValueEqual(t *testing.T) {
	cpy := DeepClone(testData)

	if !ValueEqual(cpy, testData) {
		t.Errorf("\n%s\n--not equal--\n%s", Stringify(cpy), Stringify(testData))
	}
}
