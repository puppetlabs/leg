package clock_test

import (
	"testing"
	"time"

	"github.com/puppetlabs/leg/timeutil/pkg/clock"
	"github.com/puppetlabs/leg/timeutil/pkg/clock/k8sext"
	"github.com/stretchr/testify/assert"
	testclock "k8s.io/apimachinery/pkg/util/clock"
)

func TestTimerCallbackClock(t *testing.T) {
	var last int
	var lastDuration time.Duration

	uc := testclock.NewFakeClock(time.Now())
	mc := clock.NewTimerCallbackClock(k8sext.NewClock(uc), func(d time.Duration) {
		last++
		lastDuration = d
	})

	// Create a new timer.
	timer := mc.NewTimer(5 * time.Second)
	assert.Equal(t, 1, last)
	assert.Equal(t, 5*time.Second, lastDuration)

	// Stop and reset the timer.
	timer.Stop()
	timer.Reset(10 * time.Second)
	assert.Equal(t, 2, last)
	assert.Equal(t, 10*time.Second, lastDuration)
}
