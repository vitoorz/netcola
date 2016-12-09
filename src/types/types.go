package types

import (
	"fmt"
	"strconv"
)

type ObjectID uint64

func (id ObjectID) ToIdString() IdString {
	return IdString(fmt.Sprintf("0x%x", id))
}

type PlayerId ObjectID

func (p PlayerId) ToIdString() IdString {
	return ObjectID(p).ToIdString()
}

type ServerId ObjectID

func (s ServerId) ToIdString() IdString {
	return ObjectID(s).ToIdString()
}

//unix timestamp
type UnixTS int64

type IdString string

func (id IdString) ToObjectID() ObjectID {
	value, err := strconv.ParseUint(string(id), 0, 64)
	if err != nil {
		fmt.Printf("prase idstring %s to object id error: %s\n", id, err.Error())
	}
	return ObjectID(value)
}

func (id IdString) ToPlayerId() PlayerId {
	return PlayerId(id.ToObjectID())
}

func (id IdString) ToServerId() ServerId {
	return ServerId(id.ToObjectID())
}


const (
	Inner_MsgTypeUnknown  = 0
	Inner_MsgTypeTelnet   = 1
	Inner_MsgTypeC2G      = 2
	Inner_MsgTypeG2C      = 3
	Inner_MsgTypeG2S      = 4
	Inner_MsgTypeS2G      = 5
)