package task

import (
	dm "library/core/datamsg"
)

const (
	taskUnknown = iota
	taskExit
	taskGetNow
	taskInfo
	taskHelp
	taskLater
)

func doNothing(in *dm.DataMsg) {
	in.Payload = nil
}

var Route = map[int](func(*dm.DataMsg)){
	taskUnknown: doNothing,
	taskExit:    doExit,
	taskGetNow:  doGetNow,
	taskInfo:    doInfo,
	taskHelp:    doHelp,
	taskLater:   doLater,
}
