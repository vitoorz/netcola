package main

import (
	"library/idgen"
	"library/logger"
	"types"
)

func main() {
	logger.Info("start with first id:%v,%v", idgen.LogicObjectDummy, types.TEST)
}
