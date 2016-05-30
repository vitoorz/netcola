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
	"service"
	"service/engine"
	"service/job"
	"service/privatetcp"
	"service/timer"
)

func stopAndCleanMemory() {
	memstat := &runtime.MemStats{}
	runtime.ReadMemStats(memstat)
	logger.Info("before gc:memstat.Alloc:%d K", memstat.Alloc/1024)
	runtime.GC()
	runtime.ReadMemStats(memstat)
	logger.Info("after gc:memstat.Alloc:%d K", memstat.Alloc/1024)
}

func main() {
	logger.EnableDebugLog(true)
	logger.Info("main start")
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		logger.Warn("Only tested in linux and mac os operating system, be carefally using in %v", runtime.GOOS)
	}

	rand.Seed(time.Now().UTC().Unix())

	idgen.InitIDGen("1")

	support.RegisterSignalHandler(os.Interrupt, InterruptHandler, nil)
	support.RegisterSignalHandler(syscall.SIGTERM, SIGTERMHandler, nil)
	// good idea to stop the world and clean memory before get job
	stopAndCleanMemory()

	enginesrv := engine.NewEngine()
	service.StartService(enginesrv, "engine", nil)

	tcpsrv := privatetcp.NewPrivateTCPServer()
	service.StartService(tcpsrv, "tcpserver", enginesrv.BUS)

	jobsrv := job.NewJob()
	service.StartService(jobsrv, "job", enginesrv.BUS)

	timersrv := timer.NewTimer()
	service.StartService(timersrv, "timer", enginesrv.BUS)

	for {
		select {
		case sigMsg, ok := <-support.SignalService.Echo:
			if !ok {
				logger.Error("sigint echo error %t", ok)
				continue
			}
			logger.Info("receive:signal echo:%v", sigMsg)
			enginesrv.Exit()
			return
		case sigMsg, ok := <-support.SignalService.Echo:
			if !ok {
				logger.Error("sigterm echo error %t", ok)
				continue
			}
			logger.Info("receive:signal echo:%v", sigMsg)
			enginesrv.Exit()
			return
		}
	}
}
