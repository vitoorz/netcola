package main

import (
	"math/rand"
	"os"
	"runtime"
	"syscall"
	"time"
)

import (
	"library/idgen"
	"library/logger"
	"server/support"
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

	sigintEcho := support.RegistorSignalHandler(os.Interrupt, InterruptHandler, nil)
	sigtermEcho := support.RegistorSignalHandler(syscall.SIGTERM, SIGTERMHandler, nil)
	// good idea to stop the world and clean memory before get job
	stopAndCleanMemory()

	for {
		select {
		case sigMsg, ok := <-sigintEcho:
			if !ok {
				logger.Error("SIGINTecho error %t", ok)
				continue
			}
			logger.Info("receive:signal echo:%v", sigMsg)
			return
		case sigMsg, ok := <-sigtermEcho:
			if !ok {
				logger.Error("SIGTERMecho error %t", ok)
				continue
			}
			logger.Info("receive:signal echo:%v", sigMsg)
			return
		}
	}
}
