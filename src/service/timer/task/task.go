package task

import (
	dm "library/core/datamsg"
	"library/logger"
)

func DoLater(in *dm.DataMsg) {
	logger.Info("timer:task info")
	in.Receiver = "tcpserver"
	in.Payload = []byte("already 5 sec...\n")
}
