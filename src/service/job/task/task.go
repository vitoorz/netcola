package task

import (
	dm "library/core/datamsg"
	"library/logger"
	"os/exec"
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

func doMongoInsert(in *dm.DataMsg) {
	in.SetMeta("mongo", ts.MongoDirty{Action: ts.MongoActionCreate})
	in.Receiver = "mongo"
}
