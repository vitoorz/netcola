package env

import (
	"os"
	"strings"
)

//SysEnv will hold all the key-value parameters from the os environment
var SysEnv *Config = nil

func init() {
	SysEnv = NewConfig("")
	osEnv := os.Environ()
	for _, env := range osEnv {
		pair := strings.Split(env, "=")
		SysEnv.Set(pair[0], pair[1])
	}
}
