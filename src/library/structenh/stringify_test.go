package structenh

import "testing"

func TestStringify(t *testing.T) {
	t.Logf(Stringify(testData))
	t.Logf(StringifyValue(testData))
}
