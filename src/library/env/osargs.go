package env

import (
	"os"
	"strings"
)

//CmdArgs will hold all the key-value parameters that were given in command line
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
			k := strings.Trim(arg, "-")
			if i == argc-1 || strings.Contains(os.Args[i+1], "-") {
				CmdArgs.Set(k, true)
			} else {
				i = i + 1
				CmdArgs.Set(k, os.Args[i])
			}
		} else {
			CmdArgs.Set(arg, true)
		}
	}
}
