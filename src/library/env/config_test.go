package env

import (
	"testing"
)

func check(t *testing.T, k string, v interface{}, err error) {
	if err != nil {
		t.Errorf("%s: error: %s", k, err.Error())
	} else {
		t.Logf("%s = %v", k, v)
	}
}

func TestNewConfig(t *testing.T) {
	c := NewConfig("test.conf")
	s, err := c.ValString("string")
	check(t, "string", s, err)

	b, err := c.ValString("bool")
	check(t, "bool", b, err)

	i, err := c.ValString("int")
	check(t, "int", i, err)

	f, err := c.ValString("float")
	check(t, "float", f, err)

	check(t, "k-v? ", len(c.KV), nil)

	for k, v := range CmdArgs.KV {
		t.Logf("args: %s = %s", k, v)
	}

	for k, v := range SysEnv.KV {
		t.Logf("osenv: %s = %s", k, v)
	}
}
