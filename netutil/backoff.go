package netutil

import (
	"context"
	"math"
	"time"

	"github.com/puppetlabs/insights-stdlib/netutil/errors"
)

func Backoff(i int) time.Duration {
	base := 300 * time.Millisecond
	return time.Duration(math.Exp2(float64(i))) * base
}

func Wait(wait time.Duration, ctx context.Context) errors.Error {
	select {
	case <-ctx.Done():
		// timed out or context was cancelled
		return errors.NewBackoffContextCancelledError()
	case <-time.After(wait):
		return nil
	}
}
