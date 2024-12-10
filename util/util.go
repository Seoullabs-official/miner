package utils

import (
	"time"

	"github.com/sirupsen/logrus"
)

func TrackLoopRate(loopCount *uint64) {
	var prevCount uint64
	for range time.Tick(time.Second) {
		currentCount := *loopCount
		loopsPerSecond := currentCount - prevCount
		prevCount = currentCount
		logrus.Infof("Loop Rate: %d loops/sec\n", loopsPerSecond)
	}
}
