// +build !race

package api

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestSpeedLimiter_Stop(t *testing.T) {
	sl := NewSpeedLimiter(10, time.Nanosecond)
	var reads int

	for range sl.Channel() {
		reads++
		if reads == 25 {
			sl.Stop()
		}
	}

	require.EqualValues(t, 30, reads)
}
