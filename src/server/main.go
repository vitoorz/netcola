package main

import (
	"math/rand"
	"os"
	"runtime"
	"syscall"
	"time"
)

import (
	cm "library/core/controlmsg"
	"library/idgen"
	"library/logger"
	"server/support"
	"service"
	"service/job"
	"service/mongo"
	"service/privatetcp"
	"service/ticker"
	"service/timer"
	"service/watchdog"
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

	enginesrv := service.NewEngine("engine")
	enginesrv.Start(nil)

	tickersrv := ticker.NewTicker("ticker")
	service.StartService(tickersrv, enginesrv.BUS)

	jobsrv := job.NewJob("job")
	service.StartService(jobsrv, enginesrv.BUS)

	timersrv := timer.NewTimer("timer")
	service.StartService(timersrv, enginesrv.BUS)

	mongosrv := mongo.NewMongo("mongo", "127.0.0.1", "27017")
	service.StartService(mongosrv, enginesrv.BUS)

	tcpsrv := privatetcp.NewPrivateTCPServer("tcpserver", "0.0.0.0", "7171")
	service.StartService(tcpsrv, enginesrv.BUS)

	watcher := watchdog.NewWatcher("watcher")
	service.StartService(watcher, enginesrv.BUS)

	tickersrv.BookedBy(watcher, 1)

	for {
		select {
		case sigMsg, ok := <-support.SignalService.Echo:
			if !ok {
				logger.Error("sigint echo error %t", ok)
				continue
			}
			logger.Info("receive:signal echo:%v", sigMsg)
			enginesrv.Cmd <- &cm.ControlMsg{MsgType: cm.ControlMsgExit}
			<-enginesrv.Echo
			return
		case sigMsg, ok := <-support.SignalService.Echo:
			if !ok {
				logger.Error("sigterm echo error %t", ok)
				continue
			}
			logger.Info("receive:signal echo:%v", sigMsg)
			enginesrv.Cmd <- &cm.ControlMsg{MsgType: cm.ControlMsgExit}
			<-enginesrv.Echo
			return
		}
	}
}
