package types

//over-all data types that would be used in multiple package
type StateT int32

const (
	StateUnknown StateT = 0
	StateOK           = 1
	StateError        = 2
)

//unix timestamp
type UnixTS int64
