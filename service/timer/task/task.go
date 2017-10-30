package task

import (
	dm "netcola/library/core/datamsg"
	"netcola/library/logger"
)

func DoLater(in *dm.DataMsg) {
	logger.Info("timer:task info")
	in.Receiver = "tcpserver"
	in.Payload = []byte("already 5 sec...\n")
}
