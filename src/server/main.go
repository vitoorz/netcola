package main

import (
	"math/rand"
	"os"
	"runtime"
	"syscall"
	"time"
)

import (
	cm "library/controlmsg"
	"library/idgen"
	"library/logger"
	"server/engine"
	"server/support"
)

func stopAndCleanMemory() {
	memstat := &runtime.MemStats{}
	runtime.ReadMemStats(memstat)
	logger.Info("before gc:memstat.Alloc:%d M", memstat.Alloc/1024)
	runtime.GC()
	runtime.ReadMemStats(memstat)
	logger.Info("after gc:memstat.Alloc:%d M", memstat.Alloc/1024)
}

func main() {
	logger.Info("main start")
	rand.Seed(time.Now().UTC().Unix())

	idgen.InitIDGen("1")

	sigintEcho := support.RegisterSignalHandler(os.Interrupt, InterruptHandler, nil)
	sigtermEcho := support.RegisterSignalHandler(syscall.SIGTERM, SIGTERMHandler, nil)
	// good idea to stop the world and clean memory before get job
	stopAndCleanMemory()

	e := engine.NewEngineDefine()
	e.StartEngine(support.NetPipe)
	ekey := e.ControlEntry()

	for {
		select {
		case sigMsg, ok := <-sigintEcho:
			if !ok {
				logger.Error("sigint echo error %t", ok)
				continue
			}
			logger.Info("receive:signal echo:%v", sigMsg)
			ekey.Cmd <- &cm.ControlMsg{MsgType: cm.ControlMsgExit}
			<-ekey.Echo
			return
		case sigMsg, ok := <-sigtermEcho:
			if !ok {
				logger.Error("sigterm echo error %t", ok)
				continue
			}
			logger.Info("receive:signal echo:%v", sigMsg)
            ekey.Cmd <- &cm.ControlMsg{MsgType: cm.ControlMsgExit}
            <-ekey.Echo
			return
		}
	}
}
