package job

import (
//"time"
)

import (
	dm "library/core/datamsg"
	"library/logger"
	"service"
	"service/job/task"
	//ts "types/service"
	"types"
)

func (t *jobType) dataEntry(msg *dm.DataMsg) (operate int, funCode int) {
	defer func() {
		if x := recover(); x != nil {
			logger.Error("job panic: %v", x)
			logger.Stack()
		}
		operate = Continue
		funCode = service.FunPanic
	}()

	logger.Info("%s:get data msg:%d,payload:%v", t.Name, msg.MsgType, msg.Payload.([]byte))

	switch msg.MsgType {
	case types.MsgTypeTelnet:
		choosetask := task.Parse(string(msg.Payload.([]byte)))
		task.Route[choosetask](msg)
	}

	if msg.Payload != nil {
		ok := t.Output.WritePipeNB(msg)
		if !ok {
			// channel full
			return Continue, service.FunDataPipeFull
		}
	}
	return Continue, service.FunOK
}
