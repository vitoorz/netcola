package core

//over-all data types that would be used in multiple package
type StateT int32

// todo: prababally , state should be put to service
const (
	StateUnknown StateT = 0
	StateOK             = 1
	StateError          = 2
	StateRunning        = 3
	StateStop           = 4
	StateInit           = 5
)
