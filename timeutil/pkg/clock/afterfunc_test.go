package clock_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/puppetlabs/leg/timeutil/pkg/clock"
	"github.com/puppetlabs/leg/timeutil/pkg/clock/k8sext"
	"github.com/stretchr/testify/assert"
	testclock "k8s.io/apimachinery/pkg/util/clock"
)

func TestAfterFunc(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var calls uint32
	ch := make(chan uint32)

	uc := testclock.NewFakeClock(time.Now())
	mc := k8sext.NewClock(uc)
	timer := clock.AfterFunc(mc, 10*time.Minute, func() {
		// Clear anything currently on the channel. Next read will be the call
		// value.
	clear:
		for {
			select {
			case <-ch:
			default:
				break clear
			}
		}
		ch <- atomic.AddUint32(&calls, 1)
	})

	// Forward the clock 30 minutes. The timer should fire exactly once.
	uc.Step(30 * time.Minute)
	select {
	case i := <-ch:
		assert.Equal(t, uint32(1), i)
	case <-ctx.Done():
		assert.Fail(t, "context expired waiting for timer")
	}

	assert.False(t, timer.Stop())

	// Now reset the timer, which will put the function in a state to be called
	// again.
	assert.False(t, timer.Reset(10*time.Minute))
	uc.Step(5 * time.Minute)

	// Stop the timer and move the clock past.
	assert.True(t, timer.Stop())
	uc.Step(30 * time.Minute)

	// Finally reset once more.
	assert.False(t, timer.Reset(10*time.Minute))

	// Move forward 5 minutes, then reset the timer again. The function should
	// move with us.
	uc.Step(5 * time.Minute)
	assert.True(t, timer.Reset(10*time.Minute))

	uc.Step(5 * time.Minute)
	uc.Step(30 * time.Minute)
	select {
	case i := <-ch:
		assert.Equal(t, uint32(2), i)
	case <-ctx.Done():
		assert.Fail(t, "context expired waiting for timer")
	}
}
