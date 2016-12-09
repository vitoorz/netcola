package idgen

import (
	"strconv"
	"sync/atomic"
	"time"
)

import (
	"library/logger"
	. "types"
)

const LogicObjectDummy ObjectID = 0

// nodeSN stores machine id generated once and used in subsequent calls
// to NewObjectId function.
var nodeSN int = 0

// objectIdCounter is atomically incremented when generating a new ObjectId
var objectIDCounter uint32 = 0

//input string should be numbers in base 10
func InitIDGen(seed string) bool {
	// read seed generates machine id and puts it into the nodeSN global
	idint, err := strconv.Atoi(seed)
	if err != nil {
		nodeSN = 0
		logger.Fatal("wrong config in seed,Atoi fail:%v,err:%v", seed, err.Error())
		return false
	}
	if idint > 0x3FFF || idint <= 0 {
		nodeSN = 0
		logger.Fatal("wrong config in seed range:%v", seed)
		return false
	}
	logger.Info("Init idgen with seed:%d", idint)
	nodeSN = idint
	return true
}

func NewObjectID() (id ObjectID) {
	if nodeSN == 0 {
		logger.Fatal("something wrong in idgen seed")
		return 0
	}
	currentTime := time.Now().UTC().Unix()
	i := atomic.AddUint32(&objectIDCounter, 1) & 0x3FFFFF

	return ObjectID(uint64(currentTime)<<36 | uint64(i)<<14 | uint64(nodeSN))
}

func DummyObjectID() ObjectID {
	return LogicObjectDummy
}
