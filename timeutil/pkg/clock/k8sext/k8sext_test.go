package k8sext_test

import (
	"context"
	"testing"
	"time"

	"github.com/puppetlabs/leg/timeutil/pkg/clock/k8sext"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/clock"
)

func TestTimerWithNoDurationImmediatelyFires(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fc := clock.NewFakeClock(time.Now())
	timer := k8sext.NewClock(fc).NewTimer(0)

	// Timer should fire without further intervention.
	select {
	case <-timer.C():
	case <-ctx.Done():
		assert.Fail(t, "context expired waiting for timer")
	}
}
