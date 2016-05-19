package support

import (
	"library/netmsg"
)

var NetPipe *netmsg.NetMsgPipe = netmsg.NewNetMsgPipe(0, 0)
