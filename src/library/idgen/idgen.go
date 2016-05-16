package idgen

import (
	"strconv"
	"sync/atomic"
	"time"
)

import (
	"library/logger"
)

const LogicObjectDummy ObjectID = 0

type ObjectID uint64

// nodeSN stores machine id generated once and used in subsequent calls
// to NewObjectId function.
var nodeSN int = 0

// objectIdCounter is atomically incremented when generating a new ObjectId
var objectIDCounter uint32 = 0

func InitIDGen(seed string) {
	// read seed generates machine id and puts it into the nodeSN global
	idint, err := strconv.Atoi(seed)
	if err != nil || idint > 0x3FFF || idint <= 0 {
		nodeSN = 0
		logger.Fatal("wrong config in seed:%v", idint)
		panic("wrong config in idgen seed")
		return
	}
	logger.Info("Init idgen with seed:%d", idint)
	nodeSN = idint
}

func NewObjectID() (id ObjectID) {
	if nodeSN == 0 {
		panic("wrong config in seed")
		return 0
	}
	currentTime := time.Now().UTC().Unix()
	i := atomic.AddUint32(&objectIDCounter, 1) & 0x3FFFFF

	return ObjectID(uint64(currentTime)<<36 | uint64(i)<<14 | uint64(nodeSN))
}

func DummyObjectID() ObjectID {
	return LogicObjectDummy
}
