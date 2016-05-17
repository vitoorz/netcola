package main

import (
	"math/rand"
	"runtime"
	"time"
)

import (
	"library/idgen"
	"library/logger"
	//"types"
)

func stopAndCleanMemory() {
	memstat := &runtime.MemStats{}
	runtime.ReadMemStats(memstat)
	logger.Info("before gc:memstat.Alloc:%d M", memstat.Alloc/1024/1024)
	runtime.GC()
	runtime.ReadMemStats(memstat)
	logger.Info("after gc:memstat.Alloc:%d M", memstat.Alloc/1024/1024)
}

func main() {
	logger.Info("main start")
	rand.Seed(time.Now().UTC().Unix())

	idgen.InitIDGen("1")

	// good idea to stop the world and clean memory before get job
	stopAndCleanMemory()
}
