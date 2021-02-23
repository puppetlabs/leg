package clockctx_test

import (
	"context"
	"testing"
	"time"

	"github.com/puppetlabs/leg/timeutil/pkg/clock/k8sext"
	"github.com/puppetlabs/leg/timeutil/pkg/clockctx"
	"github.com/stretchr/testify/assert"
	testclock "k8s.io/apimachinery/pkg/util/clock"
)

func TestWithDeadlineOfClock(t *testing.T) {
	// Parent context will expire after 10 actual seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uc := testclock.NewFakeClock(time.Now())
	mc := k8sext.NewClock(uc)

	tctx, _ := clockctx.WithDeadlineOfClock(context.Background(), mc, mc.Now().Add(10*time.Minute))

	// Forward by 15 minutes.
	uc.Step(15 * time.Minute)
	select {
	case <-tctx.Done():
	case <-ctx.Done():
		assert.Fail(t, "context expired waiting for deadline of context under test")
	}

	assert.Equal(t, context.DeadlineExceeded, tctx.Err())
}
