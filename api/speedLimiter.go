package api

import (
	"time"
)

type SpeedLimiter struct {
	limit            int
	ticker           time.Ticker
	pool             int
	allowedToExecute chan bool
	stop             chan bool
}

func NewSpeedLimiter(maxTimes int, aTimeUnit time.Duration) *SpeedLimiter {
	sl := &SpeedLimiter{
		limit:            maxTimes,
		ticker:           *time.NewTicker(aTimeUnit),
		pool:             maxTimes,
		allowedToExecute: make(chan bool, maxTimes),
		stop:             make(chan bool),
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
		for {
			select {
			case <-sl.ticker.C:
				fillBuffer(sl)
			case <-sl.stop:
				return
			}
		}
	}()
}

func (sl *SpeedLimiter) Stop() {
	sl.ticker.Stop()
	sl.stop <- true
	close(sl.allowedToExecute)
	close(sl.stop)
}

func (sl *SpeedLimiter) Channel() <-chan bool {
	return sl.allowedToExecute
}
