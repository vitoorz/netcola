package types

type StateT int32

const (
	StateUnknown StateT = 0
	StateOK           = 1
	StateError        = 2
)
