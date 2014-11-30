package engine

import (
	"sync/atomic"
)

var engineID uint64 = 100

func EngineID() uint64 {
	atomic.AddUint64(&engineID, 1)
	return engineID
}
