package support

import (
	dm "library/core/datamsg"
)

var NetPipe *dm.DataMsgPipe = dm.NewNetMsgPipe(0, 0)
