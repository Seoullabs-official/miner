package utils

import (
	"time"

	logger "github.com/sirupsen/logrus"
)

func TrackLoopRate(loopCount *uint64, logger *logger.Logger) {
	var prevCount uint64
	for range time.Tick(time.Second) {
		currentCount := *loopCount
		loopsPerSecond := currentCount - prevCount
		prevCount = currentCount
		logger.Infof("Loop Rate: %d loops/sec\n", loopsPerSecond)
	}
}
