package main

import (
	"math/rand"
	"time"
)

import (
	"library/idgen"
	"library/logger"
	//"types"
)

func main() {
	logger.Info("main start")
	rand.Seed(time.Now().UTC().Unix())

	idgen.InitIDGen("1")
}
