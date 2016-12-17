package types

import (
	"fmt"
	"strconv"
)

const (
	DataMsgFlagC2G      = 1
	DataMsgFlagG2C      = 2
	DataMsgFlagG2S      = 3
	DataMsgFlagS2G      = 4
	DataMsgFlagAsync     = 5
)

//unix timestamp
type UnixTS int64

type ObjectID uint64

func (id ObjectID) ToIdString() IdString {
	return IdString(fmt.Sprintf("0x%x", id))
}

type IdString string

func (id IdString) ToObjectID() ObjectID {
	value, err := strconv.ParseUint(string(id), 0, 64)
	if err != nil {
		fmt.Printf("prase idstring %s to object id error: %s\n", id, err.Error())
	}
	return ObjectID(value)
}
