package api

import (
	"time"
)

type SpeedLimiter struct {
	limit            int
	ticker           time.Ticker
	pool             int
	allowedToExecute chan bool
}

func NewSpeedLimiter(maxTimes int, aTimeUnit time.Duration) *SpeedLimiter {
	sl := &SpeedLimiter{
		limit:            maxTimes,
		ticker:           *time.NewTicker(aTimeUnit),
		pool:             maxTimes,
		allowedToExecute: make(chan bool, maxTimes),
	}
	startFillingBuffer(sl)
	return sl
}

func fillBuffer(sl *SpeedLimiter) {
	for i := 0; i < sl.limit; i++ {
		sl.allowedToExecute <- true
	}
}

func startFillingBuffer(sl *SpeedLimiter) {
	go func() {
		defer func() { recover() }()

		for range sl.ticker.C {
			fillBuffer(sl)
		}
	}()
}

func (sl *SpeedLimiter) Stop() {
	sl.ticker.Stop()
	close(sl.allowedToExecute)
}

func (sl *SpeedLimiter) Channel() <-chan bool {
	return sl.allowedToExecute
}
