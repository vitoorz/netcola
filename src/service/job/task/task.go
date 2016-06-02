package task

import (
	"fmt"
	dm "library/core/datamsg"
	"library/env"
	"library/logger"
	"os/exec"
	"service"
	"time"
	ts "types/service"
)

func doExit(in *dm.DataMsg) {
	logger.Info("job:task exit")
	in.Receiver = "tcpserver"
	in.Payload = []byte("doExit")
}

func doGetNow(in *dm.DataMsg) {
	logger.Info("job:task now")
	in.Receiver = "tcpserver"
	in.Payload = []byte(time.Now().String() + "\n")
}

func doInfo(in *dm.DataMsg) {
	logger.Info("job:task info")
	var out []byte = []byte("")
	out, err := exec.Command("uname", "-a").Output()
	if err != nil {
		logger.Error("doInfo failed:%s", err.Error())
	}
	in.Receiver = "tcpserver"
	in.Payload = out
}

func doHelp(in *dm.DataMsg) {
	logger.Info("job:task help")
	usage := `
Usage:
    exit: make server exit
    now: get the server time
    info: get the server uname information
    later: echo 'cheers!' 10 sec later
    help: show this help information
`
	in.Receiver = "tcpserver"
	in.Payload = []byte(usage)
}

func doLater(in *dm.DataMsg) {
	in.SetMeta("timer", ts.Event{When: time.Now().Unix() + 5})
	in.Receiver = "timer"
}

func doMongoCreate(in *dm.DataMsg) {
	logger.Info("job:doMongoCreate")
	r := record{}
	//r.Dirty = ts.Dirty{Action: ts.MongoActionCreate}
	in.SetMeta("mongo", &r)
	//in.SetMeta("mongo", ts.MongoDirty{Action: ts.MongoActionCreate})
	in.Receiver = "mongo"
}

func doSysEnv(in *dm.DataMsg) {
	envVal := "System environment variables:\n"
	for key, value := range env.SysEnv.KV {
		envVal = fmt.Sprintf("%s %s = %s\n", envVal, key, value)
	}
	in.Receiver = "tcpserver"
	in.Payload = []byte(envVal)
}

func doService(in *dm.DataMsg) {
	list := "Service list in App:\n"
	for _, name := range service.ServicePool.NameList() {
		list = fmt.Sprintf("%s %s\n", list, name)
	}
	in.Receiver = "tcpserver"
	in.Payload = []byte(list)
}

type record struct {
	//ts.Dirty
}

func (t *record) CRUD() bool {
	logger.Info("I'm talking to mongo")
	return true
}
func (t *record) Inspect() string { return "inspect" }
