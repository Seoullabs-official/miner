package utils

import (
	"fmt"
	"time"
)

func TrackLoopRate(loopCount *uint64) {
	var prevCount uint64
	for range time.Tick(time.Second) {
		currentCount := *loopCount
		loopsPerSecond := currentCount - prevCount
		prevCount = currentCount
		fmt.Printf("Loop Rate: %d loops/sec\n", loopsPerSecond)
	}
}
