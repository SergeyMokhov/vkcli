package api

import (
	"time"
)

type SpeedLimiter struct {
	limit            int
	ticker           time.Ticker
	allowedToExecute chan bool
	stop             chan bool
}

func NewSpeedLimiter(maxTimes int, every time.Duration) *SpeedLimiter {
	sl := newSpeedLimiter(maxTimes, every)
	topUpBuffer(sl)
	startFillingBufferByTimer(sl)
	return sl
}

func newSpeedLimiter(maxTimes int, aTimeUnit time.Duration) *SpeedLimiter {
	return &SpeedLimiter{
		limit:            maxTimes,
		ticker:           *time.NewTicker(aTimeUnit),
		allowedToExecute: make(chan bool, maxTimes),
		stop:             make(chan bool),
	}
}

func topUpBuffer(sl *SpeedLimiter) {
	topUpTimes := sl.limit - len(sl.allowedToExecute)
	for i := 0; i < topUpTimes; i++ {
		sl.allowedToExecute <- true
	}
}

func startFillingBufferByTimer(sl *SpeedLimiter) {
	go func() {
		for {
			select {
			case <-sl.ticker.C:
				topUpBuffer(sl)
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
