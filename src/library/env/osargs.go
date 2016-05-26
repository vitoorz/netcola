package env

import (
	"os"
	"strings"
)

var CmdArgs *Config = nil

func init() {
	CmdArgs = NewConfig("")
	argc := len(os.Args)
	for i := 1; i < argc; i++ {
		arg := os.Args[i]
		kv := strings.Split(arg, "=")
		if len(kv) == 2 {
			CmdArgs.Set(strings.Trim(kv[0], "-"), kv[1])
		} else if strings.Contains(arg, "-") {
			i = i + 1
			CmdArgs.Set(strings.Trim(arg, "-"), os.Args[i])
		} else {
			CmdArgs.Set(arg, true)
		}
	}
}
