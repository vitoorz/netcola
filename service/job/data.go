package job

import (
	dm "netcola/library/core/datamsg"
	"netcola/library/logger"
	"netcola/service/job/task"
	"netcola/types"
)

func (t *jobType) DataHandler(msg *dm.DataMsg) bool {
	logger.Info("this is handle:%+v", msg)
	switch msg.MsgType {
	case types.MsgTypeTelnet:
		switch msg.Sender {
		case "mongo":
			msg.Receiver = "tcpserver"
		case "tcpserver":
			choosetask := task.Parse(string(msg.Payload.([]byte)))
			task.Route[choosetask](msg)
		}
	case types.MsgTypeUnknown:
		fallthrough
	default:
		logger.Warn("%s:not handle:get data msg from:%s,type:%d", t.Name, msg.Sender, msg.MsgType)
		msg.Receiver = dm.NoReceiver
	}
	if msg.Receiver != dm.NoReceiver {
		t.Output.WritePipeNoBlock(msg)
	}
	return true
}
