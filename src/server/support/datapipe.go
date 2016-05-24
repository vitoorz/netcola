package support

import (
	"library/datamsg"
)

var NetPipe *datamsg.DataMsgPipe = datamsg.NewNetMsgPipe(0, 0)
