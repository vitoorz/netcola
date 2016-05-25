package main

import (
	"math/rand"
	"os"
	"runtime"
	"syscall"
	"time"
)

import (
	//cm "library/core/controlmsg"
	dm "library/core/datamsg"
	"library/idgen"
	"library/logger"
	"server/support"
	"service"
	"service/engine"
	"service/job"
	"service/privatetcp"
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
	logger.Info("main start")
	rand.Seed(time.Now().UTC().Unix())

	idgen.InitIDGen("1")

	support.RegisterSignalHandler(os.Interrupt, InterruptHandler, nil)
	support.RegisterSignalHandler(syscall.SIGTERM, SIGTERMHandler, nil)
	// good idea to stop the world and clean memory before get job
	stopAndCleanMemory()

	bus := dm.NewDataMsgPipe(0, 0)

	enginesrv := engine.NewEngine(bus)
	service.StartService(enginesrv)

	tcpsrv := privatetcp.NewPrivateTCPServer(bus)
	service.StartService(tcpsrv)

	jobsrv := job.NewJob(nil)
	service.StartService(jobsrv)

	for {
		select {
		case sigMsg, ok := <-support.SignalService.Echo:
			if !ok {
				logger.Error("sigint echo error %t", ok)
				continue
			}
			logger.Info("receive:signal echo:%v", sigMsg)
			//enginesrv.Cmd <- &cm.ControlMsg{MsgType: cm.ControlMsgExit}
			//<-enginesrv.Echo
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
