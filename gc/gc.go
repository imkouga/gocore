package gc

import (
	"runtime"
	"time"
)

func AutoGC() {

	gcTick := time.Tick(1 * time.Minute)
	for {
		select {
		case <-gcTick:
			runtime.GC()
		}
	}
}
