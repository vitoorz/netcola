package service

import (
	cm "netcola/library/core/controlmsg"
	dm "netcola/library/core/datamsg"
)

type IService interface {
	Start(bus *dm.DataMsgPipe) bool
	Pause() bool
	Resume() bool
	Exit() bool
	Self() *Service
	ControlHandler(*cm.ControlMsg) (int, int)
	DataHandler(*dm.DataMsg) bool
}
