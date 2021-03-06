package task

import (
	dm "netcola/library/core/datamsg"
)

const (
	taskUnknown = iota
	taskExit
	taskGetNow
	taskInfo
	taskHelp
	taskLater
	taskMongoCreate
	taskEnv
	taskServiceList
	taskPProf
)

func doNothing(in *dm.DataMsg) {
	in.Receiver = dm.NoReceiver
}

var Route = map[int](func(*dm.DataMsg)){
	taskUnknown:     doNothing,
	taskExit:        doExit,
	taskGetNow:      doGetNow,
	taskInfo:        doInfo,
	taskHelp:        doHelp,
	taskLater:       doLater,
	taskMongoCreate: doMongoCreate,
	taskEnv:         doSysEnv,
	taskServiceList: doService,
	taskPProf:       doPProf,
}
