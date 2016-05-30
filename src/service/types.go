package service

//over-all data types that would be used in multiple package
type StateT int32

// todo: prababally , state should be put to service
const (
	StateUnknown StateT = iota //in develop stage, we use iota, when mature will change
	StateOK
	StateError
	StateRunning
	StateStop
	// following is corresponding to service interface
	StateInit
	StateStart
	StatePause
	StateExit
	StateKill
)

const (
	FunOK = iota
	FunUnknown
	FunPanic
	FunDataPipeFull
)
