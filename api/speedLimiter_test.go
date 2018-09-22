package api

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestSpeedLimiter_Stop(t *testing.T) {
	sl := NewSpeedLimiter(10, time.Millisecond)
	var reads int

	for range sl.Channel() {
		reads++
		if reads == 25 {
			sl.Stop()
		}
	}

	require.EqualValues(t, 30, reads)
}

func TestNewSpeedLimiter_DoesNotGoOverOrUnderLimit(t *testing.T) {
	reads := 0
	timeUnit := time.Millisecond
	limit := 10
	every := timeUnit * 4
	stopAfter := every * 10
	expectedRedsMin := limit * 10
	expectedRedsMax := expectedRedsMin + limit

	sl := NewSpeedLimiter(limit, every)

	stopTimer := time.NewTimer(stopAfter)
	go func() {
		<-stopTimer.C
		sl.Stop()
	}()

	for range sl.Channel() {
		reads++
	}

	require.True(t, reads <= expectedRedsMax, reads)
	require.True(t, reads >= expectedRedsMin, reads)
}

func TestNewSpeedLimiter_BufferGetsToppedUpByMissingAmount(t *testing.T) {
	maxTimes := 10
	aTimeUnit := time.Second
	sl := newSpeedLimiter(maxTimes, aTimeUnit)

	topUpBuffer(sl)
	require.EqualValues(t, maxTimes, len(sl.allowedToExecute))

	for i := 0; i < 5; i++ {
		<-sl.allowedToExecute
	}
	require.EqualValues(t, 5, len(sl.allowedToExecute))
	topUpBuffer(sl)
	require.EqualValues(t, maxTimes, len(sl.allowedToExecute))
}
