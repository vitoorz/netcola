package privatetcp

import (
	//"net"
	"sync"
)

type PrivateTCPClient struct {
	Lock sync.RWMutex
	ID   uint64
}
